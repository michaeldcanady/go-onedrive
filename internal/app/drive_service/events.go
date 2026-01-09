package driveservice

import "github.com/microsoftgraph/msgraph-sdk-go/models"

const (
	DriveChildrenLoadedTopic = "drive.children.loaded"
)

type DriveEvent struct {
	topic string
	path  string
	items []models.DriveItemable
}

func (e DriveEvent) Topic() string {
	return e.topic
}

func (e DriveEvent) Path() string {
	return e.path
}

func (e DriveEvent) Items() []models.DriveItemable {
	return e.items
}

func newDriveChildrenLoadedEvent(path string, items []models.DriveItemable) DriveEvent {
	return DriveEvent{
		topic: DriveChildrenLoadedTopic,
		path:  path,
		items: items,
	}
}
