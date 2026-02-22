package file

const (
	// baseURL is the primary endpoint for the Microsoft Graph API v1.0.
	baseURL = "https://graph.microsoft.com/v1.0"
	// rootChildrenURITemplate2 defines the URL template for listing top-level children
	// of a drive.
	rootChildrenURITemplate2 = "{+baseurl}/drives/{drive_id}/root/children"
	// rootRelativeChildrenURITemplate2 defines the URL template for listing children
	// of a drive item at a relative path.
	rootRelativeChildrenURITemplate2 = "{+baseurl}/drives/{drive_id}/root:{path}:/children"
	// rootRelativeURITemplate2 defines the URL template for fetching metadata for
	// a drive item at a relative path.
	rootRelativeURITemplate2 = "{+baseurl}/drives/{drive_id}/root:{path}:"
	// rootURITemplate2 defines the URL template for fetching metadata for the
	// root of a drive.
	rootURITemplate2 = "{+baseurl}/drives/{drive_id}/root"
	// rootRelativeContentURITemplate2 defines the URL template for downloading or
	// uploading content to a drive item at a relative path.
	rootRelativeContentURITemplate2 = "{+baseurl}/drives/{drive_id}/root:{path}:/content"
)
