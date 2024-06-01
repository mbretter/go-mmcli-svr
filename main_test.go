package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mbretter/go-mmcli-svr/backend"
	"net/http"
	"testing"
)

func TestRegisterSmsRoutes(t *testing.T) {
	backendMock := backend.NewBackendMock(t)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Get All Sms", http.MethodGet, "/sms/", http.StatusOK},
		{"Get Sms By ID", http.MethodGet, "/sms/1", http.StatusOK},
		{"Send Sms", http.MethodPost, "/sms", http.StatusOK},
		{"Delete Sms", http.MethodDelete, "/sms/1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			registerSmsRoutes(r, backendMock)

			r.Routes()
		})
	}
}
