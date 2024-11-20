package helpers

import (
	"fmt"
	"net/http"
	"reservation/internal/pkg/config"
	"runtime/debug"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error wiht status ", status)
	http.Error(w, http.StatusText(status), status)
}
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(500), 500)
}
