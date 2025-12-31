package app

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"path"

	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	rootChildrenURITemplate         = "https://graph.microsoft.com/v1.0/drives/%s/root/children"
	rootRelativeChildrenURITemplate = "https://graph.microsoft.com/v1.0/drives/%s/root:%s:/children"
)

type DriveService struct {
	graph Clienter
}

func NewDriveService(graph Clienter) *DriveService {
	return &DriveService{graph: graph}
}

// normalizePath ensures paths like "Documents", "/Documents", "Documents/" all become "/Documents"
func normalizePath(p string) string {
	if p == "" || p == "/" || p == "." {
		return ""
	}
	p = path.Clean("/" + p)
	return p
}

// getUserDriveID retrieves the user's default drive ID
func (s *DriveService) getUserDriveID(ctx context.Context) (string, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		return "", errors.Join(errors.New("unable to instantiate client"), err)
	}

	drive, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user drive: %w", err)
	}
	id := drive.GetId()
	if id == nil || *id == "" {
		return "", fmt.Errorf("user drive ID is empty")
	}
	return *id, nil
}

func (s *DriveService) getChildren(ctx context.Context, folderPath string) (models.DriveItemCollectionResponseable, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		return nil, err
	}

	// Normalize path
	folderPath = normalizePath(folderPath)

	// Get drive ID
	driveID, err := s.getUserDriveID(ctx)
	if err != nil {
		return nil, err
	}

	// Build the raw URL for the children request
	// Example: https://graph.microsoft.com/v1.0/drives/{id}/root:/Documents:/children
	rawURL := ""
	if folderPath != "" {
		rawURL = fmt.Sprintf(rootRelativeChildrenURITemplate, driveID, folderPath)
	} else {
		rawURL = fmt.Sprintf(rootChildrenURITemplate, driveID)
	}

	// Execute the children request
	return drives.
		NewItemItemsRequestBuilder(rawURL, client.RequestAdapter).
		Get(ctx, nil)
}

func (s *DriveService) ChildrenIterator(ctx context.Context, folderPath string) iter.Seq2[string, error] {
	// Retrieve the first page of children
	resp, err := s.getChildren(ctx, folderPath)
	if err != nil {
		return func(yield func(string, error) bool) {
			yield("", fmt.Errorf("unable to retrieve children: %w", err))
		}
	}

	// Get the Graph client (needed for paging)
	client, err := s.graph.Client(ctx)
	if err != nil {
		return func(yield func(string, error) bool) {
			yield("", err)
		}
	}

	// Build the page iterator
	pageIterator, err := msgraphcore.NewPageIterator[models.DriveItemable](
		resp,
		client.GetAdapter(),
		models.CreateDriveItemFromDiscriminatorValue,
	)
	if err != nil {
		return func(yield func(string, error) bool) {
			yield("", fmt.Errorf("unable to create page iterator: %w", err))
		}
	}

	// Return the iterator function
	return func(yield func(string, error) bool) {
		iterErr := pageIterator.Iterate(ctx, func(item models.DriveItemable) bool {
			var name string
			if item.GetName() != nil {
				name = *item.GetName()
			}
			return yield(name, nil)
		})

		if iterErr != nil {
			yield("", iterErr)
		}
	}
}
