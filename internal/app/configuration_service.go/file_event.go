package configurationservice

const (
	FileChangedTopic = "configuration.file.changed"
)

type FileEvent struct {
	topic string
	old   string
	path  string
}

func newFileEvent(topic, old, path string) *FileEvent {
	return &FileEvent{
		topic: topic,
		old:   old,
	}
}

// Topic returns the event topic.
func (e *FileEvent) Topic() string {
	return e.topic
}

// Old returns the old file path associated with the event.
func (e *FileEvent) Old() string {
	return e.old
}

// Path returns the new file path associated with the event.
func (e *FileEvent) Path() string {
	return e.path
}

func newConfigurationFileChangedEvent(old, path string) *FileEvent {
	return newFileEvent(FileChangedTopic, old, path)
}
