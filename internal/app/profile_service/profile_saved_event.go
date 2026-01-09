package profileservice

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

const (
	ProfileSavedTopic = "profile.saved"
)

func newProfileSavedEvent(record azidentity.AuthenticationRecord) ProfileEvent {
	return newProfileEvent(ProfileSavedTopic, record)
}
