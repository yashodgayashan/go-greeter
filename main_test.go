package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGreetHandler tests the greet function
func TestGreetHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/greet?name=Alice", nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	expected := "Hello, Alice!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestGreetHandlerWithoutName tests the greet function when no name is provided
func TestGreetHandlerWithoutName(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/greet", nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	expected := "Hello, Stranger!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}
