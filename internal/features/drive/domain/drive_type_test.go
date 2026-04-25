package drive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDriveType_String(t *testing.T) {
	tests := []struct {
		name string
		dt   DriveType
		want string
	}{
		{name: "personal", dt: DriveTypePersonal, want: "personal"},
		{name: "business", dt: DriveTypeBusiness, want: "business"},
		{name: "sharepoint", dt: DriveTypeSharePoint, want: "sharepoint"},
		{name: "unknown", dt: DriveTypeUnknown, want: "unknown"},
		{name: "out of range", dt: DriveType(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.dt.String())
		})
	}
}

func TestNewDriveType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  DriveType
	}{
		{name: "personal", input: "personal", want: DriveTypePersonal},
		{name: "business", input: "business", want: DriveTypeBusiness},
		{name: "sharepoint", input: "sharepoint", want: DriveTypeSharePoint},
		{name: "unknown string", input: "invalid", want: DriveTypeUnknown},
		{name: "empty string", input: "", want: DriveTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewDriveType(tt.input))
		})
	}
}
