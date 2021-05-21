package main

import (
	"log"
	"net/http"

	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/config"
	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/handlers"
	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/render"
)

func main() {

	var app config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}

	// store the template cache in the site-wide config struct
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	// call handlers
	http.HandleFunc("/", handlers.Repo.Home)

	http.HandleFunc("/about", handlers.Repo.About)

	_ = http.ListenAndServe(":8080", nil)
}
