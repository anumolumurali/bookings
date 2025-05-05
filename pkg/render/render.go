package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/anumolumurali/bookings/pkg/config"
	"github.com/anumolumurali/bookings/pkg/models"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData To add the default data to each template
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Render Templates using html templates
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var templateCache map[string]*template.Template

	// get the template from the app config
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Couldn't get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Couldn't parse the template from template cache")
	}

	// render the teamplte
	_, err1 := buf.WriteTo(w)
	if err1 != nil {
		log.Println("Error writing template to browser")
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files names *.page.tmpl from ./templates folder
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
