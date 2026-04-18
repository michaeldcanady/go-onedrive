package onedrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	// uploadChunkSize is the size of each chunk in a resumable upload (must be multiple of 320 KiB).
	uploadChunkSize = 320 * 1024 * 10 // 3.2 MiB
)

func writeLargeFile(ctx context.Context, b *Backend, driveID, itemPath string, r io.Reader, opts fs.WriteOptions) (fs.Item, error) {
	log := b.log.WithContext(ctx).With(
		logger.String("method", "writeLargeFile"),
		logger.String("path", itemPath),
	)

	url := expandURI("", rootRelativeCreateSessionURITemplate, driveID, itemPath)
	adapter, err := b.platform.Adapter(ctx)
	if err != nil {
		return fs.Item{}, mapError(err, itemPath)
	}

	// 1. Create Upload Session
	sessionReq := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	itemProps := models.NewDriveItemUploadableProperties()
	name := path.Base(itemPath)
	itemProps.SetName(&name)
	sessionReq.SetItem(itemProps)

	builder := drives.NewItemItemsItemCreateUploadSessionRequestBuilder(url, adapter)
	session, err := builder.Post(ctx, sessionReq, nil)
	if err != nil {
		return fs.Item{}, mapError(err, itemPath)
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
				return fs.Item{}, mapError(err, itemPath)
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
				return fs.Item{}, mapError(err, itemPath)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				return fs.Item{}, mapError(fmt.Errorf("chunk upload failed with status %d", resp.StatusCode), itemPath)
			}

			uploaded += int64(n)

			if resp.StatusCode == 201 || resp.StatusCode == 200 {
				// Final chunk uploaded, response contains the DriveItem (potentially)
				// For simplicity, we Stat the item to get the final metadata.
				return b.Stat(ctx, itemPath)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fs.Item{}, mapError(err, itemPath)
		}
	}

	// If we got here, we might have finished without a 200/201 (e.g. totalSize was unknown)
	return b.Stat(ctx, itemPath)
}
