package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestPingController(t *testing.T) {

	tests := []struct {
		name   string
		method string
		status int
	}{
		{"fail when receive POST req", http.MethodPost, http.StatusMethodNotAllowed},
		{"fail when receive CONNECT req", http.MethodConnect, http.StatusMethodNotAllowed},
		{"fail when receive DELETE req", http.MethodDelete, http.StatusMethodNotAllowed},
		{"fail when receive HEAD req", http.MethodHead, http.StatusMethodNotAllowed},
		{"fail when receive OPTIONS req", http.MethodOptions, http.StatusMethodNotAllowed},
		{"fail when receive PATCH req", http.MethodPatch, http.StatusMethodNotAllowed},
		{"fail when receive PUT req", http.MethodPut, http.StatusMethodNotAllowed},
		{"fail when receive TRACE req", http.MethodTrace, http.StatusMethodNotAllowed},
		{"succeed when receive GET req", http.MethodGet, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewPingController()

			router := mux.NewRouter()
			router.HandleFunc("/", ctrl.HandlePingRequest)

			req, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Errorf("threw error when calling endpoint, got %s", err.Error())
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
			}
		})
	}
}
