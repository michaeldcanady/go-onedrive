package onedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	baseURL                              = "https://graph.microsoft.com/v1.0"
	rootURITemplate                      = "{+baseurl}/drives/{drive_id}/root"
	rootRelativeURITemplate              = "{+baseurl}/drives/{drive_id}/root:{path}:"
	rootChildrenURITemplate              = "{+baseurl}/drives/{drive_id}/root/children"
	rootRelativeChildrenURITemplate      = "{+baseurl}/drives/{drive_id}/root:{path}:/children"
	rootRelativeContentURITemplate       = "{+baseurl}/drives/{drive_id}/root:{path}:/content"
	rootRelativeCreateSessionURITemplate = "{+baseurl}/drives/{drive_id}/root:{path}:/createUploadSession"
	uploadThreshold                      = 4 * 1024 * 1024
)

const (
	driveIDOptionKey = "drive_id"
)

type Backend struct {
	driveID string
}

func NewBackend(opts map[string]string) *Backend {
	driveID := opts[driveIDOptionKey]

	return &Backend{
		driveID: driveID,
	}
}

func (b *Backend) Name() string {
	return "onedrive"
}

func (b *Backend) createAdapter(ctx context.Context, rawToken string) (abstractions.RequestAdapter, error) {
	var token identity.AccessToken

	if err := json.Unmarshal([]byte(rawToken), &token); err != nil {
		return nil, err
	}

	cred := microsoft.NewStaticTokenCredential(token)

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{
		"Files.ReadWrite.All",
		"User.Read",
		"offline_access",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication provider: %w", err)
	}

	return client.RequestAdapter, nil
}

func (b *Backend) Stat(ctx context.Context, token, driveID, path string) (fs.Item, error) {
	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, path)
	adapter, err := b.createAdapter(ctx, token)
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

func (b *Backend) List(ctx context.Context, token, driveID, path string) ([]fs.Item, error) {
	url := expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, b.driveID, path)
	adapter, err := b.createAdapter(ctx, token)
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
		items = append(items, mapItemToSharedItem(it, joinPath(path, name)))
	}
	return items, nil
}

func (b *Backend) Capabilities() fs.Capabilities {
	return fs.Capabilities{CanMove: true, CanCopy: false, CanRecursive: true}
}

func (b *Backend) Open(ctx context.Context, token, driveID, path string) (io.ReadCloser, error) {
	url := expandURI(rootRelativeContentURITemplate, rootRelativeContentURITemplate, driveID, path)
	adapter, err := b.createAdapter(ctx, token)
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

func (b *Backend) Create(ctx context.Context, token, driveID, path string, r io.Reader) (fs.Item, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return fs.Item{}, mapError(err, path)
	}

	url := expandURI("", rootRelativeContentURITemplate, driveID, path)
	adapter, err := b.createAdapter(ctx, token)
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

func (b *Backend) Mkdir(ctx context.Context, token, driveID, itemPath string) error {
	parentPath := path.Dir(itemPath)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}
	name := path.Base(itemPath)

	requestBody := models.NewDriveItem()
	requestBody.SetName(&name)
	requestBody.SetFolder(models.NewFolder())

	url := expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, parentPath)
	adapter, err := b.createAdapter(ctx, token)
	if err != nil {
		return mapError(err, itemPath)
	}
	builder := drives.NewItemItemsRequestBuilder(url, adapter)

	_, err = builder.Post(ctx, requestBody, nil)
	return mapError(err, itemPath)
}

func (b *Backend) Remove(ctx context.Context, token, driveID, path string) error {
	url := expandURI(rootURITemplate, rootRelativeURITemplate, driveID, path)
	adapter, err := b.createAdapter(ctx, token)
	if err != nil {
		return mapError(err, path)
	}

	builder := drives.NewItemItemsDriveItemItemRequestBuilder(url, adapter)
	return builder.Delete(ctx, nil)
}

func (b *Backend) Move(ctx context.Context, token, driveID, src, dst string) error {
	newName := path.Base(dst)
	parentPath := path.Dir(dst)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}

	parent, err := b.Stat(ctx, token, driveID, parentPath)
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
	adapter, err := b.createAdapter(ctx, token)
	if err != nil {
		return mapError(err, src)
	}

	builder := drives.NewItemItemsDriveItemItemRequestBuilder(url, adapter)
	_, err = builder.Patch(ctx, requestBody, nil)
	return mapError(err, src)
}

func (b *Backend) Copy(ctx context.Context, token, driveID, src, dst string) error {
	return fmt.Errorf("copy not supported on OneDrive backend")
}

func joinPath(base, name string) string {
	return path.Join(base, name)
}
