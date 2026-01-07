package root

import (
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

var (
	credentialService  CredentialService
	profileService     ProfileService
	graphClientService ClientService
	driveService       driveChildIterator
	logger             logging.Logger
	graphClient        *msgraphsdkgo.GraphServiceClient
)
