package cli

import (
	"context"
	"errors"
	"fmt"
	"testing"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Types -------------------------------------------------------------

func AsType[T any](args mock.Arguments, index int) T {
	obj := args.Get(index)
	var s T
	var ok bool
	if obj == nil {
		return s
	}
	if s, ok = obj.(T); !ok {
		panic(fmt.Sprintf("assert: arguments: AsType[%T](%d) failed because object wasn't correct type: %v", s, index, obj))
	}
	return s
}

type mockURI struct{ mock.Mock }

func (m *mockURI) FromString(s string) (*fs.URI, error) {
	args := m.Called(s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fs.URI), args.Error(1)
}

type mockMountLister struct{ mock.Mock }

func (m *mockMountLister) ListMounts(ctx context.Context) ([]mount.MountConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).([]mount.MountConfig), args.Error(1)
}

type mockItemLister struct{ mock.Mock }

func (m *mockItemLister) List(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ListOptions) ([]pkgfs.Item, error) {
	args := m.Called(ctx, uri, opts)

	return AsType[[]pkgfs.Item](args, 0), args.Error(1)
}

// ---------------------------------------------------------------------------

func TestProviderPathCompletion(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		toComplete     string
		mockSetup      func(*mockURI, *mockMountLister, *mockItemLister)
		expected       []string
		expectedDirect cobra.ShellCompDirective
	}{
		// -------------------------------------------------------------------
		{
			name:       "mount listing fallback when URI parse fails and no slash",
			args:       []string{},
			toComplete: "lo",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				u.On("FromString", "lo").Return(nil, errors.New("bad uri"))
				m.On("ListMounts", mock.Anything).Return([]mount.MountConfig{
					{Path: "/local"},
					{Path: "/remote"},
				}, nil)
			},
			expected:       []string{"local:"},
			expectedDirect: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp,
		},
		// -------------------------------------------------------------------
		{
			name:       "URI parse error but contains slash → no mount fallback",
			args:       []string{},
			toComplete: "abc/def",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				u.On("FromString", "abc/def").Return(nil, errors.New("bad uri"))
			},
			expected:       nil,
			expectedDirect: cobra.ShellCompDirectiveNoFileComp,
		},
		// -------------------------------------------------------------------
		{
			name:       "empty path returns root suggestion",
			args:       []string{},
			toComplete: "local:",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				u.On("FromString", "local:").Return(&fs.URI{
					Provider: "local",
					Path:     "",
				}, nil)
			},
			expected:       []string{"local:/"},
			expectedDirect: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp,
		},
		// -------------------------------------------------------------------
		{
			name:       "path without slash returns root",
			args:       []string{},
			toComplete: "local:abc",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				u.On("FromString", "local:abc").Return(&fs.URI{
					Provider: "local",
					Path:     "abc",
				}, nil)
			},
			expected:       []string{"local:/"},
			expectedDirect: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp,
		},
		// -------------------------------------------------------------------
		{
			name:       "directory listing returns matching items",
			args:       []string{},
			toComplete: "local:/foo/ba",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				u.On("FromString", "local:/foo/ba").Return(&fs.URI{
					Provider: "local",
					Path:     "/foo/ba",
				}, nil)

				i.On("List", mock.Anything, &pkgfs.URI{
					Provider: "local",
					Path:     "/foo",
				}, pkgfs.ListOptions{}).Return([]pkgfs.Item{
					{Name: "bar", Type: pkgfs.TypeFile},
					{Name: "baz", Type: pkgfs.TypeFolder},
					{Name: "zzz", Type: pkgfs.TypeFile},
				}, nil)
			},
			expected:       []string{"/foo/bar", "/foo/baz/"},
			expectedDirect: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp,
		},
		// -------------------------------------------------------------------
		{
			name:       "itemLister error returns no file comp",
			args:       []string{},
			toComplete: "local:/foo/bar",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				u.On("FromString", "local:/foo/bar").Return(&fs.URI{
					Provider: "local",
					Path:     "/foo/bar",
				}, nil)

				i.On("List", mock.Anything, &pkgfs.URI{
					Provider: "local",
					Path:     "/foo",
				}, pkgfs.ListOptions{}).Return(nil, errors.New("list error"))
			},
			expected:       nil,
			expectedDirect: cobra.ShellCompDirectiveNoFileComp,
		},
		// -------------------------------------------------------------------
		{
			name:       "prefix stripping when previous arg ends with colon",
			args:       []string{"local:"},
			toComplete: "ba",
			mockSetup: func(u *mockURI, m *mockMountLister, i *mockItemLister) {
				// prefixToStrip = "local:"
				u.On("FromString", "local:ba").Return(&fs.URI{
					Provider: "local",
					Path:     "ba",
				}, nil)
			},
			expected:       []string{"/"},
			expectedDirect: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp,
		},
	}

	// -----------------------------------------------------------------------

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURI := &mockURI{}
			mockMount := &mockMountLister{}
			mockItems := &mockItemLister{}

			if tt.mockSetup != nil {
				tt.mockSetup(mockURI, mockMount, mockItems)
			}

			cmd := &cobra.Command{}
			cmd.SetContext(context.Background())

			comp := ProviderPathCompletion(mockItems, mockURI, mockMount)

			got, directive := comp(cmd, tt.args, tt.toComplete)

			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.expectedDirect, directive)

			mockURI.AssertExpectations(t)
			mockMount.AssertExpectations(t)
			mockItems.AssertExpectations(t)
		})
	}
}
