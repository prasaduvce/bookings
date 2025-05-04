package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prasaduvce/bookings/internal/config"
	"github.com/prasaduvce/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//what to store in session
	gob.Register(models.Reservation{})

	//change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = false
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.HttpOnly = true

	testApp.Session = session
	
	appConfig = &testApp

	os.Exit(m.Run())
}