package main

import (
	"booking/pkg/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(repo *handlers.Repository) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(writeToConsole)
	mux.Use(Nosurf)
	mux.Use(LoadSession)
	mux.Get("/", repo.Home)
	mux.Get("/about", repo.About)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))
	return mux
}
