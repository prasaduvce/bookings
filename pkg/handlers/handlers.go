package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prasaduvce/bookings/pkg/config"
	"github.com/prasaduvce/bookings/pkg/models"
	"github.com/prasaduvce/bookings/pkg/render"
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

func (Repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	//log.Println("Inside Home")
	fmt.Fprintf(w, "This is home page")
}

func (Repo *Repository) HomeHtml(w http.ResponseWriter, r *http.Request) {
	//log.Println("Inside HomeHtml")
	remoteIP := r.RemoteAddr
	//log.Println("remoteIP inside Home ", remoteIP)
	Repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderHtml(w, "home.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) AboutHtml(w http.ResponseWriter, r *http.Request) {
	//log.Println("Inside AboutHtml")
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, World!"

	// Get the remote IP address from the session
	remoteIP := Repo.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//log.Println("remoteIP inside about ", remoteIP)
	render.RenderHtml(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	}, r)
}

func (Repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "make-reservation.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "search-availability.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) PostSearch(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Posted to search availability, start: %s, end: %s", start, end)))
}

type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (Repo *Repository) AvilabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := JsonResponse{
		OK:      false,
		Message: "Available!",
	}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (Repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "contact.page.tmpl", &models.TemplateData{}, r)
}
