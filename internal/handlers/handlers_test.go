package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"PostSearch", "/search-availability", "POST", []postData{
		{"start", "2025-10-01"}, 
		{"end", "2025-10-02"}}, 
		http.StatusOK,
	},
	{"PostSearchAvailJson", "/search-availability-json", "POST", []postData{
		{"start", "2025-10-01"}, 
		{"end", "2025-10-02"}}, 
		http.StatusOK,
	},
	{"MakeReservation", "/make-reservation", "POST", []postData{
		{"first_name", "John"},
		{"last_name", "Doe"},
		{"email:", "test@test.com"},
		{"phone", "1234567890"},
	}, http.StatusOK},
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
			values := url.Values {}
				for _, param := range e.params {
					values.Add(param.key, param.value)
			}
			//values.Add("csrf_token", "test")
			res, err := ts.Client().PostForm(ts.URL + e.url, values)
			if err != nil {
				t.Errorf("Error getting %s: %v", e.name, err)
				t.Fatal(err)
			}

			if res.StatusCode != e.expectedststcode {
				t.Errorf("For %s, expected %d but got %d", e.name, e.expectedststcode, res.StatusCode)
			} 
		}
	}
}