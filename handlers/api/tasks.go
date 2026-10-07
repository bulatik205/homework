package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"homework/models"
)

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()

		subject := q.Get("subject")
		recency := q.Get("recency")

		query := `SELECT id, subject, task, date_from, date_to, status, created_at, updated_at
		          FROM tasks WHERE 1=1`
		args := []any{}

		if subject != "" {
			query += " AND subject = ?"
			args = append(args, subject)
		}

		switch recency {
		case "last":
			query += " ORDER BY id DESC"
		default:
			query += " ORDER BY date_to IS NULL, date_to ASC, id DESC"
		}

		query += " LIMIT 100"

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "db error")
			return
		}
		defer rows.Close()

		tasks := []models.Task{}
		for rows.Next() {
			var t models.Task

			var dateFrom, dateTo sql.NullTime
			var status sql.NullString
			var cf, uf sql.NullTime

			if err := rows.Scan(
				&t.ID, &t.Subject, &t.Task,
				&dateFrom, &dateTo, &status,
				&cf, &uf,
			); err != nil {
				log.Printf("getTasks scan: %v", err)
				writeError(w, http.StatusInternalServerError, "scan error")
				return
			}

			if dateFrom.Valid {
				s := dateFrom.Time.Format("2006-01-02")
				t.DateFrom = &s
			}
			if dateTo.Valid {
				s := dateTo.Time.Format("2006-01-02")
				t.DateTo = &s
			}
			if status.Valid {
				s := status.String
				t.Status = &s
			}
			if cf.Valid {
				t.CreatedAt = cf.Time.Format("2006-01-02 15:04:05")
			}
			if uf.Valid {
				t.UpdatedAt = uf.Time.Format("2006-01-02 15:04:05")
			}
			tasks = append(tasks, t)
		}

		if err := rows.Err(); err != nil {
			writeError(w, http.StatusInternalServerError, "rows error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[[]models.Task]{
			Success: true,
			Body:    tasks,
		})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, models.Response[any]{
		Success: false,
		Error:   msg,
		Code:    code,
	})
}
