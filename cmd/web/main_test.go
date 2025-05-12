package main

import "testing"

func TestRun(t *testing.T) {
	// Test the run function
	_, err := run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

}