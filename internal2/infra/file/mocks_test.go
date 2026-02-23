package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	"github.com/stretchr/testify/mock"
)

// MockContentsCache mocks the ContentsCache interface
type MockContentsCache struct {
	mock.Mock
}

func (m *MockContentsCache) Get(ctx context.Context, id string) (*file.Contents, bool) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*file.Contents), args.Bool(1)
}

func (m *MockContentsCache) Put(ctx context.Context, id string, c *file.Contents) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}

func (m *MockContentsCache) Invalidate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockMetadataCache mocks the MetadataCache interface
type MockMetadataCache struct {
	mock.Mock
}

func (m *MockMetadataCache) Get(ctx context.Context, id string) (*file.Metadata, bool) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*file.Metadata), args.Bool(1)
}

func (m *MockMetadataCache) Put(ctx context.Context, meta *file.Metadata) error {
	args := m.Called(ctx, meta)
	return args.Error(0)
}

func (m *MockMetadataCache) Invalidate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockListingCache mocks the ListingCache interface
type MockListingCache struct {
	mock.Mock
}

func (m *MockListingCache) Get(ctx context.Context, path string) (*Listing, bool) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*Listing), args.Bool(1)
}

func (m *MockListingCache) Put(ctx context.Context, path string, l *Listing) error {
	args := m.Called(ctx, path, l)
	return args.Error(0)
}

func (m *MockListingCache) Invalidate(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

// MockRequestAdapter mocks the abstractions.RequestAdapter interface
type MockRequestAdapter struct {
	mock.Mock
}

func (m *MockRequestAdapter) Send(ctx context.Context, requestInfo *abstractions.RequestInformation, constructor serialization.ParsableFactory, errorMappings abstractions.ErrorMappings) (serialization.Parsable, error) {
	args := m.Called(ctx, requestInfo, constructor, errorMappings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(serialization.Parsable), args.Error(1)
}

func (m *MockRequestAdapter) SendEnum(ctx context.Context, requestInfo *abstractions.RequestInformation, parser serialization.EnumFactory, errorMappings abstractions.ErrorMappings) (any, error) {
	args := m.Called(ctx, requestInfo, parser, errorMappings)
	return args.Get(0), args.Error(1)
}

func (m *MockRequestAdapter) SendCollection(ctx context.Context, requestInfo *abstractions.RequestInformation, constructor serialization.ParsableFactory, errorMappings abstractions.ErrorMappings) ([]serialization.Parsable, error) {
	args := m.Called(ctx, requestInfo, constructor, errorMappings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]serialization.Parsable), args.Error(1)
}

func (m *MockRequestAdapter) SendEnumCollection(ctx context.Context, requestInfo *abstractions.RequestInformation, parser serialization.EnumFactory, errorMappings abstractions.ErrorMappings) ([]any, error) {
	args := m.Called(ctx, requestInfo, parser, errorMappings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]any), args.Error(1)
}

func (m *MockRequestAdapter) SendPrimitive(ctx context.Context, requestInfo *abstractions.RequestInformation, typeName string, errorMappings abstractions.ErrorMappings) (any, error) {
	args := m.Called(ctx, requestInfo, typeName, errorMappings)
	return args.Get(0), args.Error(1)
}

func (m *MockRequestAdapter) SendPrimitiveCollection(ctx context.Context, requestInfo *abstractions.RequestInformation, typeName string, errorMappings abstractions.ErrorMappings) ([]any, error) {
	args := m.Called(ctx, requestInfo, typeName, errorMappings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]any), args.Error(1)
}

func (m *MockRequestAdapter) SendNoContent(ctx context.Context, requestInfo *abstractions.RequestInformation, errorMappings abstractions.ErrorMappings) error {
	args := m.Called(ctx, requestInfo, errorMappings)
	return args.Error(0)
}

func (m *MockRequestAdapter) ConvertToNativeRequest(ctx context.Context, requestInfo *abstractions.RequestInformation) (any, error) {
	args := m.Called(ctx, requestInfo)
	return args.Get(0), args.Error(1)
}

func (m *MockRequestAdapter) GetSerializationWriterFactory() serialization.SerializationWriterFactory {
	return nil
}

func (m *MockRequestAdapter) EnableBackingStore(factory store.BackingStoreFactory) {
}

func (m *MockRequestAdapter) SetBaseUrl(baseUrl string) {
}

func (m *MockRequestAdapter) GetBaseUrl() string {
	return ""
}
