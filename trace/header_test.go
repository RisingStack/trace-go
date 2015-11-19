package trace

import (
	"net/http"
	"testing"
)

const (
	HeaderName  = "My-Header"
	HeaderValue = "My-Value"
)

func TestSetHeaders(t *testing.T) {
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	r.Header.Set(HeaderName, HeaderValue)
	s := NewRootSpanID()
	r2 := SetHeaders(r, s)
	if r2.Header.Get(HeaderName) != HeaderValue {
		t.Error("SetHeaders should copy headers from original request.")
	}
	if r2.Header.Get(HeaderTraceID) != s.Trace.String() {
		t.Errorf("Trace ID in request doesn't match: %016x != %016x", r2.Header.Get(HeaderTraceID), s.Trace)
	}
	if r2.Header.Get(HeaderSpanID) != s.Span.String() {
		t.Errorf("Span ID in request doesn't match: %016x != %016x", r2.Header.Get(HeaderSpanID), s.Trace)
	}
	s2 := ParseSpanID(r2)
	if s2 != s {
		t.Errorf("The two span doesn't match: %v != %v", s, s2)
	}
}

func TestNewFromEmptyRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	s := NewSpanIDFromRequest(r)
	if s.Parent != ID(0) {
		t.Errorf("Span from empty request should be root. This is not: %v", s)
	}
	if s.Trace == ID(0) {
		t.Errorf("Span from request should have non-empty Trace ID. %v", s)
	}
	if s.Span == ID(0) {
		t.Errorf("Span from request should have non-empty Span ID. %v", s)
	}
}

func TestNewFromNonEmptyRequest(t *testing.T) {
	so := NewRootSpanID()
	r, err := http.NewRequest("GET", "http://localhost:12300", nil)
	if err != nil {
		t.Error("Failed to create request. ", err)
	}
	r2 := SetHeaders(r, so)
	s := NewSpanIDFromRequest(r2)
	if s.Parent != so.Span {
		t.Errorf("Span from request should have the parent set to the Span ID. Original: %v. New: %v", so, s)
	}
	if s.Trace != so.Trace {
		t.Errorf("Span from request should have Trace ID from request. Original: %v. New: %v", so, s)
	}
	if s.Span == so.Span || s.Span == ID(0) {
		t.Errorf("Span from request should have non-empty, uniqe Span ID. Original: %v. New: %v", so, s)
	}
}
