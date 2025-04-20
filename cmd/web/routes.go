package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prasaduvce/bookings/pkg/config"
	"github.com/prasaduvce/bookings/pkg/handlers"
)

func routes(appConfig *config.AppConfig) http.Handler{
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	//mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.HomeHtml)
	mux.Get("/about", handlers.Repo.AboutHtml)
	return mux
}