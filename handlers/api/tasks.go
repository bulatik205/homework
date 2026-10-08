package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"homework/domain"
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
		dateStr := q.Get("date")

		query := `SELECT id, subject, task, date_from, date_to, status, instead_of, created_at, updated_at
		          FROM tasks WHERE 1=1`
		args := []any{}

		if subject != "" {
			query += " AND subject = ?"
			args = append(args, subject)
		}

		if dateStr != "" {
			query += " AND date_to = ?"
			args = append(args, dateStr)
		}

		switch recency {
		case "last":
			query += " ORDER BY id DESC"
		default:
			query += " ORDER BY subject, id"
		}

		query += " LIMIT 200"

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
			var status, insteadOf sql.NullString
			var cf, uf sql.NullTime

			if err := rows.Scan(
				&t.ID, &t.Subject, &t.Task,
				&dateFrom, &dateTo, &status, &insteadOf,
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
			if insteadOf.Valid {
				s := insteadOf.String
				t.InsteadOf = &s
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
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type item struct {
		Code    string `json:"code"`
		Display string `json:"display"`
	}

	list := make([]item, 0, len(domain.Names))
	for code, name := range domain.Names {
		list = append(list, item{Code: code, Display: name})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Display < list[j].Display
	})

	writeJSON(w, http.StatusOK, models.Response[[]item]{
		Success: true,
		Body:    list,
	})
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
