package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prasaduvce/bookings/internal/config"
	"github.com/prasaduvce/bookings/internal/handlers"
)

func routes(appConfig *config.AppConfig) http.Handler {
	appConfig.UseCache = false
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	//mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.HomeHtml)
	mux.Get("/about", handlers.Repo.AboutHtml)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	
	mux.Get("/search-availability", handlers.Repo.Search)
	mux.Post("/search-availability", handlers.Repo.PostSearch)
	mux.Post("/search-availability-json", handlers.Repo.AvilabilityJson)
	
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
