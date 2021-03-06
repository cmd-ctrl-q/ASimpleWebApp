package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/config"
	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/handlers"
	"github.com/cmd-ctrl-q/ASimpleWebApp/pkg/render"
)

const (
	portNumber = ":8080"
)

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {

	// change this to true when in production (there is a better way to do this)
	app.InProduction = false

	// initialize sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour // let sessions last for 24 hours
	// should the cookie persist after the client browser window is closed by the end user?
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // encrypted cookie. false for development

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Staring application on port %s\n", portNumber)

	// router
	router := routes(&app)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: router,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
