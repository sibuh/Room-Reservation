package middle

import (
	"booking/internal/pkg/config"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager

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
