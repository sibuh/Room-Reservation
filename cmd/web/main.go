package main

import (
	"booking/internal/helpers"
	"booking/internal/pkg/config"
	"booking/internal/pkg/handlers"
	"booking/internal/pkg/models"
	"booking/internal/pkg/render"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const PortNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	gob.Register(models.Reservation{})
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
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errLog
	app.UseCashe = false
	app.TemplateCashe = tc

	render.NewApp(&app)
	repo := handlers.NewRepository(&app)
	helpers.NewHelpers(&app)
	mux := routes(repo)
	fmt.Println("server starting at port 8080")
	log.Fatal(http.ListenAndServe(PortNumber, mux))
}
