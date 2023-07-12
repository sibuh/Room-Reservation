package render

import (
	"booking/internal/pkg/config"
	"booking/internal/pkg/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{"format": func(t time.Time) string {
	return t.Format("2006-01-02")
}}

var app *config.AppConfig

func NewApp(a *config.AppConfig) {
	app = a
}
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	td = AddDefaultData(td, r)
	var tc map[string]*template.Template
	if app.UseCashe {
		tc = app.TemplateCashe
	} else {
		tc, _ = CreateTemplateCash()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatalf("no template corresponding to the %s", tmpl)
	}
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writting template to browser:", err)
	}

}

func CreateTemplateCash() (map[string]*template.Template, error) {
	myCash := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCash, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCash, err
		}
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCash, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCash, err
			}
		}
		myCash[name] = ts

	}
	return myCash, nil
}
