package trace

import "net/http"

const (
	HeaderSpanID  = "Trace-Span"
	HeaderTraceID = "Trace-Trace"
)

var DefaultCollector = NewMemoryCollector()

type Tracer struct {
	collector Collector
	handler   http.Handler
}

func New(h http.Handler) Tracer {
	return Tracer{
		collector: DefaultCollector,
		handler:   h,
	}
}

func (t Tracer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	span := NewSpanIDFromRequest(r)
	event := NewEvent(span, RequestReceived)
	t.collector.Record(event)
	r2 := SetHeaders(r, span)
	t.handler.ServeHTTP(rw, r2)
	event = NewEvent(span, RequestCompleted)
	t.collector.Record(event)
}

func Trace(fn http.HandlerFunc) http.HandlerFunc {
	return TraceWithCollector(fn, DefaultCollector)
}

func TraceWithCollector(fn http.HandlerFunc, collector Collector) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		span := NewSpanIDFromRequest(r)
		event := NewEvent(span, RequestReceived)
		collector.Record(event)
		r2 := SetHeaders(r, span)
		fn(rw, r2)
		event = NewEvent(span, RequestCompleted)
		collector.Record(event)
	}
}
