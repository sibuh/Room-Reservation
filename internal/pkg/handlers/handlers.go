package handlers

import (
	"booking/internal/forms"
	"booking/internal/helpers"
	"booking/internal/pkg/config"
	"booking/internal/pkg/models"
	"booking/internal/pkg/render"
	"encoding/json"
	"fmt"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{MapString: stringMap})

}
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{MapString: stringMap})
}
func (repo *Repository) Middle(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "middle.page.html", &models.TemplateData{})
}
func (repo *Repository) Economic(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "economic.page.html", &models.TemplateData{})
}
func (repo *Repository) Reserve(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "reserve.page.html", &models.TemplateData{})
}

func (repo *Repository) PostReserve(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.InfoLog.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName:   r.Form.Get("first_name"),
		LastName:    r.Form.Get("last_name"),
		PhoneNumber: r.Form.Get("phone_number"),
		Email:       r.Form.Get("email"),
	}
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	data["reservation"] = reservation
	if r.Form.Get("first_name") == "" {
		form.Errors.Add("first_name", "this filed is mandatory")
	}
	fmt.Println(form.Errors.Get("first_name"))
	fmt.Println(form.Valid())
	if !form.Valid() {
		render.RenderTemplate(w, r, "reserved.page.html",
			&models.TemplateData{
				Form: form,
				Data: data,
			})
	}
	render.RenderTemplate(w, r, "reserve.page.html", &models.TemplateData{})
}
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "availability.page.html", &models.TemplateData{})
}
func (repo *Repository) Business(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "business.page.html", &models.TemplateData{})
}
func (repo *Repository) Contacts(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contacts.page.html", &models.TemplateData{})
}
func (repo *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.PostFormValue("start")
	end := r.PostFormValue("end")
	response := struct {
		Start string
		End   string
	}{
		Start: start,
		End:   end,
	}
	jresponse, err := json.MarshalIndent(response, "", "")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Write(jresponse)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	// render.RenderTemplate(w, "availability.page.html", &models.TemplateData{})
}
