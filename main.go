package main

import (
	"fmt"
	"homework/config"
	"homework/db"
	"homework/handlers/api"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})
	mux.HandleFunc("/api/v1/ping", api.Ping)
	mux.HandleFunc("/api/v1/tasks", api.GetTasks(database))

	log.Printf("listening on 127.0.0.1:%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+cfg.ServerPort, mux))
}
