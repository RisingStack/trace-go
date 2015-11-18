package trace

import "net/http"

const (
	// HeaderSpanID is the header name in requests for the SpanID.
	HeaderSpanID = "Trace-Span"
	// HeaderTraceID is the header name in requests for the TraceID.
	HeaderTraceID = "Trace-Trace"
)

// DefaultCollector is a MemoryCollector. Trace defaults back to DefaultCollector if none given.
var DefaultCollector = NewMemoryCollector()

// Trace is used to add incoming Trace functionality.
type Trace struct {
	collector Collector
}

func (t *Trace) process(rw http.ResponseWriter, r *http.Request, f http.HandlerFunc) {
	if t.collector == nil {
		t.collector = DefaultCollector
	}

	span := NewSpanIDFromRequest(r)
	event := NewEvent(span, RequestReceived)
	t.collector.Record(event)

	r2 := SetHeaders(r, span)
	f(rw, r2)

	event = NewEvent(span, RequestCompleted)
	t.collector.Record(event)
}

// Handler wraps an existing http.Handler with Trace functionality.
func (t *Trace) Handler(h http.Handler) http.Handler {
	f := h.ServeHTTP
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		t.process(rw, r, f)
	})
}

// HandlerFunc wraps an existing http.HandlerFunc with Trace functionality.
func (t *Trace) HandlerFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t.process(rw, r, f)
	}
}

// HandlerFuncNext middleware function can be used for Negroni integration.
func (t *Trace) HandlerFuncNext(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	t.process(rw, r, next)
}
