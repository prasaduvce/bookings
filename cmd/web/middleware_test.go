package main

import (
	"log"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var myH myHandler

	testHandler := NoSurf(&myH)

	switch v := testHandler.(type) {
	case http.Handler:
		// Test passed
		log.Println("TestNoSurf: Passed")
	default:
		t.Errorf("Expected http.Handler, got %T", v) 
	}
}

func TestWriteToConsole(t *testing.T) {

	var myH myHandler

	testHandler := WriteToConsole(&myH)

	switch v := testHandler.(type) {
	case http.Handler:
		// Test passed
		log.Println("TestWriteToConsole: Passed")
	default:
		t.Errorf("Expected http.Handler, got %T", v) 
	}
}

func TestSessionLoad(t *testing.T) {

	var myH myHandler

	testHandler := SessionLoad(&myH)

	switch v := testHandler.(type) {
	case http.Handler:
		// Test passed
		log.Println("TestSessionLoad: Passed")
	default:
		t.Errorf("Expected http.Handler, got %T", v) 
	}
}