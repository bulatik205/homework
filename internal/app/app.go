package app

import (
	"homework/internal/handler"
	"homework/internal/handler/status"
	"net/http"
	"time"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/api/v1/getStatus", status.Get)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}
