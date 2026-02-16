package file

const (
	baseURL                          = "https://graph.microsoft.com/v1.0"
	rootChildrenURITemplate2         = "{+baseurl}/drives/{drive_id}/root/children"
	rootRelativeChildrenURITemplate2 = "{+baseurl}/drives/{drive_id}/root:{path}:/children"
	rootRelativeURITemplate2         = "{+baseurl}/drives/{drive_id}/root:{path}:"
	rootURITemplate2                 = "{+baseurl}/drives/{drive_id}/root"
)
