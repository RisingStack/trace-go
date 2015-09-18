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
		trace.GetId(r),
	}
	log.Println("Id in 2nd endpoint: " + trace.GetId(r) + " on URL: " + r.URL.String())
	json.NewEncoder(rw).Encode(resp)
}

func handle(rw http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		trace.GetId(r),
	}
	log.Println("Id in 1st endpoint: " + trace.GetId(r) + " on URL: " + r.URL.String())
	_, err := trace.Get("http://localhost:9876/test2", r)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(rw).Encode(resp)
}

func main() {
	http.HandleFunc("/test", trace.Instrument(handle))
	http.HandleFunc("/test2", trace.Instrument(handle2))
	log.Println("Listening on :9876")
	log.Println(http.ListenAndServe(":9876", nil))
}
