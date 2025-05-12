package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prasaduvce/bookings/internal/config"
	driver "github.com/prasaduvce/bookings/internal/drivers"
	"github.com/prasaduvce/bookings/internal/handlers"
	"github.com/prasaduvce/bookings/internal/helpers"
	"github.com/prasaduvce/bookings/internal/models"
	"github.com/prasaduvce/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		fmt.Println("Error starting application: ", err)
		return
	}
	defer db.SQL.Close()

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

func run() (*driver.DB,  error) {

	//what to store in session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})

	//change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction // set to true in production
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.HttpOnly = true

	app.Session = session

	//connect to Database
	log.Println("Connecting to Database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=renuka.bm password=")
	if err != nil {
		log.Fatal("Cannot connect to Database !!!", err)
	}
	log.Println("connected to Database")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		fmt.Println("Error creating template cache: ", err)
		return nil, err
	}


	app.Templates = tc
	app.UseCache = true
	render.NewRenderer(&app)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	return db, nil;
}