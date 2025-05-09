package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/prasaduvce/bookings/internal/config"
	"github.com/prasaduvce/bookings/internal/models"
	"github.com/prasaduvce/bookings/internal/render"
	// Removed cyclic dependency by not importing handlers
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var infoLog *log.Logger
var errorLog *log.Logger

func getRoutes() http.Handler {
	// what to store in session
	gob.Register(models.Reservation{})

	// change this to true when in production
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

	tc, err := CreateTestTemplateCache()

	if err != nil {
		fmt.Println("Error creating template cache: ", err)
	}

	app.Templates = tc
	app.UseCache = true
	render.NewTemplates(&app)
	// Avoid using handlers directly to prevent cyclic dependency
	repo := NewRepo(&app)
	NewHandlers(repo)
	//render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	//mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.HomeHtml)
	mux.Get("/about", Repo.AboutHtml)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	
	mux.Get("/search-availability", Repo.Search)
	mux.Post("/search-availability", Repo.PostSearch)
	mux.Post("/search-availability-json", Repo.AvilabilityJson)
	
	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

//NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure: app.InProduction,
		Path: "/",
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//loads the session and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)

	myCache := map[string]*template.Template{}

	//get all the files names *.page.tmpl from templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	//range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		//parse the page and base layout
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println("Error parsing template: ", err)
			return myCache, err
		}

		//parse the base layout
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			log.Println("Error parsing base layout: ", err)
			return myCache, err
		}
		//parse the partials
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				log.Println("Error parsing base layout: ", err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}