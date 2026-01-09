package profileservice

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

const (
	ProfileLoadedTopic  = "profile.loaded"
	ProfileSavedTopic   = "profile.saved"
	ProfileClearedTopic = "profile.cleared"
)

type ProfileEvent struct {
	topic  string
	record azidentity.AuthenticationRecord
}

func newProfileEvent(topic string, record azidentity.AuthenticationRecord) ProfileEvent {
	return ProfileEvent{topic: topic, record: record}
}

func (e ProfileEvent) Topic() string {
	return e.topic
}

func (e ProfileEvent) Record() azidentity.AuthenticationRecord {
	return e.record
}

func newProfileClearedEvent() ProfileEvent {
	return ProfileEvent{topic: ProfileClearedTopic, record: azidentity.AuthenticationRecord{}}
}

func newProfileLoadedEvent(record azidentity.AuthenticationRecord) ProfileEvent {
	return newProfileEvent(ProfileLoadedTopic, record)
}

func newProfileSavedEvent(record azidentity.AuthenticationRecord) ProfileEvent {
	return newProfileEvent(ProfileSavedTopic, record)
}
