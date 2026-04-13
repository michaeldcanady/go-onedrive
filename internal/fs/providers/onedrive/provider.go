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
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	platform "github.com/michaeldcanady/go-onedrive/internal/identity/providers/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
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
	driveSvc drive.Service
	log      logger.Logger
}

// NewProvider creates a new instance of the OneDrive filesystem provider.
func NewProvider(p platform.PlatformProvider, state state.Service, driveSvc drive.Service, log logger.Logger) *Provider {
	return &Provider{
		platform: p,
		state:    state,
		driveSvc: driveSvc,
		log:      log.With(logger.String("provider", providerName)),
	}
}

func (p *Provider) Name() string {
	return providerName
}

// mapReadError translates errors from the OneDrive SDK into a ReadError wrapping domain-specific errors.
func (p *Provider) mapReadError(err error, itemPath *shared.URI) error {
	if err == nil {
		return nil
	}

	wrapped := p.mapToDomainError(err, itemPath)
	return NewReadError(fs.ToManagerPath(itemPath), wrapped)
}

// mapWriteError translates errors from the OneDrive SDK into a WriteError wrapping domain-specific errors.
func (p *Provider) mapWriteError(err error, itemPath *shared.URI) error {
	if err == nil {
		return nil
	}

	wrapped := p.mapToDomainError(err, itemPath)
	return NewWriteError(fs.ToManagerPath(itemPath), wrapped)
}

// mapGenericError translates errors from the OneDrive SDK into domain-specific errors wrapped in an AppError.
func (p *Provider) mapGenericError(err error, itemPath *shared.URI) error {
	if err == nil {
		return nil
	}

	domainErr := p.mapToDomainError(err, itemPath)

	// If mapToDomainError already returned one of our custom error types, return it directly.
	// This supports the "multi-modal branching patterns" desired by the user.
	if domainErr != err {
		return domainErr
	}

	code := coreerrors.CodeInternal
	safeMsg := "an unexpected OneDrive API error occurred"
	hint := "Check your internet connection and authentication status."

	appErr := coreerrors.NewAppError(code, domainErr, safeMsg, hint)
	if itemPath != nil {
		appErr.WithContext(coreerrors.KeyPath, fs.ToManagerPath(itemPath))
	}

	driveID, _ := p.resolveDrive(itemPath)
	if driveID != "" {
		appErr.WithContext(coreerrors.KeyDriveID, driveID)
	}

	return appErr
}

// mapToDomainError converts low-level SDK errors into domain-specific error types.
func (p *Provider) mapToDomainError(err error, itemPath *shared.URI) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, microsoft.ErrNotAuthenticated) {
		return coreerrors.NewUnauthorizedError(err)
	}

	// 1. Handle OData specific errors (detailed errors from Graph)
	var odataErr *odataerrors.ODataError
	if errors.As(err, &odataErr) {
		return p.mapStatusCodeToError(odataErr.GetStatusCode(), err, itemPath)
	}

	// 2. Handle broader Kiota API errors
	var apiErr *abstractions.ApiError
	if errors.As(err, &apiErr) {
		return p.mapStatusCodeToError(apiErr.GetStatusCode(), err, itemPath)
	}

	// 3. Handle Azure SDK errors (e.g., token acquisition failures)
	// Some auth failures might not be 401s yet but are still auth-related.
	errMsg := err.Error()
	if strings.Contains(errMsg, "AuthenticationFailed") ||
		strings.Contains(errMsg, "AADSTS") ||
		strings.Contains(errMsg, "token") ||
		strings.Contains(errMsg, "expired") {
		return coreerrors.NewUnauthorizedError(err)
	}

	return err
}

// mapStatusCodeToError maps an HTTP status code to a domain-specific error.
func (p *Provider) mapStatusCodeToError(statusCode int, err error, itemPath *shared.URI) error {
	switch statusCode {
	case 400:
		return coreerrors.NewBadRequestError(err)
	case 401:
		return coreerrors.NewUnauthorizedError(err)
	case 403:
		return coreerrors.NewForbiddenError(fs.ToManagerPath(itemPath), err)
	case 404:
		return coreerrors.NewNotFoundError(fs.ToManagerPath(itemPath), err)
	case 409:
		return coreerrors.NewConflictError(fs.ToManagerPath(itemPath), err)
	case 410:
		return NewGoneError(fs.ToManagerPath(itemPath), err)
	case 412:
		return coreerrors.NewPreconditionFailedError(fs.ToManagerPath(itemPath), err)
	case 423:
		return NewLockedError(fs.ToManagerPath(itemPath), err)
	case 429:
		return coreerrors.NewTransientError("too many requests: OneDrive is throttling your requests", err)
	case 500:
		return coreerrors.NewInternalError(err)
	case 503, 504:
		return coreerrors.NewTransientError("OneDrive service is temporarily unavailable", err)
	case 507:
		return NewInsufficientStorageError(err)
	}
	return err
}

// resolveDrive determines the drive ID and relative path for a given item path, handling aliases and defaults.
func (p *Provider) resolveDrive(uri *shared.URI) (string, string) {
	driveID := uri.DriveRef
	if driveID == "" {
		driveID = "me"
	}
	return driveID, p.normalizePath(uri.Path)
}

// expandURI constructs the full API endpoint URI based on the provided templates and parameters.
func (p *Provider) expandURI(rootTemplate, relativeTemplate, driveID, itemPath string) string {
	normalized := itemPath
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
func (p *Provider) Get(ctx context.Context, uri *shared.URI) (shared.Item, error) {
	log := p.log.WithContext(ctx)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(uri)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID), logger.String("clean_path", cleanPath))

	apiUrl := p.expandURI(rootURITemplate, rootRelativeURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	log.Debug("retrieving adapter")
	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return shared.Item{}, p.mapGenericError(err, uri)
	}

	builder := drives.NewItemRootRequestBuilder(apiUrl, adapter)

	log.Debug("sending get item metadata request")
	it, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to get item metadata", logger.Error(err))
		return shared.Item{}, p.mapGenericError(err, uri)
	}

	log.Debug("retrieved item metadata", logger.String("item_id", *it.GetId()))
	return p.mapItemToSharedItem(it, uri.String()), nil
}

// List enumerates the contents of a directory in OneDrive.
func (p *Provider) List(ctx context.Context, uri *shared.URI, opts shared.ListOptions) ([]shared.Item, error) {
	log := p.log.WithContext(ctx)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(uri)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	apiUrl := p.expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return nil, p.mapGenericError(err, uri)
	}
	builder := drives.NewItemItemsRequestBuilder(apiUrl, adapter)

	log.Debug("sending list children request")
	resp, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to list children", logger.Error(err))
		return nil, p.mapGenericError(err, uri)
	}

	var items []shared.Item
	for _, it := range resp.GetValue() {
		childPath := path.Join(uri.Path, *it.GetName())
		childURI := &shared.URI{
			Provider: uri.Provider,
			DriveRef: driveID,
			Path:     childPath,
		}
		items = append(items, p.mapItemToSharedItem(it, childURI.String()))

		if opts.Recursive && it.GetFolder() != nil {
			children, err := p.List(ctx, childURI, opts)
			if err == nil {
				items = append(items, children...)
			}
		}
	}

	log.Debug("listed items", logger.Int("count", len(items)))
	return items, nil
}

// ReadFile opens a read stream for a file's content in OneDrive.
func (p *Provider) ReadFile(ctx context.Context, uri *shared.URI, opts shared.ReadOptions) (io.ReadCloser, error) {
	log := p.log.WithContext(ctx)

	log.Debug("resolving drive id and path", logger.String("path", uri.String()))
	driveID, cleanPath := p.resolveDrive(uri)
	log.Info("resolved drive id and path", logger.String("drive_id", driveID), logger.String("path", cleanPath))

	if cleanPath == "" || strings.HasSuffix(cleanPath, "/") {
		return nil, NewBadPathError(uri.String(), "cannot read content from a directory", nil)
	}

	apiUrl := p.expandURI("", rootRelativeContentURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return nil, p.mapReadError(err, uri)
	}
	builder := drives.NewItemRootContentRequestBuilder(apiUrl, adapter)

	log.Debug("sending download content request")
	content, err := builder.Get(ctx, nil)
	if err != nil {
		log.Error("failed to download content", logger.Error(err))
		return nil, p.mapReadError(err, uri)
	}

	log.Debug("downloaded content", logger.Int("size", len(content)))
	return io.NopCloser(strings.NewReader(string(content))), nil
}

// Stat returns metadata for a OneDrive item.
func (p *Provider) Stat(ctx context.Context, uri *shared.URI) (shared.Item, error) {
	return p.Get(ctx, uri)
}

// WriteFile creates or updates a file in OneDrive with the content from the reader.
func (p *Provider) WriteFile(ctx context.Context, uri *shared.URI, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	log := p.log.WithContext(ctx)

	log.Debug("resolving drive id and path", logger.String("path", uri.String()))
	driveID, cleanPath := p.resolveDrive(uri)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	if cleanPath == "" || strings.HasSuffix(cleanPath, "/") {
		return shared.Item{}, NewBadPathError(uri.String(), "cannot write content to a directory", nil)
	}

	// Use resumable upload for large files
	if opts.Size > uploadThreshold {
		return p.writeLargeFile(ctx, driveID, cleanPath, uri, r, opts)
	}

	log.Debug("reading input stream")
	data, err := io.ReadAll(r)
	if err != nil {
		return shared.Item{}, p.mapWriteError(err, uri)
	}

	// If it turns out the data is larger than threshold after reading, use large file upload
	if int64(len(data)) > uploadThreshold {
		log.Info("file size is greater than upload threshold", logger.Int("threshold", uploadThreshold))
		return p.writeLargeFile(ctx, driveID, cleanPath, uri, io.MultiReader(strings.NewReader(string(data)), r), opts)
	}

	return p.writeFile(ctx, driveID, cleanPath, uri, data, opts)
}

func (p *Provider) writeFile(ctx context.Context, driveID, cleanPath string, uri *shared.URI, data []byte, opts shared.WriteOptions) (shared.Item, error) {
	log := p.log.WithContext(ctx)

	apiUrl := p.expandURI("", rootRelativeContentURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return shared.Item{}, p.mapWriteError(err, uri)
	}

	config := &drives.ItemRootContentRequestBuilderPutRequestConfiguration{
		Headers: abstractions.NewRequestHeaders(),
	}
	if opts.IfMatch != "" {
		config.Headers.Add("If-Match", opts.IfMatch)
	}

	builder := drives.NewItemRootContentRequestBuilder(apiUrl, adapter)
	log.Debug("sending upload content request")
	it, err := builder.Put(ctx, data, config)
	if err != nil {
		log.Error("failed to upload content", logger.Error(err))
		return shared.Item{}, p.mapWriteError(err, uri)
	}

	log.Info("uploaded content", logger.Int("size", *it.GetSize()))
	return p.mapItemToSharedItem(it, fs.ToManagerPath(uri)), nil
}

// writeLargeFile handles uploading files larger than the threshold using OneDrive's resumable upload session API.
func (p *Provider) writeLargeFile(ctx context.Context, driveID, cleanPath string, uri *shared.URI, r io.Reader, opts shared.WriteOptions) (shared.Item, error) {
	log := p.log.WithContext(ctx)

	apiUrl := p.expandURI("", rootRelativeCreateSessionURITemplate, driveID, cleanPath)
	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		return shared.Item{}, p.mapWriteError(err, uri)
	}

	// 1. Create Upload Session
	sessionReq := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	itemProps := models.NewDriveItemUploadableProperties()
	name := path.Base(cleanPath)
	itemProps.SetName(&name)
	sessionReq.SetItem(itemProps)

	builder := drives.NewItemItemsItemCreateUploadSessionRequestBuilder(apiUrl, adapter)
	session, err := builder.Post(ctx, sessionReq, nil)
	if err != nil {
		return shared.Item{}, p.mapWriteError(err, uri)
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
				return shared.Item{}, p.mapWriteError(err, uri)
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
				return shared.Item{}, p.mapWriteError(err, uri)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				return shared.Item{}, p.mapWriteError(fmt.Errorf("chunk upload failed with status %d", resp.StatusCode), uri)
			}

			uploaded += int64(n)

			if resp.StatusCode == 201 || resp.StatusCode == 200 {
				// Final chunk uploaded, response contains the DriveItem (potentially)
				// For simplicity, we Stat the item to get the final metadata.
				return p.Get(ctx, uri)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return shared.Item{}, p.mapWriteError(err, uri)
		}
	}

	// If we got here, we might have finished without a 200/201 (e.g. totalSize was unknown)
	return p.Get(ctx, uri)
}

// Mkdir creates a new folder in OneDrive at the given path.
func (p *Provider) Mkdir(ctx context.Context, uri *shared.URI) error {
	log := p.log.WithContext(ctx)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(uri)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	parentPath := path.Dir(cleanPath)
	if parentPath == "." || parentPath == "/" {
		parentPath = ""
	}
	name := path.Base(cleanPath)

	requestBody := models.NewDriveItem()
	requestBody.SetName(&name)
	requestBody.SetFolder(models.NewFolder())

	apiUrl := p.expandURI(rootChildrenURITemplate, rootRelativeChildrenURITemplate, driveID, parentPath)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return p.mapGenericError(err, uri)
	}
	builder := drives.NewItemItemsRequestBuilder(apiUrl, adapter)

	log.Debug("sending create directory request")
	_, err = builder.Post(ctx, requestBody, nil)
	if err != nil {
		log.Error("failed to create directory", logger.Error(err))
		return p.mapGenericError(err, uri)
	}

	log.Info("created directory")
	return nil
}

// Remove deletes an item from OneDrive.
func (p *Provider) Remove(ctx context.Context, uri *shared.URI) error {
	log := p.log.WithContext(ctx)

	log.Debug("resolving drive id and path")
	driveID, cleanPath := p.resolveDrive(uri)
	log.Debug("resolved drive id and path", logger.String("drive_id", driveID))

	apiUrl := p.expandURI(rootURITemplate, rootRelativeURITemplate, driveID, cleanPath)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return p.mapGenericError(err, uri)
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(apiUrl, adapter)

	log.Debug("sending delete request")
	err = builder.Delete(ctx, nil)
	if err != nil {
		log.Error("failed to delete item", logger.Error(err))
		return p.mapGenericError(err, uri)
	}

	log.Info("deleted item")
	return nil
}

// Copy duplicates a file or folder within OneDrive.
func (p *Provider) Copy(ctx context.Context, src, dst *shared.URI, opts shared.CopyOptions) error {
	return NewUnsupportedOperationError(src.String(), "internal copy", coreerrors.NewInternal(fmt.Errorf("onedrive internal copy not implemented - use manager for cross-provider/fallback copy"), "copy operation not supported internally", "The FileSystemManager should handle this cross-provider or via fallback."))
}

// Move relocates or renames a file or folder within OneDrive.
func (p *Provider) Move(ctx context.Context, src, dst *shared.URI) error {
	log := p.log.WithContext(ctx)

	log.Debug("resolving source and destination drives")
	srcDriveID, cleanSrc := p.resolveDrive(src)

	dstDriveID, cleanDst := p.resolveDrive(dst)

	if srcDriveID != dstDriveID {
		return NewUnsupportedOperationError(src.String(), "cross-drive move", coreerrors.NewInvalidInput(nil, "cross-drive move not supported via internal Move", "The FileSystemManager should handle cross-drive moves via Copy and Remove."))
	}
	log.Debug("resolved drive id", logger.String("drive_id", srcDriveID))

	newName := path.Base(cleanDst)
	parentURI := shared.ParentURI(dst)
	parent, err := p.Get(ctx, parentURI)
	if err != nil {
		return p.mapGenericError(err, parentURI)
	}

	requestBody := models.NewDriveItem()
	requestBody.SetName(&newName)
	ref := models.NewItemReference()
	id := parent.ID
	ref.SetId(&id)
	requestBody.SetParentReference(ref)

	apiUrl := p.expandURI(rootURITemplate, rootRelativeURITemplate, srcDriveID, cleanSrc)
	log.Debug("expanded uri", logger.String("uri", apiUrl))

	adapter, err := p.platform.Adapter(ctx)
	if err != nil {
		log.Warn("failed to retrieve adapter", logger.Error(err))
		return p.mapGenericError(err, src)
	}
	builder := drives.NewItemItemsDriveItemItemRequestBuilder(apiUrl, adapter)

	log.Debug("sending move (patch) request")
	_, err = builder.Patch(ctx, requestBody, nil)
	if err != nil {
		log.Error("failed to move item", logger.Error(err))
		return p.mapGenericError(err, src)
	}

	log.Info("moved item")
	return nil
}

// Touch updates the timestamp of an existing file or creates an empty one.
func (p *Provider) Touch(ctx context.Context, uri *shared.URI) (shared.Item, error) {
	log := p.log.WithContext(ctx).With(
		logger.String("provider", providerName),
		logger.String("path", uri.String()),
	)

	log.Debug("touching item")
	item, err := p.Get(ctx, uri)
	if err == nil {
		log.Info("item exists, touch completed (metadata refreshed)")
		return item, nil
	}

	log.Info("item does not exist, creating empty file")
	return p.WriteFile(ctx, uri, strings.NewReader(""), shared.WriteOptions{})
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
