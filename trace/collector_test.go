package trace

import "testing"

func TestRecord(t *testing.T) {
	c := NewMemoryCollector()
	e := NewEvent(NewRootSpanID(), RequestReceived)
	c.Record(e)
	events := c.GetEvents()
	if len(events) != 1 {
		t.Errorf("Number of events should be 1")
	}
	if events[0] != e {
		t.Errorf("Event list doesn't contain recorded event")
	}
}
