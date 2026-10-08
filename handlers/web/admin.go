package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"homework/domain"
	"homework/models"
)

var moscowLoc = time.FixedZone("MSK", 3*60*60)

func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, moscowLoc)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/admin.html")
}

func AdminTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		dateStr := r.URL.Query().Get("date")

		query := `
			SELECT id, subject, task, date_from, date_to, status, instead_of, created_at, updated_at
			FROM tasks
			WHERE 1=1`
		args := []any{}

		if dateStr != "" {
			query += " AND date_to = ?"
			args = append(args, dateStr)
		}

		query += " ORDER BY subject, id"

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			log.Printf("admin tasks query: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}
		defer rows.Close()

		tasks, err := scanTasks(rows)
		if err != nil {
			log.Printf("admin tasks scan: %v", err)
			writeErr(w, http.StatusInternalServerError, "scan error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[[]models.Task]{
			Success: true,
			Body:    tasks,
		})
	}
}

type createTaskReq struct {
	Subject   string  `json:"subject"`
	Task      string  `json:"task"`
	DateFrom  *string `json:"date_from"`
	DateTo    *string `json:"date_to"`
	Status    *string `json:"status"`
	InsteadOf *string `json:"instead_of"`
}

func CreateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req createTaskReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "bad json")
			return
		}

		req.Subject = strings.TrimSpace(req.Subject)
		req.Task = strings.TrimSpace(req.Task)
		if req.Subject == "" || req.Task == "" {
			writeErr(w, http.StatusBadRequest, "subject and task required")
			return
		}
		if _, ok := domain.Names[req.Subject]; !ok {
			writeErr(w, http.StatusBadRequest, "unknown subject")
			return
		}
		if req.InsteadOf != nil {
			s := strings.TrimSpace(*req.InsteadOf)
			if s == "" {
				req.InsteadOf = nil
			} else {
				if _, ok := domain.Names[s]; !ok {
					writeErr(w, http.StatusBadRequest, "unknown instead_of")
					return
				}
				req.InsteadOf = &s
			}
		}

		res, err := db.ExecContext(r.Context(), `
			INSERT INTO tasks (subject, task, date_from, date_to, status, instead_of)
			VALUES (?, ?, ?, ?, ?, ?)`,
			req.Subject, req.Task, req.DateFrom, req.DateTo, req.Status, req.InsteadOf,
		)
		if err != nil {
			log.Printf("create task: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		id, _ := res.LastInsertId()

		t, err := getTaskByID(r.Context(), db, id)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "fetch error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[models.Task]{
			Success: true,
			Body:    t,
		})
	}
}

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch && r.Method != http.MethodPut {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		id, err := idFromPath(r.URL.Path, "/api/admin/tasks/")
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}

		var req createTaskReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "bad json")
			return
		}

		req.Subject = strings.TrimSpace(req.Subject)
		req.Task = strings.TrimSpace(req.Task)
		if req.Subject == "" || req.Task == "" {
			writeErr(w, http.StatusBadRequest, "subject and task required")
			return
		}
		if _, ok := domain.Names[req.Subject]; !ok {
			writeErr(w, http.StatusBadRequest, "unknown subject")
			return
		}
		if req.InsteadOf != nil {
			s := strings.TrimSpace(*req.InsteadOf)
			if s == "" {
				req.InsteadOf = nil
			} else {
				if _, ok := domain.Names[s]; !ok {
					writeErr(w, http.StatusBadRequest, "unknown instead_of")
					return
				}
				req.InsteadOf = &s
			}
		}

		_, err = db.ExecContext(r.Context(), `
			UPDATE tasks
			SET subject = ?, task = ?, date_from = ?, date_to = ?, status = ?, instead_of = ?
			WHERE id = ?`,
			req.Subject, req.Task, req.DateFrom, req.DateTo, req.Status, req.InsteadOf, id,
		)
		if err != nil {
			log.Printf("update task: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		t, err := getTaskByID(r.Context(), db, id)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "fetch error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[models.Task]{
			Success: true,
			Body:    t,
		})
	}
}

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		id, err := idFromPath(r.URL.Path, "/api/admin/tasks/")
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}

		_, err = db.ExecContext(r.Context(), "DELETE FROM tasks WHERE id = ?", id)
		if err != nil {
			log.Printf("delete task: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[any]{Success: true})
	}
}

func AdminSubjects(w http.ResponseWriter, r *http.Request) {
	type item struct {
		Code    string `json:"code"`
		Display string `json:"display"`
	}
	list := make([]item, 0, len(domain.Names))
	for code, name := range domain.Names {
		list = append(list, item{Code: code, Display: name})
	}

	for i := 1; i < len(list); i++ {
		for j := i; j > 0 && list[j].Display < list[j-1].Display; j-- {
			list[j], list[j-1] = list[j-1], list[j]
		}
	}
	writeJSON(w, http.StatusOK, models.Response[[]item]{Success: true, Body: list})
}

func AdminScheduleDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeErr(w, http.StatusBadRequest, "date required")
		return
	}
	d, err := parseDate(dateStr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad date")
		return
	}
	wd := int(d.Weekday())
	if wd == 0 {
		writeJSON(w, http.StatusOK, models.Response[any]{
			Success: true,
			Body: map[string]any{
				"date":     dateStr,
				"day":      -1,
				"day_name": "Вс",
				"lessons":  []any{},
			},
		})
		return
	}
	dayIdx := wd - 1

	lessons := domain.LessonsForDay(dayIdx)
	type lessonItem struct {
		Code    string `json:"code"`
		Display string `json:"display"`
	}
	out := make([]lessonItem, 0, len(lessons))
	for _, code := range lessons {
		out = append(out, lessonItem{
			Code:    code,
			Display: domain.DisplayName(code),
		})
	}

	writeJSON(w, http.StatusOK, models.Response[any]{
		Success: true,
		Body: map[string]any{
			"date":     dateStr,
			"day":      dayIdx,
			"day_name": domain.DayNames[dayIdx],
			"lessons":  out,
		},
	})
}

func AdminScheduleNext(w http.ResponseWriter, r *http.Request) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		writeErr(w, http.StatusBadRequest, "subject required")
		return
	}
	if _, ok := domain.Names[subject]; !ok {
		writeErr(w, http.StatusBadRequest, "unknown subject")
		return
	}

	today := time.Now().In(moscowLoc)
	for i := 0; i < 14; i++ {
		d := today.AddDate(0, 0, i)
		wd := int(d.Weekday())
		if wd == 0 {
			continue
		}
		dayIdx := wd - 1
		for _, code := range domain.LessonsForDay(dayIdx) {
			if code == subject {
				writeJSON(w, http.StatusOK, models.Response[any]{
					Success: true,
					Body: map[string]any{
						"date":     d.Format("2006-01-02"),
						"day":      dayIdx,
						"day_name": domain.DayNames[dayIdx],
					},
				})
				return
			}
		}
	}

	writeErr(w, http.StatusNotFound, "no lesson in next 14 days")
}

func idFromPath(path, prefix string) (int64, error) {
	s := strings.TrimPrefix(path, prefix)
	s = strings.TrimSuffix(s, "/")
	return strconv.ParseInt(s, 10, 64)
}

func scanTasks(rows *sql.Rows) ([]models.Task, error) {
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
			return nil, err
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
	return tasks, rows.Err()
}

func getTaskByID(ctx context.Context, db *sql.DB, id int64) (models.Task, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, subject, task, date_from, date_to, status, instead_of, created_at, updated_at
		FROM tasks WHERE id = ?`, id)
	if err != nil {
		return models.Task{}, err
	}
	defer rows.Close()
	list, err := scanTasks(rows)
	if err != nil {
		return models.Task{}, err
	}
	if len(list) == 0 {
		return models.Task{}, sql.ErrNoRows
	}
	return list[0], nil
}
