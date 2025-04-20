package handlers

import (
	"fmt"
	"net/http"

	"github.com/prasaduvce/bookings/pkg/config"
	"github.com/prasaduvce/bookings/pkg/render"
	"github.com/prasaduvce/bookings/pkg/models"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (Repo *Repository) Home(w http.ResponseWriter, t *http.Request) {
	//log.Println("Inside Home")
	fmt.Fprintf(w, "This is home page")
}

func (Repo *Repository) HomeHtml(w http.ResponseWriter, t *http.Request) {
	//log.Println("Inside HomeHtml")
	remoteIP := t.RemoteAddr
	//log.Println("remoteIP inside Home ", remoteIP)
	Repo.App.Session.Put(t.Context(), "remote_ip", remoteIP)
	render.RenderHtml(w, "home.page.tmpl", &models.TemplateData{})
}

func (Repo *Repository) AboutHtml(w http.ResponseWriter, t *http.Request) {
	//log.Println("Inside AboutHtml")
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, World!"

	// Get the remote IP address from the session
	remoteIP := Repo.App.Session.GetString(t.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//log.Println("remoteIP inside about ", remoteIP)
	render.RenderHtml(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	} )
}