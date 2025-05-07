package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/prasaduvce/bookings/internal/config"
)

var app *config.AppConfig

//setup the app config for the helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status:", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	//get the trace of the error and debug errors
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}