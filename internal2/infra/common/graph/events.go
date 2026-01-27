package clientservice

const (
	GraphClientInitializedTopic = "graph.client.initialized"
)

type GraphClientEvent struct {
	topic string
}

func (e GraphClientEvent) Topic() string {
	return e.topic
}

func newGraphClientInitializedEvent() GraphClientEvent {
	return GraphClientEvent{
		topic: GraphClientInitializedTopic,
	}
}
