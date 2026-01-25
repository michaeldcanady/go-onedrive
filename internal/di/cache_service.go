package di

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	cliprofileservicego "github.com/michaeldcanady/go-onedrive/internal/app/cli_profile_service.go"
	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	"github.com/michaeldcanady/go-onedrive/internal/config"
)

type CacheService interface {
	GetProfile(context.Context, string) (azidentity.AuthenticationRecord, error)
	SetProfile(ctx context.Context, name string, record azidentity.AuthenticationRecord) error
	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
	SetConfiguration(ctx context.Context, name string, record config.Configuration3) error
	GetCLIProfile(ctx context.Context, name string) (cliprofileservicego.Profile, error)
	SetCLIProfile(ctx context.Context, name string, profile cliprofileservicego.Profile) error
	driveservice.CacheService
}
