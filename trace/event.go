package trace

import (
	"fmt"
	"time"
)

const (
	RequestReceived = iota
)

type Event struct {
	EventID   ID
	Span      SpanID
	CreatedAt time.Time
	Type      int
}

func (e Event) String() string {
	return fmt.Sprintf("[ID: %s, Trace: %s, CreatedAt: %s]", e.EventID.String(), e.Span.Trace.String(), e.CreatedAt)
}

func NewEvent(span SpanID) Event {
	return Event{
		CreatedAt: time.Now(),
		EventID:   NewID(),
		Span:      span,
	}
}
