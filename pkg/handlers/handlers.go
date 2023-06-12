package handlers

import (
	"net/http"
	"webApp/pkg/config"
	"webApp/pkg/models"
	"webApp/pkg/render"
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
