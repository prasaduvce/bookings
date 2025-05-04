package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name             string
	url              string
	method           string
	params           []postData
	expectedststcode int
}{
	{"Home", "/", "GET", nil, http.StatusOK},
	{"About", "/about", "GET", nil, http.StatusOK},
	{"GeneralQuarter", "/generals-quarters", "GET", nil, http.StatusOK},
	{"MajorsSuite", "/majors-suite", "GET", nil, http.StatusOK},
	{"SearchAvailability", "/search-availability", "GET", nil, http.StatusOK},
	{"Contact", "/contact", "GET", nil, http.StatusOK},
	{"MakeReservation", "/make-reservation", "GET", nil, http.StatusOK},
	{"ReservationSummary", "/reservation-summary", "GET", nil, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes();

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if (e.method == "GET") {
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Errorf("Error getting %s: %v", e.name, err)
				t.Fatal(err)
			}

			if res.StatusCode != e.expectedststcode {
				t.Errorf("For %s, expected %d but got %d", e.name, e.expectedststcode, res.StatusCode)
			} 
		} else {
			//ts.Client().Post(ts.URL + e.url)
		}
	}
}