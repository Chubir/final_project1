package main

import (
	"log"
	"myApp/config"
	"myApp/handlers"
	repo "myApp/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	conf := config.Init()
	repo, err := repo.New(conf)
	if err != nil {
		log.Panic(err)
	}
	r := chi.NewRouter()
	r.Post("/api/task", handlers.CreateTaskHandler(repo))
	r.Get("/api/tasks", handlers.ListTaskHandler(repo))
	r.Get("/api/task", handlers.ReadTaskHandler(repo))
	r.Put("/api/task", handlers.UpdateTaskHandler(repo))
	r.Post("/api/task/done", handlers.DoneTaskHandler(repo))
	r.Delete("/api/task", handlers.DeleteTaskHandler(repo))
	r.Get("/api/nextdate", handlers.NextDateHandler())

	webDir := "./web"

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер открылся на порту:%s", conf.Port)
	err = http.ListenAndServe(":"+conf.Port, r)
	if err != nil {
		log.Panic(err)
	}
}
