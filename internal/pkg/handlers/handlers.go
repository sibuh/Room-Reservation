package handlers

import (
	"booking/internal/pkg/config"
	"booking/internal/pkg/models"
	"booking/internal/pkg/render"
	"encoding/json"
	"log"
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
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{MapString: stringMap,
		Flash: "abu",
	})

}
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "This is template data"
	remoteIP := repo.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remoteIP"] = remoteIP
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
		log.Println(err)
	}
	w.Write(jresponse)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	// render.RenderTemplate(w, "availability.page.html", &models.TemplateData{})
}
