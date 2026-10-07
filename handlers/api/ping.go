package api

import (
	"encoding/json"
	"fmt"
	"homework/models"
	"log"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := models.Response[any]{
		Success: true,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(data))
}
