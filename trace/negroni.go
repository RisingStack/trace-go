package trace

import "net/http"

type NegroniHandler struct {
	collector Collector
}

func NewNegroniHandler(c Collector) NegroniHandler {
	return NegroniHandler{
		c,
	}
}

func (g NegroniHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	span := NewSpanIDFromRequest(r)
	g.collector.Record(NewEvent(span, RequestReceived))
	next(rw, r)
	g.collector.Record(NewEvent(span, RequestCompleted))
}
