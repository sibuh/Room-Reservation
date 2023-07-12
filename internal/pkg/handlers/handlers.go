package handlers

import (
	"booking/internal/driver"
	"booking/internal/forms"
	"booking/internal/helpers"
	"booking/internal/pkg/config"
	"booking/internal/pkg/models"
	"booking/internal/pkg/render"
	"booking/internal/repository"
	"booking/internal/repository/dbrepo"
	"booking/platform"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type Repository struct {
	App   *config.AppConfig
	DB    repository.DatabaseRepo
	Redis platform.RedisInterface
}

func NewRepository(a *config.AppConfig, db *driver.DB, red platform.RedisInterface) *Repository {
	return &Repository{
		App:   a,
		DB:    dbrepo.NewPostgresDbRepo(db.SQL, a),
		Redis: red,
	}
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}
func (repo *Repository) Middle(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "middle.page.html", &models.TemplateData{})
}
func (repo *Repository) Economic(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "economic.page.html", &models.TemplateData{})
}
func (repo *Repository) Reserve(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "reserve.page.html", &models.TemplateData{})
}

func (repo *Repository) PostReserve(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.InfoLog.Println(err)
		return
	}
	res, err := repo.Redis.GetFromRedis(context.Background(), "reservation")
	if err != nil {
		helpers.ServerError(w, err)
	}

	var v models.Reservation
	err = json.Unmarshal([]byte(res), &v)
	if err != nil {
		helpers.ServerError(w, err)
	}
	fmt.Println("this is read from redis,in post reserve", v)
	var reservation = models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: v.StartDate,
		EndDate:   v.EndDate,
		RoomID:    v.RoomID,
	}

	form := forms.New(r.PostForm)
	form.Required()
	data := make(map[string]interface{})
	data["reservation"] = reservation
	if r.Form.Get("first_name") == "" {
		form.Errors.Add("first_name", "this filed is mandatory")
	}
	fmt.Println(form.Errors.Get("first_name"))
	if !form.Valid() {
		for _, e := range form.Errors {
			fmt.Println(e)
		}
		render.Template(w, r, "reserved.page.html",
			&models.TemplateData{
				Form: form,
				Data: data,
			})
	}
	resID, err := repo.DB.MakeReservation(reservation)
	if err != nil {
		fmt.Println("Reservation failed", err)
		helpers.ServerError(w, err)
		return
	}
	rr := models.RoomRestriction{
		RoomID:        v.RoomID,
		RestrictionID: 1,
		ReservationID: resID,
		StartDate:     v.StartDate,
		EndDate:       v.EndDate,
	}
	err = repo.DB.InsertRoomRestriction(rr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fmt.Println("\u001b[31m;1m", "Set reservation data to redis")
	repo.Redis.SetToRedis(context.Background(), "reservation", reservation)
	msg := models.EmailData{
		From:     "sibuh@gmail.com",
		To:       reservation.Email,
		Subject:  "Confirmation Email",
		Content:  fmt.Sprintf("This email is to confirmation that you have reserved %s from %s to %s", reservation.Room.RoomName, reservation.StartDate, reservation.EndDate),
		Template: "basic.html",
	}
	repo.App.ErrorLog.Println(msg)
	repo.App.EmailChan <- msg
	fmt.Println("sent msg to channel")
	http.Redirect(w, r, "/summary", http.StatusSeeOther)
}
func (repo *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	res, err := repo.Redis.GetFromRedis(context.Background(), "reservation")
	fmt.Println(res, "read from the redis in summary")
	if err != nil {
		helpers.ServerError(w, err)
	}
	fmt.Println(res)
	data["reserved"] = res
	render.Template(w, r, "summary.page.html", &models.TemplateData{Data: data})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "availability.page.html", &models.TemplateData{})
}
func (repo *Repository) Business(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "business.page.html", &models.TemplateData{})
}
func (repo *Repository) Contacts(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contacts.page.html", &models.TemplateData{})
}
func (repo *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var layout = "2006-01-02"
	start := r.PostFormValue("start_date")
	end := r.PostFormValue("end_date")
	roomID := r.PostFormValue("room_id")

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if roomID != "" {
		rID, err := strconv.Atoi(roomID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		reserved, err := repo.DB.SearchAvailabilityByRoomID(rID, startDate, endDate)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		if reserved {
			data := make(map[string]interface{})
			data["reserved"] = "The room is already reserved"
			render.Template(w, r, "available.page.html", &models.TemplateData{
				Data: data,
			})
		} else {
			data := make(map[string]interface{})
			data["reserved"] = "The room is free for reservation"
			render.Template(w, r, "available.page.html", &models.TemplateData{
				Data: data,
			})
		}
	} else {
		rooms, err := repo.DB.SearchAvailableRooms(startDate, endDate)
		fmt.Println(rooms)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		if len(rooms) == 0 {
			w.Write([]byte("No available rooms"))
			http.Redirect(w, r, "/availability", http.StatusSeeOther)
			return
		} else {
			res := &models.Reservation{
				StartDate: startDate,
				EndDate:   endDate,
			}

			err = repo.Redis.SetToRedis(context.Background(), "reservation", res)
			if err != nil {
				helpers.ServerError(w, err)
			}
			data := make(map[string]interface{})
			data["rooms"] = rooms
			render.Template(w, r, "availablerooms.page.html", &models.TemplateData{
				Data: data,
			})
		}
	}
}
func (repo *Repository) Chooseroom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	fmt.Println("this is room id", roomID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	value, err := repo.Redis.GetFromRedis(context.Background(), "reservation")
	if err != nil {
		helpers.ServerError(w, err)
	}
	fmt.Println(value)
	var v models.Reservation
	err = json.Unmarshal([]byte(value), &v)
	if err != nil {
		helpers.ServerError(w, err)
	}
	v.RoomID = roomID
	fmt.Println(v, "roomID added")
	err = repo.Redis.SetToRedis(context.Background(), "reservation", v)
	data := map[string]interface{}{}
	data["choose"] = v
	if err != nil {
		helpers.ServerError(w, err)
	}

	render.Template(w, r, "reserve.page.html", &models.TemplateData{
		Data: data,
	})

}
