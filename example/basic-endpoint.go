package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RisingStack/trace-go/trace"
)

func handle2(rw http.ResponseWriter, r *http.Request) {
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
	c := http.Client{Transport: trace.NewTransport(r)}
	_, err := c.Get("http://localhost:9876/test2")
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(rw).Encode(msg)
}

func eventsHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Eventshandler called")
	events := trace.DefaultCollector.GetEvents()
	msg := struct {
		Events []trace.Event `json:"events"`
	}{
		events,
	}
	json.NewEncoder(rw).Encode(msg)
}

func main() {
	http.HandleFunc("/test", trace.Trace(handle))
	http.HandleFunc("/test2", trace.Trace(handle2))
	http.HandleFunc("/events", eventsHandler)

	log.Println("Listening on :9876")
	log.Println(http.ListenAndServe(":9876", nil))
}
