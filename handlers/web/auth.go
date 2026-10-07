package web

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"homework/models"
)

const sessionCookieName = "session_str"
const sessionTTL = 30 * 24 * time.Hour

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, models.Response[any]{
		Success: false,
		Error:   msg,
		Code:    code,
	})
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func setSessionCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func AuthPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/auth.html")
}

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req registerReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "bad json")
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Key = strings.TrimSpace(req.Key)

		if len(req.Username) < 3 {
			writeErr(w, http.StatusBadRequest, "username too short")
			return
		}
		if len(req.Password) < 6 {
			writeErr(w, http.StatusBadRequest, "password too short")
			return
		}

		parts := strings.SplitN(req.Key, ":", 2)
		if len(parts) != 2 {
			writeErr(w, http.StatusBadRequest, "invalid key format")
			return
		}

		var keyID int64
		if _, err := fmt.Sscanf(parts[0], "%d", &keyID); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid key id")
			return
		}
		originalKey := parts[1]

		var keyHash, keyStatus string
		err := db.QueryRowContext(r.Context(),
			"SELECT key_hash, status FROM `keys` WHERE id = ?", keyID,
		).Scan(&keyHash, &keyStatus)

		if err == sql.ErrNoRows {
			writeErr(w, http.StatusBadRequest, "key not found")
			return
		}
		if err != nil {
			log.Printf("register key lookup: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}
		if keyStatus != "ready" {
			writeErr(w, http.StatusBadRequest, "key already used")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(keyHash), []byte(originalKey)); err != nil {
			writeErr(w, http.StatusBadRequest, "key mismatch")
			return
		}

		pwHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("register bcrypt: %v", err)
			writeErr(w, http.StatusInternalServerError, "hash error")
			return
		}

		session, err := randomHex(32)
		if err != nil {
			log.Printf("register session: %v", err)
			writeErr(w, http.StatusInternalServerError, "session error")
			return
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			log.Printf("register tx: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}
		defer tx.Rollback()

		res, err := tx.ExecContext(r.Context(),
			"INSERT INTO users (username, password_hash, role, session_str) VALUES (?, ?, 'user', ?)",
			req.Username, string(pwHash), session,
		)
		if err != nil {
			log.Printf("register insert user: %v", err)
			writeErr(w, http.StatusConflict, "username taken")
			return
		}

		userID, _ := res.LastInsertId()

		if _, err := tx.ExecContext(r.Context(),
			"UPDATE `keys` SET used_by = ?, status = 'used' WHERE id = ? AND status = 'ready'",
			userID, keyID,
		); err != nil {
			log.Printf("register update key: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		if err := tx.Commit(); err != nil {
			log.Printf("register commit: %v", err)
			writeErr(w, http.StatusInternalServerError, "commit error")
			return
		}

		setSessionCookie(w, session)

		writeJSON(w, http.StatusOK, models.Response[models.User]{
			Success: true,
			Body: models.User{
				ID:       userID,
				Username: req.Username,
				Role:     "user",
			},
		})
	}
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "bad json")
			return
		}
		req.Username = strings.TrimSpace(req.Username)

		var (
			id        int64
			pwHash    string
			role      string
			createdAt time.Time
		)
		err := db.QueryRowContext(r.Context(),
			"SELECT id, password_hash, role, created_at FROM users WHERE username = ?",
			req.Username,
		).Scan(&id, &pwHash, &role, &createdAt)

		if err == sql.ErrNoRows {
			writeErr(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if err != nil {
			log.Printf("login lookup: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(req.Password)); err != nil {
			writeErr(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		session, err := randomHex(32)
		if err != nil {
			log.Printf("login session: %v", err)
			writeErr(w, http.StatusInternalServerError, "session error")
			return
		}

		if _, err := db.ExecContext(r.Context(),
			"UPDATE users SET session_str = ? WHERE id = ?", session, id,
		); err != nil {
			log.Printf("login update session: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		setSessionCookie(w, session)

		writeJSON(w, http.StatusOK, models.Response[models.User]{
			Success: true,
			Body: models.User{
				ID:        id,
				Username:  req.Username,
				Role:      role,
				CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		c, err := r.Cookie(sessionCookieName)
		if err == nil && c.Value != "" {
			if _, err := db.ExecContext(r.Context(),
				"UPDATE users SET session_str = NULL WHERE session_str = ?", c.Value,
			); err != nil {
				log.Printf("logout: %v", err)
			}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		writeJSON(w, http.StatusOK, models.Response[any]{Success: true})
	}
}

func Me(db *sql.DB) http.HandlerFunc {
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
			log.Printf("me: %v", err)
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}

		writeJSON(w, http.StatusOK, models.Response[models.User]{
			Success: true,
			Body: models.User{
				ID:        id,
				Username:  username,
				Role:      role,
				CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}
