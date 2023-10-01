package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/gkobs/bookings/pkg/config"
	"github.com/gkobs/bookings/pkg/handlers"
	"github.com/gkobs/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	fmt.Println(fmt.Sprintf("Starting application on %s", portNumber))
	err = srv.ListenAndServe()

	log.Fatal(err)
}
