package executor

type EventVariant uint

const (
	EventStartVariant EventVariant = iota
	EventStopVariant
)

type Event struct {
	Variant EventVariant
	Start   *EventStart
	Stop    *EventStop
}

func NewStartEvent(id string) Event {
	return Event{
		Variant: EventStartVariant,
		Start:   &EventStart{ID: id},
	}
}

func NewStopEvent(id, remote string) Event {
	return Event{
		Variant: EventStopVariant,
		Stop:    &EventStop{ID: id, Remote: remote},
	}
}

type EventStart struct {
	ID string
}

type EventStop struct {
	ID     string
	Remote string
}
