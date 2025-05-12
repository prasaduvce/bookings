package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	render.Template(w, "home.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) AboutHtml(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "about.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := Repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Cannot get reservation from session"))
		return	
	}

	room, err := Repo.DB.GetRoomById(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	Repo.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	}, r)
}

// PostReservation handles the post request for reservation form
func (Repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := Repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Can't get from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		//log.Println("Error parsing form: ", err)
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}
	newReservationId, err := Repo.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	Repo.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction {
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	Repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	Repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (Repo *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.tmpl", &models.TemplateData{}, r)
}

func (Repo *Repository) PostSearch(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := Repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	for _, room := range rooms {
		Repo.App.InfoLog.Println("ROOM:", room.ID, room.RoomName)
	}

	if len(rooms) == 0 {
		//no availability
		Repo.App.InfoLog.Println("There is no availability")
		Repo.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}

	Repo.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}

type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	RoomID string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

func (Repo *Repository) AvilabilityJson(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	available, err := Repo.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomID)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := JsonResponse{
		OK:      available,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (Repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl", &models.TemplateData{}, r)
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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: map[string]interface{}{
			"reservation": reservation,
		},
		StringMap: stringMap,
	}, r)
}

func (Repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	Repo.App.Session.Get(r.Context(), "reservation")

	res, ok := Repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return	
	}
	res.RoomID = roomID
	Repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//takes URL parameters, takes user to make reservation screen
func (Repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// id, s, e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	res.RoomID = roomID

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.StartDate = startDate
	res.EndDate = endDate

	room, err := Repo.DB.GetRoomById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	Repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}