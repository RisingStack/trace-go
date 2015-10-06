package trace

import "time"

const (
	RequestReceived = iota
)

type Event struct {
	SpanID    ID
	ParentID  ID
	TraceID   ID
	CreatedAt time.Time
	Type      int
}

func NewEvent() Event {
	return Event{}
}
