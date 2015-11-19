package trace

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type TestCollector struct {
	list []Event
	lock sync.Mutex
}

func NewTestCollector() *TestCollector {
	c := TestCollector{}
	c.list = make([]Event, 0, NumberOfEventsBeforeFlush)
	c.lock = sync.Mutex{}
	return &c
}

func (c *TestCollector) Record(e Event) error {
	c.lock.Lock()
	c.list = append(c.list, e)
	c.lock.Unlock()
	return nil
}

func (c *TestCollector) GetEvents() []Event {
	c.lock.Lock()
	defer c.lock.Unlock()
	listLength := len(c.list)
	newList := make([]Event, listLength)
	copied := copy(newList, c.list)
	if copied != listLength {
		log.Panicln("Failed to copy Events from list")
	}
	return newList
}

func testHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello world!"))
}

func TestHandlerFuncRoot(t *testing.T) {
	c := NewTestCollector()
	tr := Trace{c}
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	tr.HandlerFunc(testHandler)(w, r)
	events := c.GetEvents()
	if len(events) != 2 {
		t.Errorf("After a call two events should be recorded")
	}
	if events[0].Span.Trace != events[1].Span.Trace {
		t.Errorf("Trace ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Span != events[1].Span.Span {
		t.Errorf("Span ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != events[1].Span.Parent {
		t.Errorf("Parent ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != ID(0) || events[1].Span.Parent != ID(0) {
		t.Errorf("Parent ID from empty request should be 0. %v", events)
	}
	if events[0].EventID == ID(0) || events[1].EventID == ID(0) {
		t.Errorf("EventID should be never 0: %v", events)
	}
	if events[0].EventID == events[1].EventID {
		t.Errorf("EventIDs should not be equal: %v.", events)
	}
	if events[0].Type == events[1].Type {
		t.Errorf("Types should not be equal: %v.", events)
	}
}

func TestHandlerFuncParent(t *testing.T) {
	c := NewTestCollector()
	tr := Trace{c}
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	sp := NewRootSpanID()
	r2 := SetHeaders(r, sp)
	tr.HandlerFunc(testHandler)(w, r2)
	events := c.GetEvents()
	if len(events) != 2 {
		t.Errorf("After a call two events should be recorded")
	}
	if events[0].Span.Trace != events[1].Span.Trace {
		t.Errorf("Trace ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Span != events[1].Span.Span {
		t.Errorf("Span ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != events[1].Span.Parent {
		t.Errorf("Parent ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != sp.Span || events[1].Span.Parent != sp.Span {
		t.Errorf("Parent ID from request with SpanID should not the parent's Span ID. %v", events)
	}
	if events[0].EventID == ID(0) || events[1].EventID == ID(0) {
		t.Errorf("EventID should be never 0: %v", events)
	}
	if events[0].EventID == events[1].EventID {
		t.Errorf("EventIDs should not be equal: %v.", events)
	}
	if events[0].Type == events[1].Type {
		t.Errorf("Types should not be equal: %v.", events)
	}
}

func TestHandlerRoot(t *testing.T) {
	log.Println("llsdjdsakfsf")
	c := NewTestCollector()
	tr := Trace{c}
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	tr.Handler(http.HandlerFunc(testHandler)).ServeHTTP(w, r)
	events := c.GetEvents()
	if len(events) != 2 {
		t.Errorf("After a call two events should be recorded")
	}
	if events[0].Span.Trace != events[1].Span.Trace {
		t.Errorf("Trace ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Span != events[1].Span.Span {
		t.Errorf("Span ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != events[1].Span.Parent {
		t.Errorf("Parent ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != ID(0) || events[1].Span.Parent != ID(0) {
		t.Errorf("Parent ID from empty request should be 0. %v", events)
	}
	if events[0].EventID == ID(0) || events[1].EventID == ID(0) {
		t.Errorf("EventID should be never 0: %v", events)
	}
	if events[0].EventID == events[1].EventID {
		t.Errorf("EventIDs should not be equal: %v.", events)
	}
	if events[0].Type == events[1].Type {
		t.Errorf("Types should not be equal: %v.", events)
	}
}

func TestHandlerParent(t *testing.T) {
	c := NewTestCollector()
	tr := Trace{c}
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	sp := NewRootSpanID()
	r2 := SetHeaders(r, sp)
	tr.Handler(http.HandlerFunc(testHandler)).ServeHTTP(w, r2)
	events := c.GetEvents()
	if len(events) != 2 {
		t.Errorf("After a call two events should be recorded")
	}
	if events[0].Span.Trace != events[1].Span.Trace {
		t.Errorf("Trace ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Span != events[1].Span.Span {
		t.Errorf("Span ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != events[1].Span.Parent {
		t.Errorf("Parent ID from same trace should be equal: %v", events)
	}
	if events[0].Span.Parent != sp.Span || events[1].Span.Parent != sp.Span {
		t.Errorf("Parent ID from request with SpanID should not the parent's Span ID. %v", events)
	}
	if events[0].EventID == ID(0) || events[1].EventID == ID(0) {
		t.Errorf("EventID should be never 0: %v", events)
	}
	if events[0].EventID == events[1].EventID {
		t.Errorf("EventIDs should not be equal: %v.", events)
	}
	if events[0].Type == events[1].Type {
		t.Errorf("Types should not be equal: %v.", events)
	}
}
