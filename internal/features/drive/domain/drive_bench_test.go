package drive

import (
	"context"
	"fmt"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/mock"
)

func BenchmarkDefaultService_ListDrives(b *testing.B) {
	ctx := context.Background()
	source := new(mockDriveSource)
	mounts := new(mockMountProvider)
	svc := NewDefaultService(source, mounts, &dummyLogger{})

	// Setup 10 mounts, each with 10 drives
	var mountConfigs []mount.MountConfig
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/mount%d", i)
		mountConfigs = append(mountConfigs, mount.MountConfig{Path: path, IdentityID: "user1"})

		var drives []fs.Drive
		for j := 0; j < 10; j++ {
			drives = append(drives, fs.Drive{ID: fmt.Sprintf("d%d-%d", i, j), Name: fmt.Sprintf("Drive %d-%d", i, j)})
		}
		source.On("ListDrives", mock.Anything, path).Return(drives, nil)
	}
	mounts.On("ListMounts", mock.Anything).Return(mountConfigs, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.ListDrives(ctx, "")
	}
}

func BenchmarkDefaultService_ResolveDrive(b *testing.B) {
	ctx := context.Background()
	source := new(mockDriveSource)
	mounts := new(mockMountProvider)
	svc := NewDefaultService(source, mounts, &dummyLogger{})

	// Setup 1 mount with 100 drives
	mounts.On("ListMounts", mock.Anything).Return([]mount.MountConfig{{Path: "/od", IdentityID: "user1"}}, nil)

	var drives []fs.Drive
	for j := 0; j < 100; j++ {
		drives = append(drives, fs.Drive{ID: fmt.Sprintf("d%d", j), Name: fmt.Sprintf("Drive %d", j)})
	}
	source.On("ListDrives", mock.Anything, "/od").Return(drives, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Resolve the last drive to test near worst-case
		_, _ = svc.ResolveDrive(ctx, "Drive 99", "")
	}
}

func BenchmarkNewDriveType(b *testing.B) {
	inputs := []string{"personal", "business", "sharepoint", "unknown", "invalid"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewDriveType(inputs[i%len(inputs)])
	}
}
