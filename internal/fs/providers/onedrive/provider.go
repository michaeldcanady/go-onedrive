package onedrive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"net/http"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	platform "github.com/michaeldcanady/go-onedrive/internal/identity/providers/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

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
	platform platform.PlatformProvider
	state    state.Service
	alias    alias.Service
	driveSvc drive.Service
	log      logger.Logger
}

// NewProvider creates a new instance of the OneDrive filesystem provider.
func NewProvider(p platform.PlatformProvider, state state.Service, driveSvc drive.Service, log logger.Logger) *Provider {
	return &Provider{
		platform: p,
		state:    state,
		driveSvc: driveSvc,
		log:      log,
	}
}

func (p *Provider) Name() string {
	return providerName
}

// mapError translates errors from the OneDrive SDK into domain-specific errors with appropriate kinds for better error handling in the manager and CLI.
func (p *Provider) mapError(err error, itemPath string) error {
	if err == nil {
		return nil
	}

	code := coreerrors.CodeInternal
	safeMsg := "an unexpected OneDrive API error occurred"
	hint := "Check your internet connection and authentication status."

	var apiErr *abstractions.ApiError
	if errors.As(err, &apiErr) {
		switch apiErr.ResponseStatusCode {
		case 401:
			code = coreerrors.CodeUnauthorized
			safeMsg = "unauthorized: please log in again"
			hint = "Use 'odc auth login' to re-authenticate."
		case 403:
			code = coreerrors.CodeForbidden
			safeMsg = "forbidden: you do not have permission to access this resource"
		case 404:
			code = coreerrors.CodeNotFound
			safeMsg = "resource not found in OneDrive"
			hint = ""
		case 409:
			code = coreerrors.CodeConflict
			safeMsg = "resource conflict: the item already exists or is locked"
		case 412:
			code = coreerrors.CodePrecondition
			safeMsg = "precondition failed (e.g., ETag mismatch)"
		case 429:
			code = coreerrors.CodeTransient
			safeMsg = "too many requests: OneDrive is throttling your requests"
			hint = "Please wait a moment and try again."
		case 503, 504:
			code = coreerrors.CodeTransient
			safeMsg = "OneDrive service is temporarily unavailable"
			hint = "Please wait a moment and try again."
		}
	}

	appErr := coreerrors.NewAppError(code, err, safeMsg, hint)
	if itemPath != "" {
		appErr.WithContext(coreerrors.KeyPath, itemPath)
	}

	driveID, _ := p.resolveDrive(itemPath)
	if driveID != "" {
		appErr.WithContext(coreerrors.KeyDriveID, driveID)
	}

	return appErr
}

// resolveDrive determines the drive ID and relative path for a given item path, handling aliases and defaults.
func (p *Provider) resolveDrive(itemPath string) (string, string) {
	// If path is "alias:path"
	if !strings.HasPrefix(itemPath, "/") && strings.Contains(itemPath, ":") {
		alias, cleanPath, _ := strings.Cut(itemPath, ":")
		driveID, err := p.alias.GetDriveIDByAlias(alias)
		if err == nil {
			return driveID, cleanPath
		}
	}

	// Default to active drive or "me"
	driveID, err := p.state.Get(state.KeyDrive)
	if err != nil || driveID == "" {
		// Fallback to primary drive
		return "me", itemPath
	}

	return driveID, itemPath
}

// expandURI constructs the full API endpoint URI based on the provided templates and parameters.
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

// normalizePath cleans up the item path and ensures it is in the correct format for OneDrive API calls.
func (p *Provider) normalizePath(pth string) string {
	if pth == "" || pth == "/" || pth == "." {
		return ""
	}
	// OneDrive relative paths in URI templates usually start with /
	return path.Clean("/" + pth)
}

// Get retrieves metadata for a single item by its OneDrive path.
func (p *Provider) Get(ctx context.Context, itemPath string) (shared.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
	)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(itemPath)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID), logger.String("clean_path", cleanPath))

	uri := p.expandURI(rootURITemplate, rootRelativeURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", uri))

	log.Debug("retrieving adapter")
	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return shared.Item{}, p.mapError(err, itemPath)
	}

	builder := drives.NewItemRootRequestBuilder(uri, adapter)

	log.Debug("sending get item metadata request")
	it, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to get item metadata", logger.Error(err))
		return shared.Item{}, p.mapError(err, itemPath)
	}

	log.Debug("retrieved item metadata", logger.String("item_id", *it.GetId()))
	return p.mapItemToSharedItem(it, itemPath), nil
}

// List enumerates the contents of a directory in OneDrive.
func (p *Provider) List(ctx context.Context, itemPath string, opts shared.ListOptions) ([]shared.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
		logger.Bool("recursive", opts.Recursive),
	)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(itemPath)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	uri := p.expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", uri))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return nil, p.mapError(err, itemPath)
	}
	builder := drives.NewItemItemsRequestBuilder(uri, adapter)

	log.Debug("sending list children request")
	resp, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to list children", logger.Error(err))
		return nil, p.mapError(err, itemPath)
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

	log.Debug("listed items", logger.Int("count", len(items)))
	return items, nil
}

// ReadFile opens a read stream for a file's content in OneDrive.
func (p *Provider) ReadFile(ctx context.Context, itemPath string, opts shared.ReadOptions) (io.ReadCloser, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
	)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(itemPath)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	if cleanPath == "" || cleanPath == "/" {
		return nil, coreerrors.NewInvalidInput(nil, "invalid item path for reading content", "The root folder cannot be read as a file.").WithContext(coreerrors.KeyPath, itemPath)
	}

	uri := p.expandURI("", rootRelativeContentURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", uri))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return nil, p.mapError(err, itemPath)
	}
	builder := drives.NewItemRootContentRequestBuilder(uri, adapter)

	log.Debug("sending download content request")
	content, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to download content", logger.Error(err))
		return nil, p.mapError(err, itemPath)
	}

	log.Debug("downloaded content", logger.Int("size", len(content)))
	return io.NopCloser(strings.NewReader(string(content))), nil
}

// Stat returns metadata for a OneDrive item.
func (p *Provider) Stat(ctx context.Context, itemPath string) (shared.Item, error) {
	return p.Get(ctx, itemPath)
}

// WriteFile creates or updates a file in OneDrive with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, itemPath string, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
	)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(itemPath)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	if cleanPath == "" || cleanPath == "/" {
		return shared.Item{}, coreerrors.NewInvalidInput(nil, "invalid item path for writing content", "The root folder cannot be overwritten as a file.").WithContext(coreerrors.KeyPath, itemPath)
	}

	// Use resumable upload for large files
	if opts.Size > uploadThreshold {
		return p.writeLargeFile(ctx, driveID, cleanPath, itemPath, r, opts)
	}

	log.Debug("reading input stream")
	data, err := io.ReadAll(r)
	if err != nil {
		return shared.Item{}, p.mapError(err, itemPath)
	}

	// If it turns out the data is larger than threshold after reading, use large file upload
	if int64(len(data)) > uploadThreshold {
		return p.writeLargeFile(ctx, driveID, cleanPath, itemPath, io.MultiReader(strings.NewReader(string(data)), r), opts)
	}

	uri := p.expandURI("", rootRelativeContentURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", uri))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return shared.Item{}, p.mapError(err, itemPath)
	}

	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}
	if opts.IfMatch != "" {
		config.Headers.Add("If-Match", opts.IfMatch)
	}

	builder := drives.NewItemRootContentRequestBuilder(uri, adapter)
	log.Debug("sending upload content request")
	it, err := builder.Put(ctx, data, config)
	if err != nil {
		log.Error("failed to upload content", logger.Error(err))
		return shared.Item{}, p.mapError(err, itemPath)
	}

	log.Info("uploaded content", logger.Int("size", *it.GetSize()))
	return p.mapItemToSharedItem(it, itemPath), nil
}

// writeLargeFile handles uploading files larger than the threshold using OneDrive's resumable upload session API.
func (p *Provider) writeLargeFile(ctx context.Context, driveID, cleanPath, itemPath string, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("method", "writeLargeFile"),
		logger.String("path", itemPath),
	)

	uri := p.expandURI("", rootRelativeCreateSessionURITemplate, driveID, cleanPath)
	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		return shared.Item{}, p.mapError(err, itemPath)
	}

	// 1. Create Upload Session
	sessionReq := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	itemProps := models.NewDriveItemUploadableProperties()
	name := path.Base(cleanPath)
	itemProps.SetName(&name)
	sessionReq.SetItem(itemProps)

	builder := drives.NewItemItemsItemCreateUploadSessionRequestBuilder(uri, adapter)
	session, err := builder.Post(ctx, sessionReq, nil)
	if err != nil {
		return shared.Item{}, p.mapError(err, itemPath)
	}

	uploadURL := *session.GetUploadUrl()
	log.Debug("created upload session", logger.String("url", uploadURL))

	// 2. Upload Chunks
	totalSize := opts.Size
	var uploaded int64

	buffer := make([]byte, uploadChunkSize)
	for {
		n, err := r.Read(buffer)
		if n > 0 {
			chunk := buffer[:n]
			req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, strings.NewReader(string(chunk)))
			if err != nil {
				return shared.Item{}, p.mapError(err, itemPath)
			}

			contentRange := fmt.Sprintf("bytes %d-%d/", uploaded, uploaded+int64(n)-1)
			if totalSize > 0 {
				contentRange += fmt.Sprintf("%d", totalSize)
			} else {
				contentRange += "*"
			}
			req.Header.Set("Content-Range", contentRange)
			req.Header.Set("Content-Length", fmt.Sprintf("%d", n))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return shared.Item{}, p.mapError(err, itemPath)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				return shared.Item{}, p.mapError(fmt.Errorf("chunk upload failed with status %d", resp.StatusCode), itemPath)
			}

			uploaded += int64(n)

			if resp.StatusCode == 201 || resp.StatusCode == 200 {
				// Final chunk uploaded, response contains the DriveItem (potentially)
				// For simplicity, we Stat the item to get the final metadata.
				return p.Get(ctx, itemPath)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return shared.Item{}, p.mapError(err, itemPath)
		}
	}

	// If we got here, we might have finished without a 200/201 (e.g. totalSize was unknown)
	return p.Get(ctx, itemPath)
}

// Mkdir creates a new folder in OneDrive at the given path.
func (p *Provider) Mkdir(ctx context.Context, itemPath string) error {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
	)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(itemPath)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	parentPath := path.Dir(cleanPath)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}
	name := path.Base(cleanPath)

	requestBody := models.NewDriveItem()
	requestBody.SetName(&name)
	requestBody.SetFolder(models.NewFolder())

	uri := p.expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, parentPath)
	log.Debug("expanded uri", logger.String("uri", uri))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return p.mapError(err, itemPath)
	}
	builder := drives.NewItemItemsRequestBuilder(uri, adapter)

	log.Debug("sending create directory request")
	_, err = builder.Post(ctx, requestBody, nil)
	if err != nil {
		log.Error("failed to create directory", logger.Error(err))
		return p.mapError(err, itemPath)
	}

	log.Info("created directory")
	return nil
}

// Remove deletes an item from OneDrive.
func (p *Provider) Remove(ctx context.Context, itemPath string) error {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
	)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(itemPath)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	uri := p.expandURI(rootURITemplate, rootRelativeURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", uri))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return p.mapError(err, itemPath)
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(uri, adapter)

	log.Debug("sending delete request")
	err = builder.Delete(ctx, nil)
	if err != nil {
		log.Error("failed to delete item", logger.Error(err))
		return p.mapError(err, itemPath)
	}

	log.Info("deleted item")
	return nil
}

// Copy duplicates a file or folder within OneDrive.
func (p *Provider) Copy(ctx context.Context, src, dst string, opts shared.CopyOptions) error {
	return coreerrors.NewInternal(fmt.Errorf("onedrive internal copy not implemented - use manager for cross-provider/fallback copy"), "copy operation not supported internally", "The FileSystemManager should handle this cross-provider or via fallback.").WithContext(coreerrors.KeyPath, src)
}

// Move relocates or renames a file or folder within OneDrive.
func (p *Provider) Move(ctx context.Context, src, dst string) error {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("src", src),
		logger.String("dst", dst),
	)

	log.Debug("resolving source and destination drives")
	srcDriveID, cleanSrc := p.resolveDrive(src)

	dstDriveID, cleanDst := p.resolveDrive(dst)

	if srcDriveID != dstDriveID {
		return coreerrors.NewInvalidInput(nil, "cross-drive move not supported via internal Move", "The FileSystemManager should handle cross-drive moves via Copy and Remove.").WithContext(coreerrors.KeyPath, src)
	}
	log.Debug("resolved drive id", logger.String("drive_id", srcDriveID))

	newName := path.Base(cleanDst)
	parentPath := path.Dir(cleanDst)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}

	log.Debug("retrieving parent metadata for move")
	parent, err := p.Get(ctx, parentPath)
	if err != nil {
		return coreerrors.NewAppError(coreerrors.CodeInternal, err, "failed to get destination parent metadata", "").WithContext(coreerrors.KeyPath, parentPath)
	}

	requestBody := models.NewDriveItem()
	requestBody.SetName(&newName)
	ref := models.NewItemReference()
	id := parent.ID
	ref.SetId(&id)
	requestBody.SetParentReference(ref)

	uri := p.expandURI(rootURITemplate, rootRelativeURITemplate, srcDriveID, cleanSrc)
	log.Debug("expanded uri", logger.String("uri", uri))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return p.mapError(err, src)
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(uri, adapter)

	log.Debug("sending move (patch) request")
	_, err = builder.Patch(ctx, requestBody, nil)
	if err != nil {
		log.Error("failed to move item", logger.Error(err))
		return p.mapError(err, src)
	}

	log.Info("moved item")
	return nil
}

// Touch updates the timestamp of an existing file or creates an empty one.
func (p *Provider) Touch(ctx context.Context, itemPath string) (shared.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", itemPath),
	)

	log.Debug("touching item")
	item, err := p.Get(ctx, itemPath)
	if err == nil {
		log.Info("item exists, touch completed (metadata refreshed)")
		return item, nil
	}

	log.Info("item does not exist, creating empty file")
	return p.WriteFile(ctx, itemPath, strings.NewReader(""), shared.WriteOptions{})
}

// mapItemToSharedItem converts a OneDrive DriveItem to the shared Item format.
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
