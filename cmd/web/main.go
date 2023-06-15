package main

import (
	"booking/pkg/config"
	"booking/pkg/handlers"
	"booking/pkg/render"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const PortNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.IsProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction
	app.Session = session
	tc, err := render.CreateTemplateCash()
	if err != nil {
		log.Fatal("can not template cashe")
	}
	app.UseCashe = false
	app.TemplateCashe = tc

	render.NewApp(&app)
	repo := handlers.NewRepository(&app)
	mux := routes(repo)
	fmt.Println("server starting at port 8080")
	log.Fatal(http.ListenAndServe(PortNumber, mux))
}
