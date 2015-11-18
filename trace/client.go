package trace

import "net/http"

type Transport struct {
	Transport http.RoundTripper
	Span      SpanID
	Collector Collector
}

func NewTransport(r *http.Request) Transport {
	span := ParseSpanID(r)
	return Transport{
		Span:      span,
		Collector: DefaultCollector,
	}
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var transport http.RoundTripper
	if t.Transport != nil {
		transport = t.Transport
	} else {
		transport = http.DefaultTransport
	}
	req2 := SetHeaders(req, t.Span)

	event := NewEvent(t.Span, ClientRequestSent)
	t.Collector.Record(event)

	resp, err := transport.RoundTrip(req2)

	event = NewEvent(t.Span, ClientRequestReceived)
	t.Collector.Record(event)

	return resp, err
}
