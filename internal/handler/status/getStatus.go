package status

import (
	"encoding/json"
	"fmt"
	"homework/internal/model"
	"log"
	"net/http"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := model.Response[any]{
		Success: true,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(data))
}
