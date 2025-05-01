package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prasaduvce/bookings/internal/config"
	"github.com/prasaduvce/bookings/internal/handlers"
	"github.com/prasaduvce/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction // set to true in production
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.HttpOnly = true

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		fmt.Println("Error creating template cache: ", err)
		return
	}

	app.Templates = tc
	app.UseCache = true
	render.SetAppConfig(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Printf("Starting server on port %s", portNumber)
	
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
