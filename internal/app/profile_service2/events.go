package cacheservice

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

const (
	ProfileUpdatedEventTopic = "profile.updated"
	ProfileDeletedEventTopic = "profile.deleted"
)

type ProfileEvent struct {
	topic string
	old   *azidentity.AuthenticationRecord
	new   *azidentity.AuthenticationRecord
}

// newProfileEvent creates a new profile event.
func newProfileEvent(topic string, old, new *azidentity.AuthenticationRecord) *ProfileEvent {
	return &ProfileEvent{
		topic: topic,
		old:   old,
		new:   new,
	}
}

// Topic returns the event topic.
func (e *ProfileEvent) Topic() string {
	return e.topic
}

// Profile returns the authentication profile associated with the event.
func (e *ProfileEvent) Profile() *azidentity.AuthenticationRecord {
	return e.new
}

// newProfileUpdatedEvent creates a new profile updated event.
func newProfileUpdatedEvent(old, new *azidentity.AuthenticationRecord) *ProfileEvent {
	return newProfileEvent(ProfileUpdatedEventTopic, old, new)
}

// newProfileDeletedEvent creates a new profile deleted event.
func newProfileDeletedEvent(old *azidentity.AuthenticationRecord) *ProfileEvent {
	return newProfileEvent(ProfileDeletedEventTopic, old, nil)
}
