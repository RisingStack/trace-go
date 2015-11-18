package trace

import (
	"log"
	"net/http"
)

// Generates a new SpanID from a request. If there is no span information in the request, a root SpanID will be generated.
func NewSpanIDFromRequest(req *http.Request) SpanID {
	parentSpanID := ParseSpanID(req)
	if parentSpanID.Empty() {
		return NewRootSpanID()
	}
	return NewSpanID(parentSpanID)
}

// Returns the SpanID contained in the request.
func ParseSpanID(req *http.Request) SpanID {
	spanIDStr := req.Header.Get(HeaderSpanID)
	traceIDStr := req.Header.Get(HeaderTraceID)
	spanID, err := ParseID(spanIDStr)
	if err != nil {
		log.Println("Failed to parse spanID string value, using 0.")
	}
	traceID, err := ParseID(traceIDStr)
	if err != nil {
		log.Println("Failed to parse traceID string value, using 0.")
	}
	return SpanID{
		Trace: traceID,
		Span:  spanID,
	}
}

func SetHeaders(req *http.Request, s SpanID) *http.Request {
	// https://groups.google.com/forum/#!topic/golang-nuts/-j6p12SSpXI
	req2 := copyRequest(req)
	req2.Header.Set(HeaderSpanID, s.Span.String())
	req2.Header.Set(HeaderTraceID, s.Trace.String())
	return req2
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
