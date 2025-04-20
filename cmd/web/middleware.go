package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

func WriteToConsole(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Writing to console")
		handler.ServeHTTP(w, r)
	})
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