package infra

import (
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func toDomainDrive(g models.Driveable) *domaindrive.Drive {
	if g == nil {
		return nil
	}

	return &domaindrive.Drive{
		ID:   deref(g.GetId()),
		Name: deref(g.GetName()),
		Type: domaindrive.NewDriveType(deref(g.GetDriveType())),
	}
}

func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}
