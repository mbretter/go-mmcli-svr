package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name             string
		modem            string
		expectStatusCode int
		expectBody       string
	}{
		{
			"Success",
			"",
			http.StatusOK,
			`{"modem":{"location":{"3gpp":{"cid":"000AA000"}}}}`,
		},
		{
			"Error",
			"0",
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"ModemNotFound",
			"0",
			http.StatusNotFound,
			`{"error":"error: couldn't find modem"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, backendMock := ProvideTestApi(t)

			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			r := newReadRequest(tt.modem, logger)
			w := httptest.NewRecorder()

			if tt.expectStatusCode == http.StatusInternalServerError {
				backendMock.EXPECT().ExecModem(tt.modem, "--location-get").Return(nil, errors.New("failed"))
			} else if tt.expectStatusCode == http.StatusNotFound {
				backendMock.EXPECT().ExecModem(tt.modem, "--location-get").Return(nil, errors.New("error: couldn't find modem"))
			} else {
				backendMock.EXPECT().ExecModem(tt.modem, "--location-get").Return([]byte(tt.expectBody), nil)
			}

			api.LocationGet(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectBody, strings.Trim(string(body), " \n"))
		})
	}
}

func newReadRequest(modemId string, logger *slog.Logger) *http.Request {
	r := httptest.NewRequest("GET", "http://127.0.0.1:8743/location", nil)

	ctx := context.WithValue(r.Context(), "modem", modemId)
	ctx = context.WithValue(ctx, "logger", logger)

	return r.WithContext(ctx)
}
