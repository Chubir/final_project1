package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
	_ "modernc.org/sqlite"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "scheduler.db"
	}

	var install bool
	if _, err := os.Stat(dbfile); err != nil {
		install = true
	}

	if install {
		instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			log.Fatal(err)
		}

		fSrc, err := (&file.File{}).Open("./migrations")
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithInstance("file", fSrc, "scheduler", instance)
		if err != nil {
			log.Fatal(err)
		}

		// modify for Down
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}
	r := chi.NewRouter()
	r.Post("/api/task", CreateTaskHandler(db))
	r.Get("/api/tasks", ListTaskHandler(db))
	r.Get("/api/task", ReadTaskHandler(db))
	r.Put("/api/task", UpdateTaskHandler(db))
	r.Post("/api/task/done", DoneTaskHandler(db))
	r.Delete("/api/task", DeleteTaskHandler(db))

	webDir := "./web"

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер открылся на порту:%s", port)
	http.ListenAndServe(":"+port, r)
}
