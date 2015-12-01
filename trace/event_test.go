package trace

import (
	"fmt"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	s := SpanID{
		Trace:  ID(1),
		Span:   ID(2),
		Parent: ID(3),
	}
	event := NewEvent(s, ClientRequestSent)
	if event.Type != ClientRequestSent {
		t.Error("Event type doesn't match")
	}
	if event.Span.Trace != ID(1) {
		t.Error("TraceID doesn't match")
	}
	if event.Span.Span != ID(2) {
		t.Error("SpanID doesn't match")
	}
	if event.Span.Parent != ID(3) {
		t.Error("ParentID doesn't match")
	}
	if event.EventID == ID(0) {
		t.Error("EventID = 0")
	}
	var et time.Time
	if event.CreatedAt == et {
		t.Error("CreatedAt attribute is empty")
	}
}

func TestEventString(t *testing.T) {
	s := SpanID{
		Trace:  ID(1),
		Span:   ID(2),
		Parent: ID(3),
	}
	e := NewEvent(s, RequestReceived)
	estr := e.String()
	expected := fmt.Sprintf("[ID: %s, Trace: %s, Type: %d, CreatedAt: %s]", e.EventID.String(), e.Span.Trace.String(), e.Type, e.CreatedAt)
	if estr != expected {
		t.Errorf("Event String(%v) doesn't match with expected (%v)", estr, expected)
	}
}
