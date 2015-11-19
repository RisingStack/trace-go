package trace

import (
	"fmt"
	"time"
)

const (
	RequestReceived = iota
	RequestCompleted
	ClientRequestSent
	ClientRequestReceived
)

type Event struct {
	EventID   ID
	Span      SpanID
	CreatedAt time.Time
	Type      int
}

func (e Event) String() string {
	return fmt.Sprintf("[ID: %s, Trace: %s, Type: %d, CreatedAt: %s]", e.EventID.String(), e.Span.Trace.String(), e.Type, e.CreatedAt)
}

func NewEvent(span SpanID, t int) Event {
	return Event{
		CreatedAt: time.Now().UTC(),
		EventID:   NewID(),
		Span:      span,
		Type:      t,
	}
}
