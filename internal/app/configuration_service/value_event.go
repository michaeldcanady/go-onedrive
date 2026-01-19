package configurationservice

const (
	ConfigurationUpdatedTopic = "configuration.updated"
)

type ValueEvent struct {
	topic string
	key   string
	old   interface{}
	value interface{}
}

func newEvent(topic, key string, old, value interface{}) *ValueEvent {
	return &ValueEvent{
		topic: topic,
		key:   key,
		old:   old,
		value: value,
	}
}

// Topic returns the event topic.
func (e *ValueEvent) Topic() string {
	return e.topic
}

// Key returns the configuration key associated with the event.
func (e *ValueEvent) Key() string {
	return e.key
}

// Old returns the old value associated with the event.
func (e *ValueEvent) Old() interface{} {
	return e.old
}

// Value returns the new value associated with the event.
func (e *ValueEvent) Value() interface{} {
	return e.value
}

func newConfigurationUpdatedEvent(key string, old, value interface{}) *ValueEvent {
	return newEvent(ConfigurationUpdatedTopic, key, old, value)
}
