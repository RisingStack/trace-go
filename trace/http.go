package trace

import "net/http"

const (
	HeaderSpanID  = "Trace-Span"
	HeaderTraceID = "Trace-Trace"
)

type Tracer struct {
	collector Collector
	handler   http.Handler
}

func New(h http.Handler, c Collector) Tracer {
	return Tracer{
		collector: c,
		handler:   h,
	}
}

func (t Tracer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	span := NewSpanIDFromRequest(r)
	event := NewEvent(span, RequestReceived)
	t.collector.Record(event)
	t.handler.ServeHTTP(rw, r)
	event = NewEvent(span, RequestCompleted)
}

func Trace(fn http.HandlerFunc, collector Collector) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		span := NewSpanIDFromRequest(r)
		event := NewEvent(span, RequestReceived)
		collector.Record(event)
		fn(rw, r)
		event = NewEvent(span, RequestCompleted)
		collector.Record(event)
	}
}

type Transport struct {
	Transport http.RoundTripper
	Span      SpanID
	Collector Collector
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var transport http.RoundTripper
	if t.Transport != nil {
		transport = t.Transport
	} else {
		transport = http.DefaultTransport
	}
	req2 := t.setHeaders(req)

	// TODO: Here comes the client send event
	event := NewEvent(t.Span, ClientRequestSent)
	t.Collector.Record(event)

	resp, err := transport.RoundTrip(req2)

	// TODO: Here comes the client receive event
	event = NewEvent(t.Span, ClientRequestReceived)
	t.Collector.Record(event)

	return resp, err
}

func (t *Transport) setHeaders(req *http.Request) *http.Request {
	// https://groups.google.com/forum/#!topic/golang-nuts/-j6p12SSpXI
	req2 := copyRequest(req)
	req2.Header.Set(HeaderSpanID, t.Span.Span.String())
	req2.Header.Set(HeaderTraceID, t.Span.Trace.String())
	return req2
}

func GetHeaders(req *http.Request) SpanID {
	spanID := NewSpanIDFromRequest(req)
	return spanID
}

func copyRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
