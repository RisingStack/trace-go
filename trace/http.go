package trace

import "net/http"

const (
	HeaderSpanID  = "Trace-Span"
	HeaderTraceID = "Trace-Trace"
)

var DefaultCollector = NewMemoryCollector()

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

func (t *Trace) Handler(h http.Handler) http.Handler {
	f := h.ServeHTTP
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		t.process(rw, r, f)
	})
}

func (t *Trace) HandlerFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		t.process(rw, r, f)
	}
}

// Middleware function for Negroni integration
func (t *Trace) HandlerFuncNext(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	t.process(rw, r, next)
}
