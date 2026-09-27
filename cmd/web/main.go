package main

import (
	"homework/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("fatal f: %v", err)
	}
}
