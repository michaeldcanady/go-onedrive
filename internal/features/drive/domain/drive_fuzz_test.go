package drive

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/mock"
)

func FuzzNewDriveType(f *testing.F) {
	f.Add("personal")
	f.Add("business")
	f.Add("sharepoint")
	f.Add("unknown")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		_ = NewDriveType(input)
	})
}

func FuzzResolveDrive(f *testing.F) {
	ctx := context.Background()
	source := new(mockDriveSource)
	mounts := new(mockMountProvider)
	svc := NewDefaultService(source, mounts, &dummyLogger{})

	f.Add("personal")
	f.Add("d1")
	f.Add("")

	f.Fuzz(func(t *testing.T, driveRef string) {
		// Setup expectations for each fuzz run to avoid panics on unexpected calls
		mounts.On("ListMounts", mock.Anything).Return([]mount.MountConfig{}, nil).Maybe()
		source.On("ListDrives", mock.Anything, mock.Anything).Return([]fs.Drive{}, nil).Maybe()

		_, _ = svc.ResolveDrive(ctx, driveRef, "")
	})
}
