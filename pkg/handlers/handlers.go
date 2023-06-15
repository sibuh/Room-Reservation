package handlers

import (
	"booking/pkg/config"
	"booking/pkg/models"
	"booking/pkg/render"
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
	stringMap["test"] = "This is template data"
	remoteIP := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{MapString: stringMap,
		Flash: "abu",
	})

}
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "This is template data"
	remoteIP := repo.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remoteIP"] = remoteIP
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{MapString: stringMap})
}
func (repo *Repository) Middle(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "middle.page.html", &models.TemplateData{})
}
func (repo *Repository) Economic(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "economic.page.html", &models.TemplateData{})
}
func (repo *Repository) Reserve(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "reserve.page.html", &models.TemplateData{})
}
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "availability.page.html", &models.TemplateData{})
}
func (repo *Repository) Business(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "availability.page.html", &models.TemplateData{})
}
func (repo *Repository) Contacts(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "contacts.page.html", &models.TemplateData{})
}
