package intiator

import (
	"booking/internal/driver"
	"booking/internal/helpers"
	"booking/internal/pkg/config"
	"booking/internal/pkg/handlers"
	"booking/internal/pkg/render"

	red "booking/platform/redis"

	"github.com/go-redis/redis/v8"

	"booking/routes"
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

func Initiate() {
	//gob.Register(models.Reservation{})
	app.IsProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction
	app.Session = session
	//connecting database
	db, err := driver.ConnectSql("host=localhost port=5432 dbname=booking user=postgres password=sm211612")
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
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
	ra := red.NewRedisAdapter(client)
	repo := handlers.NewRepository(&app, db, ra)
	helpers.NewHelpers(&app)
	mux := routes.Routes(repo)
	fmt.Println("server starting at port 8080")
	log.Fatal(http.ListenAndServe(PortNumber, mux))
}
