package main

import (
	"fmt"
	"homework/config"
	"homework/db"
	"homework/handlers/api"
	"homework/handlers/web"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})
	mux.HandleFunc("/api/v1/ping", api.Ping)
	mux.HandleFunc("/api/v1/tasks", api.GetTasks(database))

	mux.HandleFunc("/auth", web.AuthPage)
	mux.HandleFunc("/api/auth/register", web.Register(database))
	mux.HandleFunc("/api/auth/login", web.Login(database))
	mux.HandleFunc("/api/auth/logout", web.Logout(database))
	mux.HandleFunc("/api/auth/me", web.Me(database))

	mux.HandleFunc("/admin",
		web.RequireAuth(database, web.RequireAdmin(web.AdminPage)))

	mux.HandleFunc("/api/admin/tasks",
		web.RequireAuth(database, web.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				web.AdminTasks(database)(w, r)
			case http.MethodPost:
				web.CreateTask(database)(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	mux.HandleFunc("/api/admin/tasks/",
		web.RequireAuth(database, web.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPatch, http.MethodPut:
				web.UpdateTask(database)(w, r)
			case http.MethodDelete:
				web.DeleteTask(database)(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})))

	mux.HandleFunc("/api/admin/subjects",
		web.RequireAuth(database, web.RequireAdmin(web.AdminSubjects)))
	mux.HandleFunc("/api/admin/schedule/day",
		web.RequireAuth(database, web.RequireAdmin(web.AdminScheduleDay)))
	mux.HandleFunc("/api/admin/schedule/next",
		web.RequireAuth(database, web.RequireAdmin(web.AdminScheduleNext)))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	log.Printf("listening on 127.0.0.1:%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+cfg.ServerPort, mux))
}
