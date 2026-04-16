package onedrive

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/fs/providers"
	platform "github.com/michaeldcanady/go-onedrive/internal/identity/providers/shared"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func init() {
	providers.Register(providers.Descriptor{
		Name: providerName,
		Factory: func(deps providers.Dependencies) (fs.Service, error) {
			p, ok := deps.Get("platform")
			if !ok {
				return nil, fmt.Errorf("platform not found in dependencies")
			}
			plat, ok := p.(platform.PlatformProvider)
			if !ok {
				return nil, fmt.Errorf("invalid platform type")
			}

			dr, ok := deps.Get("drive_resolver")
			if !ok {
				return nil, fmt.Errorf("drive resolver not found in dependencies")
			}
			resolver, ok := dr.(fs.DriveResolver)
			if !ok {
				return nil, fmt.Errorf("invalid drive resolver type")
			}

			return NewProvider(plat, resolver, deps.Logger()), nil
		},
	})
}

const (
	baseURL = "https://graph.microsoft.com/v1.0"
	// URI Templates from old implementation
	rootURITemplate                      = "{+baseurl}/drives/{drive_id}/root"
	rootRelativeURITemplate              = "{+baseurl}/drives/{drive_id}/root:{path}:"
	rootChildrenURITemplate              = "{+baseurl}/drives/{drive_id}/root/children"
	rootRelativeChildrenURITemplate      = "{+baseurl}/drives/{drive_id}/root:{path}:/children"
	rootRelativeContentURITemplate       = "{+baseurl}/drives/{drive_id}/root:{path}:/content"
	rootRelativeCreateSessionURITemplate = "{+baseurl}/drives/{drive_id}/root:{path}:/createUploadSession"
	providerName                         = "onedrive"
	// uploadThreshold is the file size at which we switch to resumable uploads (4MB).
	uploadThreshold = 4 * 1024 * 1024
	// uploadChunkSize is the size of each chunk in a resumable upload (must be multiple of 320 KiB).
	uploadChunkSize = 320 * 1024 * 10 // 3.2 MiB
)

// Provider implements the filesystem Service interface for Microsoft OneDrive.
type Provider struct {
	platform      platform.PlatformProvider
	driveResolver fs.DriveResolver
	log           logger.Logger
}

// NewProvider creates a new instance of the OneDrive filesystem provider.
func NewProvider(p platform.PlatformProvider, dr fs.DriveResolver, log logger.Logger) *Provider {
	return &Provider{
		platform:      p,
		driveResolver: dr,
		log:           log,
	}
}

func (p *Provider) Name() string {
	return providerName
}

// Get retrieves metadata for a single item by its OneDrive structured URI.
func (p *Provider) Get(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	driveID := resolveDriveID(ctx, uri, p.driveResolver)
	log.Debug("resolved drive id", logger.String("drive_id", driveID))

	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, uri.Path)
	log.Debug("expanded uri", logger.String("uri", url))

	log.Debug("retrieving adapter")
	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return fs.Item{}, mapError(err, uri.Path)
	}

	builder := drives.NewItemRootRequestBuilder(url, adapter)

	log.Debug("sending get item request")
	it, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to get item", logger.Error(err))
		return fs.Item{}, mapError(err, uri.Path)
	}

	log.Info("retrieved item metadata")
	return mapItemToSharedItem(it, uri.Path), nil
}

// Stat returns metadata for an item at the specified URI.
func (p *Provider) Stat(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	return p.Get(ctx, uri)
}

// List returns the children of a directory in OneDrive.
func (p *Provider) List(ctx context.Context, uri *fs.URI, opts fs.ListOptions) ([]fs.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	driveID := resolveDriveID(ctx, uri, p.driveResolver)
	log.Debug("resolved drive id", logger.String("drive_id", driveID))

	url := expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, uri.Path)
	log.Debug("expanded uri", logger.String("uri", url))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return nil, mapError(err, uri.Path)
	}

	builder := drives.NewItemItemsRequestBuilder(url, adapter)

	log.Debug("sending list children request")
	resp, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to list children", logger.Error(err))
		return nil, mapError(err, uri.Path)
	}

	var items []fs.Item
	for _, it := range resp.GetValue() {
		name := ""
		if it.GetName() != nil {
			name = *it.GetName()
		}
		childPath := path.Join(uri.Path, name)
		items = append(items, mapItemToSharedItem(it, childPath))
	}

	log.Info("retrieved children", logger.Int("count", len(items)))
	return items, nil
}

// ReadFile provides an io.ReadCloser for the content of a file in OneDrive.
func (p *Provider) ReadFile(ctx context.Context, uri *fs.URI, opts fs.ReadOptions) (io.ReadCloser, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	driveID := resolveDriveID(ctx, uri, p.driveResolver)
	log.Debug("resolved drive id", logger.String("drive_id", driveID))

	url := expandURI(rootRelativeContentURITemplate, rootRelativeContentURITemplate, driveID, uri.Path)
	log.Debug("expanded uri", logger.String("uri", url))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return nil, mapError(err, uri.Path)
	}

	builder := drives.NewItemRootContentRequestBuilder(url, adapter)

	log.Debug("sending get content request")
	content, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to get content", logger.Error(err))
		return nil, mapError(err, uri.Path)
	}

	log.Info("retrieved content stream")
	return io.NopCloser(strings.NewReader(string(content))), nil
}

// WriteFile creates or updates a file in OneDrive with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, uri *fs.URI, r io.Reader, opts fs.WriteOptions) (fs.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	driveID := resolveDriveID(ctx, uri, p.driveResolver)
	log.Debug("resolved drive id", logger.String("drive_id", driveID))

	if uri.Path == "" || uri.Path == "/" {
		return fs.Item{}, &fs.Error{Kind: fs.ErrInvalidRequest, Path: uri.Path}
	}

	// Use resumable upload for large files
	if opts.Size > uploadThreshold {
		return writeLargeFile(ctx, p, driveID, uri.Path, r, opts)
	}

	log.Debug("reading input stream")
	data, err := io.ReadAll(r)
	if err != nil {
		return fs.Item{}, mapError(err, uri.Path)
	}

	// If it turns out the data is larger than threshold after reading, use large file upload
	if int64(len(data)) > uploadThreshold {
		return writeLargeFile(ctx, p, driveID, uri.Path, io.MultiReader(strings.NewReader(string(data)), r), opts)
	}

	url := expandURI("", rootRelativeContentURITemplate, driveID, uri.Path)
	log.Debug("expanded uri", logger.String("uri", url))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return fs.Item{}, mapError(err, uri.Path)
	}

	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}
	if opts.IfMatch != "" {
		config.Headers.Add("If-Match", opts.IfMatch)
	}

	builder := drives.NewItemRootContentRequestBuilder(url, adapter)
	log.Debug("sending upload content request")
	it, err := builder.Put(ctx, data, config)
	if err != nil {
		log.Error("failed to upload content", logger.Error(err))
		return fs.Item{}, mapError(err, uri.Path)
	}

	log.Info("uploaded content", logger.Int("size", *it.GetSize()))
	return mapItemToSharedItem(it, uri.Path), nil
}

// Mkdir creates a new folder in OneDrive at the given structured URI.
func (p *Provider) Mkdir(ctx context.Context, uri *fs.URI) error {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	driveID := resolveDriveID(ctx, uri, p.driveResolver)
	log.Debug("resolved drive id", logger.String("drive_id", driveID))

	parentPath := path.Dir(uri.Path)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}
	name := path.Base(uri.Path)

	requestBody := models.NewDriveItem()
	requestBody.SetName(&name)
	requestBody.SetFolder(models.NewFolder())

	url := expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, parentPath)
	log.Debug("expanded uri", logger.String("uri", url))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return mapError(err, uri.Path)
	}
	builder := drives.NewItemItemsRequestBuilder(url, adapter)

	log.Debug("sending create directory request")
	_, err = builder.Post(ctx, requestBody, nil)
	if err != nil {
		log.Error("failed to create directory", logger.Error(err))
		return mapError(err, uri.Path)
	}

	log.Info("created directory")
	return nil
}

// Remove deletes an item from OneDrive.
func (p *Provider) Remove(ctx context.Context, uri *fs.URI) error {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	driveID := resolveDriveID(ctx, uri, p.driveResolver)
	log.Debug("resolved drive id", logger.String("drive_id", driveID))

	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, uri.Path)
	log.Debug("expanded uri", logger.String("uri", url))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return mapError(err, uri.Path)
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(url, adapter)

	log.Debug("sending delete request")
	err = builder.Delete(ctx, nil)
	if err != nil {
		log.Error("failed to delete item", logger.Error(err))
		return mapError(err, uri.Path)
	}

	log.Info("deleted item")
	return nil
}

// Copy duplicates a file or folder within OneDrive.
func (p *Provider) Copy(ctx context.Context, src, dst *fs.URI, opts fs.CopyOptions) error {
	return &fs.Error{Kind: fs.ErrInternal, Err: func() error { return nil }(), Path: src.Path}
}

// Move relocates or renames a file or folder within OneDrive.
func (p *Provider) Move(ctx context.Context, src, dst *fs.URI) error {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("src", src.Path),
		logger.String("dst", dst.Path),
	)

	srcDriveID := resolveDriveID(ctx, src, p.driveResolver)
	dstDriveID := resolveDriveID(ctx, dst, p.driveResolver)

	if srcDriveID != dstDriveID {
		return &fs.Error{Kind: fs.ErrInvalidRequest, Err: func() error { return nil }(), Path: src.Path}
	}
	log.Debug("resolved drive id", logger.String("drive_id", srcDriveID))

	newName := path.Base(dst.Path)
	parentPath := path.Dir(dst.Path)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}

	log.Debug("retrieving parent metadata for move")
	parentURI := *dst // shallow copy
	parentURI.Path = parentPath
	parent, err := p.Get(ctx, &parentURI)
	if err != nil {
		return fmt.Errorf("failed to get destination parent metadata: %w", err)
	}

	requestBody := models.NewDriveItem()
	requestBody.SetName(&newName)
	ref := models.NewItemReference()
	id := parent.ID
	ref.SetId(&id)
	requestBody.SetParentReference(ref)

	url := expandURI(rootURITemplate, rootRelativeURITemplate, srcDriveID, src.Path)
	log.Debug("expanded uri", logger.String("uri", url))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return mapError(err, src.Path)
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(url, adapter)

	log.Debug("sending move (patch) request")
	_, err = builder.Patch(ctx, requestBody, nil)
	if err != nil {
		log.Error("failed to move item", logger.Error(err))
		return mapError(err, src.Path)
	}

	log.Info("moved item")
	return nil
}

// Touch updates the timestamp of an existing file or creates an empty one.
func (p *Provider) Touch(ctx context.Context, uri *fs.URI) (fs.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.Path),
	)

	log.Debug("touching item")
	item, err := p.Get(ctx, uri)
	if err == nil {
		log.Info("item exists, touch completed (metadata refreshed)")
		return item, nil
	}

	log.Info("item does not exist, creating empty file")
	return p.WriteFile(ctx, uri, strings.NewReader(""), fs.WriteOptions{})
}
