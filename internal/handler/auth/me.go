package auth

import (
	"encoding/json"
	"errors"
	"homework/internal/model"
	"log"
	"net/http"
)

func Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("user_session")
	if err != nil {
		resp := model.Response[[]any]{
			Success: false,
			Body: []any{
				map[string]any{"authorized": false},
			},
		}

		var status int
		if errors.Is(err, http.ErrNoCookie) {
			status = http.StatusUnauthorized
			resp.Error = "no session cookie"
		} else {
			status = http.StatusInternalServerError
			resp.Error = "internal error"
		}

		w.WriteHeader(status)

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(data)
		return
	}

	resp := model.Response[map[string]any]{
		Success: true,
		Body: map[string]any{
			"authorized": true,
			"cookie":     cookie.Value,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
