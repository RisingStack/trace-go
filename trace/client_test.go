package trace

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type server struct {
	LastSpanID SpanID
}

func (s server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	span := ParseSpanID(r)
	fmt.Fprintf(rw, "%#v", span)
}

func TestClient(t *testing.T) {
	s := server{}
	ts := httptest.NewServer(s)
	defer ts.Close()
	or, err := http.NewRequest("GET", "example.com", nil)
	if err != nil {
		t.Errorf("Failed to create Request")
	}
	os := NewRootSpanID()
	or2 := SetHeaders(or, os)
	c := http.Client{Transport: NewTransport(or2)}
	r, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Failed to create Request")
	}
	rq, err := c.Do(r)
	if err != nil {
		t.Errorf("Failed to execute request: %#v", r)
	}
	defer rq.Body.Close()
	b, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		t.Error("Failed to read request")
	}
	if fmt.Sprintf("%#v", os) != string(b) {
		t.Errorf("SpanID (%#v) is not the same as provided (%#v)", b, os)
	}
}
