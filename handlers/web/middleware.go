package web

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"homework/models"
)

type ctxKey string

const userCtxKey ctxKey = "user"

func RequireAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err != nil || c.Value == "" {
			writeErr(w, http.StatusUnauthorized, "no session")
			return
		}

		var (
			id        int64
			username  string
			role      string
			createdAt time.Time
		)
		err = db.QueryRowContext(r.Context(),
			"SELECT id, username, role, created_at FROM users WHERE session_str = ?",
			c.Value,
		).Scan(&id, &username, &role, &createdAt)

		if err == sql.ErrNoRows {
			writeErr(w, http.StatusUnauthorized, "invalid session")
			return
		}
		if err != nil {
			log.Printf("requireAuth: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		u := models.User{
			ID:        id,
			Username:  username,
			Role:      role,
			CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
		}
		ctx := context.WithValue(r.Context(), userCtxKey, u)
		next(w, r.WithContext(ctx))
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(userCtxKey).(models.User)
		if !ok {
			writeErr(w, http.StatusUnauthorized, "no user")
			return
		}
		if u.Role != "admin" {
			writeErr(w, http.StatusForbidden, "admin only")
			return
		}
		next(w, r)
	}
}

func UserFromCtx(r *http.Request) (models.User, bool) {
	u, ok := r.Context().Value(userCtxKey).(models.User)
	return u, ok
}
