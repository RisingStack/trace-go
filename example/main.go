package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RisingStack/trace-go/trace"
)

var collector = trace.NewMemoryCollector()

func handle2(rw http.ResponseWriter, r *http.Request) {
	span := trace.NewSpanIDFromRequest(r)
	log.Println("Id in 2nd endpoint: " + span.String() + " on URL: " + r.URL.String())
	resp := struct {
		Message string `json:"message"`
	}{
		"Message from Endpoint2 to Endpoint1",
	}
	json.NewEncoder(rw).Encode(resp)
}

func handle(rw http.ResponseWriter, r *http.Request) {
	msg := struct {
		Message string `json:"message"`
	}{
		"Message from Endpoint1 to Endpoint2",
	}
	span := trace.NewSpanIDFromRequest(r)
	log.Println("Id in 1st endpoint: " + span.String())
	c := http.Client{
		Transport: trace.Transport{
			Span:      span,
			Collector: collector,
		},
	}
	_, err := c.Get("http://localhost:9876/test2")
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(rw).Encode(msg)
}

func eventsHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Eventshandler called")
	events := collector.GetEvents()
	msg := struct {
		Events []trace.Event `json:"events"`
	}{
		events,
	}
	json.NewEncoder(rw).Encode(msg)
}

func main() {
	http.HandleFunc("/test", trace.Trace(handle, collector))
	http.HandleFunc("/test2", trace.Trace(handle2, collector))
	http.HandleFunc("/events", trace.Trace(eventsHandler, collector))

	log.Println("Listening on :9876")
	log.Println(http.ListenAndServe(":9876", nil))
}
