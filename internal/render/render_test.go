package render

import (
	"net/http"
	"testing"

	"github.com/prasaduvce/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()

	if (err != nil) {
		t.Fatalf("Failed to get session: %v", err)
	}
	result := AddDefaultData(&td, r)

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if result.Flash != "123" {
		t.Errorf("Expected flash message '123', got '%s'", result.Flash)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {	
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(r.Context(), r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	session.Put(r.Context(), "flash", "123")
	return r, nil
}

func TestRenderHtml(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Fatalf("Failed to create template cache: %v", err)
	}
	appConfig.Templates = tc

	r, err := getSession()
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	var ww myWriter

	err = RenderHtml(&ww, "home.page.tmpl", &models.TemplateData{}, r)

	if err != nil {
		t.Fatalf("Failed to render HTML: %v", err)
	}

	err = RenderHtml(&ww, "non-existing.page.tmpl", &models.TemplateData{}, r)
	if err == nil {
		t.Fatalf("Expected error for non-existing template, got nil")
	}
}