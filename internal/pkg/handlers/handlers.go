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
	validation "github.com/go-ozzo/ozzo-validation"
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
		return
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
	fmt.Println("\u001b[31m", "Set session data to redis")
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	err = repo.Redis.SetToRedis(ctx, "reservation", reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// msg := models.EmailData{
	// 	From:     "sibuh@gmail.com",
	// 	To:       reservation.Email,
	// 	Subject:  "Confirmation Email",
	// 	Content:  fmt.Sprintf("This email is to confirmation that you have reserved %s from %s to %s", reservation.Room.RoomName, reservation.StartDate, reservation.EndDate),
	// 	Template: "basic.html",
	// }
	// repo.App.EmailChan <- msg
	//fmt.Println("sent msg to channel in Reservation")
	http.Redirect(w, r, "/summary", http.StatusSeeOther)
}
func (repo *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	res, err := repo.Redis.GetFromRedis(context.Background(), "reservation")
	fmt.Println(res, "read from the redis in summary")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var reservation models.Reservation
	err = json.Unmarshal([]byte(res), &reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["reserved"] = reservation
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
	var Input = models.RoomAvailabilityRequest{
		StartDate: start,
		EndDate:   end,
	}
	err := validation.ValidateStruct(&Input,
		validation.Field(&Input.StartDate, validation.Required.Error("SatrtDate is Required")),
		validation.Field(&Input.EndDate, validation.Required.Error("EndDate is Required"), validation.Date(layout).Error("wrong layout of date")),
	)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
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
	rooms, err := repo.DB.SearchAvailableRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		w.Write([]byte("No available rooms"))
		return
	} else {
		res := &models.Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}
		data := make(map[string]interface{})
		data["rooms"] = rooms
		render.Template(w, r, "availablerooms.page.html", &models.TemplateData{
			Data: data,
		})
		err = repo.Redis.SetToRedis(context.Background(), "reservation", res)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

	}

}
func (repo *Repository) Chooseroom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	fmt.Println("This is room id", roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	value, err := repo.Redis.GetFromRedis(context.Background(), "reservation")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fmt.Println(value)
	var res models.Reservation
	err = json.Unmarshal([]byte(value), &res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = roomID
	fmt.Println(res, "roomID added")
	err = repo.Redis.SetToRedis(context.Background(), "reservation", res)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	data := make(map[string]interface{})
	data["choosen"] = res
	render.Template(w, r, "reserve.page.html", &models.TemplateData{Data: data})

}
func (repo *Repository) AddRooms(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var req models.AddRoomRequest
	req.RoomNumber = r.Form.Get("room_number")
	req.RoomType = r.Form.Get("room_type")
	err = validation.ValidateStruct(&req,
		validation.Field(&req.RoomNumber, validation.Required.Error("room_number is required")),
		validation.Field(&req.RoomType, validation.Required.Error("room_type is required")))
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	err = repo.DB.InsertRooms(req)
	if err != nil {
		helpers.ServerError(w, err)
		return
	} else {
		render.Template(w, r, "insertroom.page.html", &models.TemplateData{})
	}

}
func (repo *Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	var req models.LoginRequest
	req.Email = r.Form.Get("email")
	req.Password = r.Form.Get("pasword")
	err = validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required.Error("email is required")),
		validation.Field(&req.Password, validation.Required.Error("password is required"), validation.Length(8, 8)))
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	tokenString, err := repo.DB.Login(req)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	err = repo.Redis.SetToRedis(context.Background(), "token", tokenString)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	cookie := http.Cookie{
		Name:  "sesion_token",
		Value: tokenString,
	}
	http.SetCookie(w, &cookie)

}
