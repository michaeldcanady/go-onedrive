package vfs

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

type orchestrator struct {
	mounts     mount.Service
	plugins    plugins.Manager
	tokens     identity.TokenService
	identities identity.Service
	logger     logger.Service
}

// NewOrchestrator returns a new [VFS] implementation that orchestrates file operations
// by routing them to the appropriate storage plugins based on active [mount.Mount] points.
func NewOrchestrator(ms mount.Service, pm plugins.Manager, ts identity.TokenService, is identity.Service, l logger.Service) VFS {
	return &orchestrator{
		mounts:     ms,
		plugins:    pm,
		tokens:     ts,
		identities: is,
		logger:     l,
	}
}

func (o *orchestrator) List(ctx context.Context, path string) ([]*Node, error) {
	client, relPath, options, err := o.prepare(ctx, path)
	if err != nil {
		return nil, err
	}

	resp, err := client.List(ctx, &storage_proto.ListRequest{
		Path:    relPath,
		Options: options,
	})
	if err != nil {
		return nil, plugins.FromGRPC(err)
	}

	nodes := make([]*Node, len(resp.Nodes))
	for i, n := range resp.Nodes {
		nodes[i] = FromProtoNode(n)
	}
	return nodes, nil
}

func (o *orchestrator) Stat(ctx context.Context, path string) (*Node, error) {
	client, relPath, options, err := o.prepare(ctx, path)
	if err != nil {
		return nil, err
	}

	resp, err := client.Stat(ctx, &storage_proto.StatRequest{
		Path:    relPath,
		Options: options,
	})
	if err != nil {
		return nil, plugins.FromGRPC(err)
	}

	return FromProtoNode(resp.Node), nil
}

func (o *orchestrator) Mkdir(ctx context.Context, path string) error {
	client, relPath, options, err := o.prepare(ctx, path)
	if err != nil {
		return err
	}

	_, err = client.Mkdir(ctx, &storage_proto.MkdirRequest{
		Path:    relPath,
		Options: options,
	})
	return plugins.FromGRPC(err)
}

func (o *orchestrator) Remove(ctx context.Context, path string) error {
	client, relPath, options, err := o.prepare(ctx, path)
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, &storage_proto.DeleteRequest{
		Path:    relPath,
		Options: options,
	})
	return plugins.FromGRPC(err)
}

func (o *orchestrator) prepare(ctx context.Context, path string) (storage_proto.StorageServiceClient, string, map[string]string, error) {
	m, relPath, err := o.resolvePath(ctx, path)
	if err != nil {
		return nil, "", nil, err
	}

	client, err := o.getBackend(m)
	if err != nil {
		return nil, "", nil, err
	}

	token, _ := o.getToken(ctx, m)
	options := o.getOptions(m, token)

	return client, relPath, options, nil
}

func (o *orchestrator) Move(ctx context.Context, src, dst string) error {
	srcM, srcRel, err := o.resolvePath(ctx, src)
	if err != nil {
		return err
	}
	dstM, dstRel, err := o.resolvePath(ctx, dst)
	if err != nil {
		return err
	}

	if srcM.Path == dstM.Path {
		// Same mount, try native move
		client, err := o.getBackend(srcM)
		if err != nil {
			return err
		}

		token, _ := o.getToken(ctx, srcM)
		options := o.getOptions(srcM, token)

		_, err = client.Move(ctx, &storage_proto.MoveRequest{
			Source:      srcRel,
			Destination: dstRel,
			Options:     options,
		})
		return plugins.FromGRPC(err)
	}

	// Cross-mount move: copy then delete
	if err := o.Copy(ctx, src, dst); err != nil {
		return err
	}
	return o.Remove(ctx, src)
}

func (o *orchestrator) Copy(ctx context.Context, src, dst string) error {
	reader, err := o.Read(ctx, src)
	if err != nil {
		return err
	}
	defer reader.Close()

	return o.Write(ctx, dst, reader)
}

func (o *orchestrator) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	client, relPath, options, err := o.prepare(ctx, path)
	if err != nil {
		return nil, err
	}

	stream, err := client.Read(ctx, &storage_proto.ReadRequest{
		Path:    relPath,
		Options: options,
	})
	if err != nil {
		return nil, plugins.FromGRPC(err)
	}

	pr, pw := io.Pipe()
	l := logger.WithContext(o.logger, ctx)
	go func() {
		defer pw.Close()
		for {
			select {
			case <-ctx.Done():
				if err := pw.CloseWithError(ctx.Err()); err != nil {
					l.Warn("failed to close pipe with error", "error", ctx.Err(), "close_error", err)
				}
				return
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					if cerr := pw.CloseWithError(plugins.FromGRPC(err)); cerr != nil {
						l.Warn("failed to close pipe with error", "error", err, "close_error", cerr)
					}
					return
				}
				if _, err = pw.Write(resp.Chunk); err != nil {
					return
				}
			}
		}
	}()

	return pr, nil
}

func (o *orchestrator) Write(ctx context.Context, path string, reader io.Reader, options ...WriteOption) error {
	client, relPath, opts, err := o.prepare(ctx, path)
	if err != nil {
		return err
	}

	// Apply functional options
	for _, opt := range options {
		opt(opts)
	}

	stream, err := client.Write(ctx)
	if err != nil {
		return plugins.FromGRPC(err)
	}

	buf := make([]byte, 32*1024)
	first := true

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 || first {
			if err := stream.Send(&storage_proto.WriteRequest{
				Path:    relPath,
				Chunk:   buf[:n],
				Options: opts,
			}); err != nil {
				return plugins.FromGRPC(err)
			}
			first = false
		}

		if err == io.EOF {
			break
		}
	}

	_, err = stream.CloseAndRecv()
	return plugins.FromGRPC(err)
}

func (o *orchestrator) resolvePath(ctx context.Context, p string) (*mount.Mount, string, error) {
	p = path.Clean(p)
	mounts, err := o.mounts.List(ctx)
	if err != nil {
		return nil, "", err
	}

	var bestMatch *mount.Mount
	for _, m := range mounts {
		// Ensure it's a full segment match: either exact match or prefix followed by /
		if p == m.Path || strings.HasPrefix(p, m.Path+"/") {
			if bestMatch == nil || len(m.Path) > len(bestMatch.Path) {
				bestMatch = m
			}
		}
	}

	if bestMatch == nil {
		return nil, "", fmt.Errorf("no mount found for path: %s", p)
	}

	relPath := strings.TrimPrefix(p, bestMatch.Path)
	if !strings.HasPrefix(relPath, "/") {
		relPath = "/" + relPath
	}

	return bestMatch, relPath, nil
}

func (o *orchestrator) getBackend(m *mount.Mount) (storage_proto.StorageServiceClient, error) {
	pluginName := fmt.Sprintf("storage-%s", m.Type)
	return o.plugins.GetStoragePlugin(pluginName)
}

func (o *orchestrator) getToken(ctx context.Context, m *mount.Mount) (*identity.Token, error) {
	l := logger.WithContext(o.logger, ctx)

	if m.IdentityID == "" {
		return nil, nil
	}

	provider := m.IdentityProvider

	// 1. Resolve identity by ID or email/display name using centralized service
	targetIdentity, err := o.identities.FindIdentity(ctx, m.IdentityID)
	if err != nil {
		return nil, fmt.Errorf("identity %s not found: %w", m.IdentityID, err)
	}

	resolvedID := targetIdentity.ID
	if provider == "" {
		provider = targetIdentity.Provider
		l.Debug("resolved missing provider from identity", "provider", provider, "identity", resolvedID)
	}

	if provider == "" {
		provider = "azure" // ultimate fallback
	}

	// 2. Lookup token using resolved provider and ID
	token, err := o.tokens.GetToken(ctx, provider, resolvedID)
	if err != nil {
		return nil, fmt.Errorf("token not found for %s:%s: %w", provider, resolvedID, err)
	}

	return token, nil
}

func (o *orchestrator) getOptions(m *mount.Mount, t *identity.Token) map[string]string {
	opts := make(map[string]string)
	for k, v := range m.Options {
		opts[k] = v
	}
	if t != nil {
		opts["token"] = t.AccessToken
	}
	return opts
}
