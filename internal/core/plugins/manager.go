package plugins

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
	storage_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/storage"
)

// HandshakeConfig ensures that the host and the plugin processes are compatible.
var HandshakeConfig = plugin.HandshakeConfig{
	MagicCookieKey:   "ODC_PLUGIN",
	MagicCookieValue: "odc",
}

const (
	// PluginTypeStorage is the identifier for storage plugins.
	PluginTypeStorage = "storage"
	// PluginTypeIdentity is the identifier for identity plugins.
	PluginTypeIdentity = "identity"
	// HeaderRequestID is the HTTP header key for the request ID.
	HeaderRequestID = "x-request-id"
)

type cachedPlugin struct {
	client     *plugin.Client
	rpcClient  plugin.ClientProtocol
	grpcClient *grpc.ClientConn
	service    any
	cleanup    func()
}

type pluginManager struct {
	config  config.Service
	logger  logger.Service
	repo    Repository

	mu      sync.Mutex
	plugins map[string]*cachedPlugin
}

// NewPluginManager returns a new [Manager] initialized with the required dependencies.
func NewPluginManager(cs config.Service, l logger.Service, repo Repository) Manager {
	return &pluginManager{
		config:  cs,
		logger:  l,
		repo:    repo,
		plugins: make(map[string]*cachedPlugin),
	}
}


func (m *pluginManager) GetStoragePlugin(name string) (storage_proto.StorageServiceClient, error) {
	return &storageProxy{manager: m, name: name}, nil
}

func (m *pluginManager) GetIdentityPlugin(name string) (identity_proto.IdentityPluginClient, error) {
	return &identityProxy{manager: m, name: name}, nil
}

func (m *pluginManager) getRawStoragePlugin(name string) (storage_proto.StorageServiceClient, error) {
	raw, err := m.dispensePlugin(name, PluginTypeStorage)
	if err != nil {
		return nil, err
	}
	return raw.(storage_proto.StorageServiceClient), nil
}

func (m *pluginManager) getRawIdentityPlugin(name string) (identity_proto.IdentityPluginClient, error) {
	raw, err := m.dispensePlugin(name, PluginTypeIdentity)
	if err != nil {
		return nil, err
	}
	return raw.(identity_proto.IdentityPluginClient), nil
}

func (m *pluginManager) ListPlugins(ctx context.Context) ([]*Metadata, error) {
	l := logger.WithContext(m.logger, ctx)

	pluginDir := m.getPluginDir()
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var results []*Metadata
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		path := filepath.Join(pluginDir, f.Name())
		info, err := os.Stat(path)
		if err != nil {
			l.Warn("unable to stat file", "path", path, "error", err)
			continue
		}

		// Only attempt to load if it's an executable file
		if info.Mode()&0111 == 0 {
			continue
		}

		if meta := m.getPluginMetadata(ctx, f.Name()); meta != nil {
			results = append(results, meta)
			if err := m.repo.Set(f.Name(), meta); err != nil {
				l.Warn("failed to cache plugin metadata", "path", f.Name(), "error", err)
			}
		}
	}

	return results, nil
}

func (m *pluginManager) getPluginMetadata(ctx context.Context, name string) *Metadata {
	// Check cache
	if meta, err := m.repo.Get(name); err == nil && meta != nil {
		return meta
	}

	// Try as storage first
	if client, err := m.getRawStoragePlugin(name); err == nil {
		if meta, err := client.GetMetadata(ctx, &storage_proto.MetadataRequest{}); err == nil {
			return &Metadata{
				Name:               meta.Name,
				Type:               meta.Type,
				SupportedProviders: meta.SupportedProviders,
				PluginPath:         name,
			}
		}
	}

	// Try as identity
	if client, err := m.getRawIdentityPlugin(name); err == nil {
		if meta, err := client.GetMetadata(ctx, &identity_proto.MetadataRequest{}); err == nil {
			return &Metadata{
				Name:               meta.Name,
				Type:               meta.Type,
				SupportedProviders: meta.SupportedProviders,
				PluginPath:         name,
			}
		}
	}

	return nil
}

func (m *pluginManager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.plugins {
		if p.cleanup != nil {
			p.cleanup()
		}
	}
	m.plugins = make(map[string]*cachedPlugin)
	return nil
}

func (m *pluginManager) dispensePlugin(name, pluginType string) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cacheKey := fmt.Sprintf("%s:%s", name, pluginType)
	if p, ok := m.plugins[cacheKey]; ok {
		// Check if connection is still alive
		if p.grpcClient != nil && p.grpcClient.GetState() != connectivity.Shutdown {
			return p.service, nil
		}
		// Connection lost or shutdown, cleanup and re-dispense
		if p.cleanup != nil {
			p.cleanup()
		}
		delete(m.plugins, cacheKey)
	}

	client, cleanup, err := m.getClient(name, pluginType)
	if err != nil {
		return nil, err
	}

	rpcClient, err := client.Client()
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to get rpc client: %w", err)
	}

	// Get gRPC Client connection for health checks
	grpcClient := rpcClient.(*plugin.GRPCClient).Conn

	raw, err := rpcClient.Dispense(pluginType)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to dispense %s plugin: %w", pluginType, err)
	}

	m.plugins[cacheKey] = &cachedPlugin{
		client:     client,
		rpcClient:  rpcClient,
		grpcClient: grpcClient,
		service:    raw,
		cleanup:    cleanup,
	}

	return raw, nil
}

func (m *pluginManager) getClient(name string, pluginType string) (*plugin.Client, func(), error) {
	pluginDir := m.getPluginDir()
	logDir := m.getLogDir()

	path := filepath.Join(pluginDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("plugin not found: %s", path)
	}

	pluginMap := make(map[string]plugin.Plugin)
	if pluginType == PluginTypeStorage {
		pluginMap[PluginTypeStorage] = &StorageGRPCPlugin{}
	} else {
		pluginMap[PluginTypeIdentity] = &IdentityGRPCPlugin{}
	}

	cmd := exec.Command(path)
	cmd.Stdin = os.Stdin

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, fmt.Sprintf("plugin-%s.log", name))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open plugin log file: %w", err)
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		SyncStderr:      logFile,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Level:  hclog.Warn,
			Output: logFile,
		}),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		GRPCDialOptions: []grpc.DialOption{
			grpc.WithUnaryInterceptor(requestIDUnaryInterceptor),
			grpc.WithStreamInterceptor(requestIDStreamInterceptor),
		},
	})

	return client, func() {
		client.Kill()
		logFile.Close()
	}, nil
}

func requestIDUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(addRequestIDToOutgoing(ctx), method, req, reply, cc, opts...)
}

func requestIDStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return streamer(addRequestIDToOutgoing(ctx), desc, cc, method, opts...)
}

func addRequestIDToOutgoing(ctx context.Context) context.Context {
	if rid := logger.GetRequestID(ctx); rid != "" {
		return metadata.AppendToOutgoingContext(ctx, HeaderRequestID, rid)
	}
	return ctx
}

// CustomGRPCServer returns a gRPC server configured with request ID interceptors for distributed tracing.
func CustomGRPCServer(opts []grpc.ServerOption) *grpc.Server {
	opts = append(opts,
		grpc.UnaryInterceptor(requestIDServerUnaryInterceptor),
		grpc.StreamInterceptor(requestIDServerStreamInterceptor),
	)
	return grpc.NewServer(opts...)
}

func requestIDServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return handler(extractRequestIDFromIncoming(ctx), req)
}

func requestIDServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := extractRequestIDFromIncoming(ss.Context())
	return handler(srv, &wrappedStream{ss, ctx})
}

func extractRequestIDFromIncoming(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get(HeaderRequestID); len(ids) > 0 {
			return logger.WithRequestID(ctx, ids[0])
		}
	}
	return ctx
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// StorageGRPCPlugin implements the [plugin.GRPCPlugin] interface for storage services.
type StorageGRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl storage_proto.StorageServiceServer
}

// GRPCServer registers the storage service implementation with the gRPC server.
func (p *StorageGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	storage_proto.RegisterStorageServiceServer(s, p.Impl)
	return nil
}

// GRPCClient returns a storage service client for the given gRPC connection.
func (p *StorageGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return storage_proto.NewStorageServiceClient(c), nil
}

// IdentityGRPCPlugin implements the [plugin.GRPCPlugin] interface for identity services.
type IdentityGRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl identity_proto.IdentityPluginServer
}

// GRPCServer registers the identity service implementation with the gRPC server.
func (p *IdentityGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	identity_proto.RegisterIdentityPluginServer(s, p.Impl)
	return nil
}

// GRPCClient returns an identity service client for the given gRPC connection.
func (p *IdentityGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return identity_proto.NewIdentityPluginClient(c), nil
}

func (m *pluginManager) getPluginDir() string {
	if val, err := m.config.Get(config.KeyCorePluginsDir); err == nil && val != nil {
		if dir := fmt.Sprintf("%v", val); dir != "" {
			return dir
		}
	}

	// Default: ~/.config/odc/plugins
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "odc", "plugins")
}

func (m *pluginManager) getLogDir() string {
	if val, err := m.config.Get(config.KeyCoreLogDir); err == nil && val != nil {
		if dir := fmt.Sprintf("%v", val); dir != "" {
			return dir
		}
	}

	// Default: ~/.config/odc/logs
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "odc", "logs")
}
