package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/config"
)

// create custom functions to be passed into the template library
var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplate sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string) {

	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// else rebuild the template cache
		tc, _ = CreateTemplateCache()
	}

	// get template out of map
	t, ok := tc[tmpl]
	if !ok {
		log.Fatalf("Could not get template %v from template cache", tmpl)
	}

	buf := new(bytes.Buffer)

	// add the template to the buffer
	_ = t.Execute(buf, nil)

	// write content to the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// golang map
	myCache := map[string]*template.Template{}

	// get all files that end with .page.html
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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

		// add template to cache
		myCache[name] = ts
	}
	return myCache, nil
}
