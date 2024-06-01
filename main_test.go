package main

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSmsRoutes(t *testing.T) {
	tests := []struct {
		name               string
		route              string
		method             string
		expectedStatusCode int
	}{
		{
			name:               "SmsGet multiple",
			route:              "/sms/",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "SmsGet single",
			route:              "/sms/abc123",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "SmsGet single invalid sms id",
			route:              "/sms/:;-ZZZ",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "SmsSend",
			route:              "/sms",
			method:             http.MethodPost,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "SmsDelete",
			route:              "/sms/abc123",
			method:             http.MethodDelete,
			expectedStatusCode: http.StatusOK,
		},
	}

	handlersMock := NewSmsHandlersInterfaceMock(t)
	r := chi.NewRouter()
	registerSmsRoutes(r, handlersMock)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			if tt.expectedStatusCode != http.StatusNotFound {
				if tt.method == http.MethodGet {
					handlersMock.EXPECT().SmsGet(rr, mock.Anything)
				}

				if tt.method == http.MethodPost {
					handlersMock.EXPECT().SmsSend(rr, mock.Anything)
				}

				if tt.method == http.MethodDelete {
					handlersMock.EXPECT().SmsDelete(rr, mock.Anything)
				}
			}

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
