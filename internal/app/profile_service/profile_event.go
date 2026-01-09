package profileservice

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

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
