package intiator

import (
	"booking/internal/driver"
	"booking/internal/helpers"
	"booking/internal/pkg/config"
	"booking/internal/pkg/handlers"
	"booking/internal/pkg/models"
	"booking/internal/pkg/render"
	"booking/routes"
	"io/ioutil"
	"strings"

	red "booking/platform/redis"

	"github.com/go-redis/redis/v8"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mail "github.com/xhit/go-simple-mail"

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
	//
	ListenToEmailChan()
	fmt.Println("created message channel")
	app.EmailChan = make(chan models.EmailData)
	fmt.Println("stopped here")
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	time.Sleep(1 * time.Second)
	//create template cache
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
func SendEmail(msg models.EmailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 5 * time.Second
	server.SendTimeout = 5 * time.Second
	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	basic, err := ioutil.ReadFile(fmt.Sprintf("./email-template/%s", msg.Template))
	//fmt.Println(string(basic))
	if err != nil {
		app.ErrorLog.Println(err)
	}
	if string(basic) == "" {
		email.SetBody(mail.TextPlain, msg.Content)
	} else {
		bodyToSend := string(basic)
		bodyString := strings.Replace(bodyToSend, "[%BODY%]", msg.Content, 1)
		email.SetBody(mail.TextHTML, bodyString)
	}
	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
	}
}
func ListenToEmailChan() {
	go func() {
		msg := <-app.EmailChan
		SendEmail(msg)
	}()
}
