package profileservice

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

const (
	ProfileLoadedTopic = "profile.loaded"
)

func newProfileLoadedEvent(record azidentity.AuthenticationRecord) ProfileEvent {
	return newProfileEvent(ProfileLoadedTopic, record)
}
