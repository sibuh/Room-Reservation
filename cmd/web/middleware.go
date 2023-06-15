package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func writeToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("you hit the page ")
		next.ServeHTTP(w, r)
	})
}
func Nosurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}
func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
