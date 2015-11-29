package trace

import "net/http"

// DefaultCollector is a MemoryCollector. Trace defaults back to DefaultCollector if none given.
var DefaultCollector = NewMemoryCollector()

// Trace is used to add incoming Trace functionality.
type Trace struct {
	collector Collector
}

func (t *Trace) process(rw http.ResponseWriter, r *http.Request, f http.Handler) {
	if t.collector == nil {
		t.collector = DefaultCollector
	}

	span := NewSpanIDFromRequest(r)
	event := NewEvent(span, RequestReceived)
	t.collector.Record(event)

	r2 := SetHeaders(r, span)
	f.ServeHTTP(rw, r2)

	event = NewEvent(span, RequestCompleted)
	t.collector.Record(event)
}

// Handler wraps an existing http.Handler with Trace functionality.
func (t *Trace) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		t.process(rw, r, h)
	})
}

// HandlerFunc wraps an existing http.HandlerFunc with Trace functionality.
func (t *Trace) HandlerFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t.process(rw, r, http.HandlerFunc(f))
	}

}
