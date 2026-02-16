package file

import abstractions "github.com/microsoft/kiota-abstractions-go"

type RelativePathContentsBuilder struct {
	abstractions.BaseRequestBuilder
}

const (
	relativePathContentsURLTemplate = "{+baseurl}/drives/{drive_id}/root:{path}:"
)

func NewRelativePathContentsBuilder(pathParameters map[string]string, requestAdapter abstractions.RequestAdapter) *RelativePathContentsBuilder {
	return &RelativePathContentsBuilder{
		BaseRequestBuilder: *abstractions.NewBaseRequestBuilder(requestAdapter, relativePathContentsURLTemplate, pathParameters),
	}
}
