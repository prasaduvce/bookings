package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prasaduvce/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	testRoutes := routes(&app)

	switch v := testRoutes.(type) {
	case *chi.Mux:
		// Test passed
	default:
		t.Errorf("Expected *chi.Mux, got %T", v)
	}
}