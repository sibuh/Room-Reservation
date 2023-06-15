package render

import (
	"booking/pkg/config"
	"booking/pkg/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

func NewApp(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
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
	fmt.Println("template data:", td)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writting template to browser:", err)
	}

}

func CreateTemplateCash() (map[string]*template.Template, error) {
	myCash := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.html")
	fmt.Println("all pages:", pages)
	if err != nil {
		return myCash, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		fmt.Println("current page:", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("the error source", name, err)
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
