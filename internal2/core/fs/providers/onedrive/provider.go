package onedrive

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/core/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

const (
	baseURL = "https://graph.microsoft.com/v1.0"
	// URI Templates from old implementation
	rootURITemplate                 = "{+baseurl}/drives/{drive_id}/root"
	rootRelativeURITemplate         = "{+baseurl}/drives/{drive_id}/root:{path}:"
	rootChildrenURITemplate         = "{+baseurl}/drives/{drive_id}/root/children"
	rootRelativeChildrenURITemplate = "{+baseurl}/drives/{drive_id}/root:{path}:/children"
	rootRelativeContentURITemplate  = "{+baseurl}/drives/{drive_id}/root:{path}:/content"
)

// Provider implements the filesystem Service interface for Microsoft OneDrive.
type Provider struct {
	graph    *microsoft.GraphProvider
	state    state.Service
	driveSvc drive.Service
	log      logger.Logger
}

// NewProvider creates a new instance of the OneDrive filesystem provider.
func NewProvider(graph *microsoft.GraphProvider, state state.Service, driveSvc drive.Service, log logger.Logger) *Provider {
	return &Provider{
		graph:    graph,
		state:    state,
		driveSvc: driveSvc,
		log:      log,
	}
}

func (p *Provider) resolveDrive(ctx context.Context, itemPath string) (string, string, error) {
	// If path is "alias:path"
	if !strings.HasPrefix(itemPath, "/") && strings.Contains(itemPath, ":") {
		alias, cleanPath, _ := strings.Cut(itemPath, ":")
		driveID, err := p.state.GetDriveAlias(alias)
		if err == nil {
			return driveID, cleanPath, nil
		}
	}

	// Default to active drive or "me"
	driveID, err := p.state.Get(state.KeyDrive)
	if err != nil || driveID == "" {
		// Fallback to primary drive
		return "me", itemPath, nil
	}

	return driveID, itemPath, nil
}

func (p *Provider) expandURI(rootTemplate, relativeTemplate, driveID, itemPath string) string {
	normalized := p.normalizePath(itemPath)
	urlTemplate := rootTemplate
	subs := make(stduritemplate.Substitutions)
	subs["baseurl"] = baseURL
	subs["drive_id"] = driveID

	if normalized != "" {
		urlTemplate = relativeTemplate
		subs["path"] = normalized
	}

	uri, _ := stduritemplate.Expand(urlTemplate, subs)
	return uri
}

func (p *Provider) normalizePath(pth string) string {
	if pth == "" || pth == "/" || pth == "." {
		return ""
	}
	// OneDrive relative paths in URI templates usually start with /
	return path.Clean("/" + pth)
}

// Get retrieves metadata for a single item by its OneDrive path.
func (p *Provider) Get(ctx context.Context, itemPath string) (shared.Item, error) {
	p.log.Debug("onedrive.Get", logger.String("path", itemPath))

	driveID, cleanPath, err := p.resolveDrive(ctx, itemPath)
	if err != nil {
		return shared.Item{}, err
	}

	uri := p.expandURI(rootURITemplate, rootRelativeURITemplate, driveID, cleanPath)
	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return shared.Item{}, err
	}
	builder := drives.NewItemRootRequestBuilder(uri, adapter)

	it, err := builder.Get(ctx, nil)
	if err != nil {
		return shared.Item{}, fmt.Errorf("failed to get item metadata: %w", err)
	}

	return p.mapItemToSharedItem(it, itemPath), nil
}

// List enumerates the contents of a directory in OneDrive.
func (p *Provider) List(ctx context.Context, itemPath string, opts shared.ListOptions) ([]shared.Item, error) {
	p.log.Debug("onedrive.List", logger.String("path", itemPath))

	driveID, cleanPath, err := p.resolveDrive(ctx, itemPath)
	if err != nil {
		return nil, err
	}

	uri := p.expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, cleanPath)
	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return nil, err
	}
	builder := drives.NewItemItemsRequestBuilder(uri, adapter)

	resp, err := builder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}

	var items []shared.Item
	for _, it := range resp.GetValue() {
		childPath := path.Join(itemPath, *it.GetName())
		items = append(items, p.mapItemToSharedItem(it, childPath))

		if opts.Recursive && it.GetFolder() != nil {
			children, err := p.List(ctx, childPath, opts)
			if err == nil {
				items = append(items, children...)
			}
		}
	}

	return items, nil
}

// ReadFile opens a read stream for a file's content in OneDrive.
func (p *Provider) ReadFile(ctx context.Context, itemPath string, opts shared.ReadOptions) (io.ReadCloser, error) {
	p.log.Debug("onedrive.ReadFile", logger.String("path", itemPath))

	driveID, cleanPath, err := p.resolveDrive(ctx, itemPath)
	if err != nil {
		return nil, err
	}

	// For content, we use a slightly different template
	uri := p.expandURI("", rootRelativeContentURITemplate, driveID, cleanPath)
	// If path is empty (root), rootRelativeContentURITemplate won't work well,
	// but usually we don't cat the root.
	if cleanPath == "" || cleanPath == "/" {
		return nil, fmt.Errorf("cannot read content of root directory")
	}

	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return nil, err
	}
	// NewItemRootContentRequestBuilder works for expanded URIs
	builder := drives.NewItemRootContentRequestBuilder(uri, adapter)

	content, err := builder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download content: %w", err)
	}

	return io.NopCloser(strings.NewReader(string(content))), nil
}

// Stat returns metadata for a OneDrive item.
func (p *Provider) Stat(ctx context.Context, itemPath string) (shared.Item, error) {
	return p.Get(ctx, itemPath)
}

// WriteFile creates or updates a file in OneDrive with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, itemPath string, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	p.log.Debug("onedrive.WriteFile", logger.String("path", itemPath))

	driveID, cleanPath, err := p.resolveDrive(ctx, itemPath)
	if err != nil {
		return shared.Item{}, err
	}

	if cleanPath == "" || cleanPath == "/" {
		return shared.Item{}, fmt.Errorf("cannot write to root directory")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return shared.Item{}, err
	}

	uri := p.expandURI("", rootRelativeContentURITemplate, driveID, cleanPath)
	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return shared.Item{}, err
	}

	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}
	if opts.IfMatch != "" {
		config.Headers.Add("If-Match", opts.IfMatch)
	}

	builder := drives.NewItemRootContentRequestBuilder(uri, adapter)
	it, err := builder.Put(ctx, data, config)
	if err != nil {
		return shared.Item{}, fmt.Errorf("failed to upload content: %w", err)
	}

	return p.mapItemToSharedItem(it, itemPath), nil
}

// Mkdir creates a new folder in OneDrive at the given path.
func (p *Provider) Mkdir(ctx context.Context, itemPath string) error {
	p.log.Debug("onedrive.Mkdir", logger.String("path", itemPath))

	driveID, cleanPath, err := p.resolveDrive(ctx, itemPath)
	if err != nil {
		return err
	}

	parentPath := path.Dir(cleanPath)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}
	name := path.Base(cleanPath)

	requestBody := models.NewDriveItem()
	requestBody.SetName(&name)
	requestBody.SetFolder(models.NewFolder())

	uri := p.expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, parentPath)
	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return err
	}
	builder := drives.NewItemItemsRequestBuilder(uri, adapter)

	_, err = builder.Post(ctx, requestBody, nil)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// Remove deletes an item from OneDrive.
func (p *Provider) Remove(ctx context.Context, itemPath string) error {
	p.log.Debug("onedrive.Remove", logger.String("path", itemPath))

	driveID, cleanPath, err := p.resolveDrive(ctx, itemPath)
	if err != nil {
		return err
	}

	uri := p.expandURI(rootURITemplate, rootRelativeURITemplate, driveID, cleanPath)
	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return err
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(uri, adapter)

	err = builder.Delete(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// Copy duplicates a file or folder within OneDrive.
func (p *Provider) Copy(ctx context.Context, src, dst string, opts shared.CopyOptions) error {
	return fmt.Errorf("onedrive internal copy not implemented - use manager for cross-provider/fallback copy")
}

// Move relocates or renames a file or folder within OneDrive.
func (p *Provider) Move(ctx context.Context, src, dst string) error {
	p.log.Debug("onedrive.Move", logger.String("src", src), logger.String("dst", dst))

	srcDriveID, cleanSrc, err := p.resolveDrive(ctx, src)
	if err != nil {
		return err
	}

	dstDriveID, cleanDst, err := p.resolveDrive(ctx, dst)
	if err != nil {
		return err
	}

	if srcDriveID != dstDriveID {
		return fmt.Errorf("cross-drive move not supported via internal Move - use manager")
	}

	newName := path.Base(cleanDst)
	parentPath := path.Dir(cleanDst)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}

	// Need parent ID for move
	parent, err := p.Get(ctx, parentPath)
	if err != nil {
		return fmt.Errorf("failed to get destination parent metadata: %w", err)
	}

	requestBody := models.NewDriveItem()
	requestBody.SetName(&newName)
	ref := models.NewItemReference()
	id := parent.ID
	ref.SetId(&id)
	requestBody.SetParentReference(ref)

	uri := p.expandURI(rootURITemplate, rootRelativeURITemplate, srcDriveID, cleanSrc)
	adapter, err := p.graph.Adapter(ctx)
	if err != nil {
		return err
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(uri, adapter)

	_, err = builder.Patch(ctx, requestBody, nil)
	if err != nil {
		return fmt.Errorf("failed to patch item for move: %w", err)
	}

	return nil
}

// Touch updates the timestamp of an existing file or creates an empty one.
func (p *Provider) Touch(ctx context.Context, itemPath string) (shared.Item, error) {
	p.log.Debug("onedrive.Touch", logger.String("path", itemPath))

	item, err := p.Get(ctx, itemPath)
	if err == nil {
		return item, nil
	}

	return p.WriteFile(ctx, itemPath, strings.NewReader(""), shared.WriteOptions{})
}

func (p *Provider) mapItemToSharedItem(it models.DriveItemable, itemPath string) shared.Item {
	if it == nil {
		return shared.Item{}
	}

	id := ""
	if it.GetId() != nil {
		id = *it.GetId()
	}

	name := ""
	if it.GetName() != nil {
		name = *it.GetName()
	}

	size := int64(0)
	if it.GetSize() != nil {
		size = *it.GetSize()
	}

	itemType := shared.TypeFolder
	if it.GetFile() != nil {
		itemType = shared.TypeFile
	}

	modifiedAt := it.GetLastModifiedDateTime()
	etag := ""
	if it.GetETag() != nil {
		etag = *it.GetETag()
	}

	var mTime time.Time
	if modifiedAt != nil {
		mTime = *modifiedAt
	}

	return shared.Item{
		ID:         id,
		Name:       name,
		Path:       itemPath,
		Type:       itemType,
		Size:       size,
		ModifiedAt: mTime,
		ETag:       etag,
	}
}
