package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RisingStack/trace-go/trace"
	"github.com/codegangsta/negroni"
)

var collector = trace.NewMemoryCollector()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!\n")
	})

	mux.HandleFunc("/events", func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(collector.GetEvents())
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Use(trace.NewNegroniHandler(collector))
	n.Run(":3000")
}
