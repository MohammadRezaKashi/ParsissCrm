package main

import (
	"ParsissCrm/internal/config"
	"ParsissCrm/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/."))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	fileServer = http.FileServer(http.Dir("./node_modules/."))
	mux.Handle("/node_modules/*", http.StripPrefix("/node_modules", fileServer))

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/report", handlers.Repo.Report)

	mux.Get("/about", handlers.Repo.About)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/report/add-new-report", handlers.Repo.AddNewReport)

	mux.Post("/report/post-add-new-report", handlers.Repo.PostAddNewReport)

	mux.Get("/report/detail/{id}/show", handlers.Repo.ShowDetail)

	mux.Post("/report/post-update-report/{id}", handlers.Repo.PostUpdateReport)

	mux.Post("/report/filters/show", handlers.Repo.ShowFilters)

	return mux
}
