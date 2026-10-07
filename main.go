package main

import (
	"fmt"
	"homework/handlers/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})
	mux.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		api.Ping(w, r)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8082", mux))
}
