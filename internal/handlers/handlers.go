package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prasaduvce/bookings/internal/config"
	driver "github.com/prasaduvce/bookings/internal/drivers"
	"github.com/prasaduvce/bookings/internal/forms"
	"github.com/prasaduvce/bookings/internal/helpers"
	"github.com/prasaduvce/bookings/internal/models"
	"github.com/prasaduvce/bookings/internal/render"
	"github.com/prasaduvce/bookings/internal/repository"
	"github.com/prasaduvce/bookings/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB repository.DataBaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostGresRepo(db.SQL, a),
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
	render.RenderHtml(w, "home.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) AboutHtml(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "about.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderHtml(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

// PostReservation handles the post request for reservation form
func (Repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//log.Println("Error parsing form: ", err)
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderHtml(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}
	Repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
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
		helpers.ServerError(w, err)
		return
		//log.Println("Error marshalling JSON: ", err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (Repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	// Get the reservation from the session
	reservation, ok := Repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		Repo.App.ErrorLog.Println("Cannot get reservation from session")
		//log.Println("Cannot get reservation from session")
		Repo.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Clear the reservation from the session
	Repo.App.Session.Remove(r.Context(), "reservation")
	render.RenderHtml(w, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: map[string]interface{}{
			"reservation": reservation,
		},
	}, r)
}
