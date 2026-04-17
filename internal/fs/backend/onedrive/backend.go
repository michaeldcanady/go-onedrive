package onedrive

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	platform "github.com/michaeldcanady/go-onedrive/internal/identity/providers/shared"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	baseURL = "https://graph.microsoft.com/v1.0"
	// URI Templates
	rootURITemplate                      = "{+baseurl}/drives/{drive_id}/root"
	rootRelativeURITemplate              = "{+baseurl}/drives/{drive_id}/root:{path}:"
	rootChildrenURITemplate              = "{+baseurl}/drives/{drive_id}/root/children"
	rootRelativeChildrenURITemplate      = "{+baseurl}/drives/{drive_id}/root:{path}:/children"
	rootRelativeContentURITemplate       = "{+baseurl}/drives/{drive_id}/root:{path}:/content"
	rootRelativeCreateSessionURITemplate = "{+baseurl}/drives/{drive_id}/root:{path}:/createUploadSession"
	// uploadThreshold is the file size at which we switch to resumable uploads (4MB).
	uploadThreshold = 4 * 1024 * 1024
)

// Backend implements the fs.Backend and fs.AdvancedBackend interfaces for Microsoft OneDrive.
type Backend struct {
	platform      platform.PlatformProvider
	driveID       string
	driveResolver fs.DriveResolver
	log           logger.Logger
}

// NewBackend creates a new instance of the OneDrive filesystem backend.
// If driveID is empty, it uses the driveResolver to find the active drive at runtime.
func NewBackend(p platform.PlatformProvider, driveID string, dr fs.DriveResolver, log logger.Logger) *Backend {
	return &Backend{
		platform:      p,
		driveID:       driveID,
		driveResolver: dr,
		log:           log,
	}
}

func (b *Backend) Name() string {
	return "onedrive"
}

func (b *Backend) getDriveID(ctx context.Context) (string, error) {
	if b.driveID != "" {
		return b.driveID, nil
	}
	if b.driveResolver == nil {
		return "", fmt.Errorf("no drive ID or resolver provided")
	}
	return b.driveResolver.GetActiveDriveID(ctx)
}

func (b *Backend) Stat(ctx context.Context, path string) (fs.Item, error) {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, path)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	builder := drives.NewItemRootRequestBuilder(url, adapter)
	it, err := builder.Get(ctx, nil)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	return mapItemToSharedItem(it, path), nil
}

func (b *Backend) List(ctx context.Context, path string) ([]fs.Item, error) {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return nil, mapError(err, path)
	}

	url := expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, path)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return nil, mapError(err, path)
	}

	builder := drives.NewItemItemsRequestBuilder(url, adapter)
	resp, err := builder.Get(ctx, nil)
	if err != nil {
		return nil, mapError(err, path)
	}

	var items []fs.Item
	for _, it := range resp.GetValue() {
		name := ""
		if it.GetName() != nil {
			name = *it.GetName()
		}
		childPath := joinPath(path, name)
		items = append(items, mapItemToSharedItem(it, childPath))
	}

	return items, nil
}

func (b *Backend) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return nil, mapError(err, path)
	}

	url := expandURI(rootRelativeContentURITemplate, rootRelativeContentURITemplate, driveID, path)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return nil, mapError(err, path)
	}

	builder := drives.NewItemRootContentRequestBuilder(url, adapter)
	content, err := builder.Get(ctx, nil)
	if err != nil {
		return nil, mapError(err, path)
	}

	return io.NopCloser(strings.NewReader(string(content))), nil
}

func (b *Backend) Create(ctx context.Context, path string, r io.Reader) (fs.Item, error) {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	// For simplicity, we'll read everything to check size, or just use the reader.
	// In a real scenario we might want to know the size beforehand.
	// For now, let's just use the basic Put if it's small.
	data, err := io.ReadAll(r)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	if int64(len(data)) > uploadThreshold {
		return writeLargeFile(ctx, b, driveID, path, strings.NewReader(string(data)), fs.WriteOptions{Size: int64(len(data))})
	}

	url := expandURI("", rootRelativeContentURITemplate, driveID, path)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	builder := drives.NewItemRootContentRequestBuilder(url, adapter)
	it, err := builder.Put(ctx, data, nil)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	return mapItemToSharedItem(it, path), nil
}

func (b *Backend) Mkdir(ctx context.Context, itemPath string) error {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return mapError(err, itemPath)
	}

	parentPath := path.Dir(itemPath)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}
	name := path.Base(itemPath)

	requestBody := models.NewDriveItem()
	requestBody.SetName(&name)
	requestBody.SetFolder(models.NewFolder())

	url := expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, parentPath)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return mapError(err, itemPath)
	}
	builder := drives.NewItemItemsRequestBuilder(url, adapter)

	_, err = builder.Post(ctx, requestBody, nil)
	return mapError(err, itemPath)
}

func (b *Backend) Remove(ctx context.Context, path string) error {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return mapError(err, path)
	}

	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, path)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return mapError(err, path)
	}

	builder := drives.NewItemItemsDriveItemItemRequestBuilder(url, adapter)
	return builder.Delete(ctx, nil)
}

func (b *Backend) Capabilities() fs.Capabilities {
	return fs.Capabilities{
		CanMove:      true,
		CanCopy:      false, // OneDrive copy is async, complex to implement here
		CanRecursive: true,  // OneDrive API handles recursive operations well if we use their native calls
	}
}

func (b *Backend) Move(ctx context.Context, src, dst string) error {
	driveID, err := b.getDriveID(ctx)
	if err != nil {
		return mapError(err, src)
	}

	newName := path.Base(dst)
	parentPath := path.Dir(dst)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}

	parent, err := b.Stat(ctx, parentPath)
	if err != nil {
		return err
	}

	requestBody := models.NewDriveItem()
	requestBody.SetName(&newName)
	ref := models.NewItemReference()
	id := parent.ID
	ref.SetId(&id)
	requestBody.SetParentReference(ref)

	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, src)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return mapError(err, src)
	}

	builder := drives.NewItemItemsDriveItemItemRequestBuilder(url, adapter)
	_, err = builder.Patch(ctx, requestBody, nil)
	return mapError(err, src)
}

func joinPath(base, name string) string {
	return path.Join(base, name)
}
