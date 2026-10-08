package web

import (
	"database/sql"
	"log"
	"net/http"

	"homework/models"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/template/admin.html")
}

func AdminTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT id, subject, task, date_from, date_to, status, instead_of, created_at, updated_at
			FROM tasks
			ORDER BY id DESC
			LIMIT 500
		`)
		if err != nil {
			log.Printf("admin tasks query: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
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
				log.Printf("admin tasks scan: %v", err)
				writeErr(w, http.StatusInternalServerError, "scan error")
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
			log.Printf("admin tasks rows: %v", err)
			writeErr(w, http.StatusInternalServerError, "rows error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[[]models.Task]{
			Success: true,
			Body:    tasks,
		})
	}
}
