package routes

import (
	"net/http"
	middle "reservation/internal/middleware"
	"reservation/internal/pkg/handlers"

	"github.com/go-chi/chi"
)

func Routes(repo *handlers.Repository) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middle.Nosurf)
	mux.Get("/", repo.Home)
	mux.Get("/about", repo.About)
	mux.Get("/business", repo.Business)
	mux.Get("/middle", repo.Middle)
	mux.Get("/economic", repo.Economic)
	mux.Get("/contacts", repo.Contacts)

	mux.Get("/reserve", repo.Reserve)
	mux.Post("/postreserve", repo.PostReserve)
	mux.Get("/summary", repo.Summary)
	//user handlers
	mux.Get("/login", repo.Login)
	mux.Get("/loginpage", repo.LoginPage)

	mux.Get("/availability", repo.Availability)
	mux.Post("/checkRooms", repo.CheckAvailability)
	mux.Get("/chooseroom/{id}", repo.Chooseroom)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
