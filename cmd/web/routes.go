package main

import (
	"net/http"
	"webApp/pkg/handlers"

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
	return mux
}
