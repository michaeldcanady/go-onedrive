package auth

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

const (
	CredentialLoadedTopic = "credential.loaded"
)

type CredentialEvent struct {
	topic      string
	credential azcore.TokenCredential
}

func newCredentialEvent(topic string, cred azcore.TokenCredential) CredentialEvent {
	return CredentialEvent{topic: topic, credential: cred}
}

func (e CredentialEvent) Topic() string {
	return e.topic
}

func (e CredentialEvent) Credential() azcore.TokenCredential {
	return e.credential
}

func newCredentialLoadedEvent(cred azcore.TokenCredential) CredentialEvent {
	return newCredentialEvent(CredentialLoadedTopic, cred)
}
