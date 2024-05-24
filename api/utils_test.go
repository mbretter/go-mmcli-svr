package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUtils_Ping(t *testing.T) {
	u := ProvideUtilsApi()
	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "NormalPing",
			expectedCode: http.StatusOK,
			expectedBody: "boing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/ping", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(u.Ping)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")
			assert.Equal(t, tt.expectedBody, rr.Body.String(), "handler returned unexpected body")
		})
	}
}
