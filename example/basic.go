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
	_, err := c.Get("http://index.hu")
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
	r := http.NewServeMux()
	r.HandleFunc("/test", handle)
	r.HandleFunc("/test2", handle2)
	r.HandleFunc("/events", eventsHandler)

	t := trace.New(r)

	log.Println("Listening on :9876")
	log.Println(http.ListenAndServe(":9876", t))
}
