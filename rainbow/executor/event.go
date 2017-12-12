package executor

type EventVariant uint

const (
	EventStartVariant EventVariant = iota
	EventStopVariant
)

type Event struct {
	Variant EventVariant `json:"variant"`
	Start   *EventStart  `json:"start"`
	Stop    *EventStop   `json:"stop"`
}

func NewStartEvent(namespace, instanceID string) Event {
	return Event{
		Variant: EventStartVariant,
		Start:   &EventStart{Namespace: namespace, InstanceID: instanceID},
	}
}

func NewStopEvent(namespace, instanceID, instanceRemote string) Event {
	return Event{
		Variant: EventStopVariant,
		Stop:    &EventStop{Namespace: namespace, InstanceID: instanceID, InstanceRemote: instanceRemote},
	}
}

type EventStart struct {
	Namespace  string `json:"namespace"`
	InstanceID string `json:"instance_id"`
}

type EventStop struct {
	Namespace      string `json:"namespace"`
	InstanceID     string `json:"instance_id"`
	InstanceRemote string `json:"instance_remote"`
}
