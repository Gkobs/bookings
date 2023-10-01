package render

import (
	"bytes"
	"github.com/gkobs/bookings/pkg/models"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gkobs/bookings/pkg/config"
)

var app *config.AppConfig

// NewTemplates sets config for templates package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// addDefaultData adds default data to the template data
func addDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Template renders a template in ResponseWriter
func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	err := t.Execute(buf, addDefaultData(td))
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
