package events

type DomainEvent interface {
	EventType() string
	OccurredAtUnix() int64
}

type EventBase struct {
	eventType      string
	occurredAtUnix int64
}

func NewEventBase(eventType string, occurredAtUnix int64) EventBase {
	return EventBase{
		eventType:      eventType,
		occurredAtUnix: occurredAtUnix,
	}
}

func (e EventBase) EventType() string {
	return e.eventType
}

func (e EventBase) OccurredAtUnix() int64 {
	return e.occurredAtUnix
}