package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestModemList(t *testing.T) {
	newRequest := func(logger *slog.Logger) *http.Request {
		r := httptest.NewRequest("GET", "http://127.0.0.1:8743/modem/", nil)
		return r.WithContext(context.WithValue(r.Context(), "logger", logger))
	}

	tests := []struct {
		name             string
		expectStatusCode int
		expectBody       string
	}{
		{
			"Success",
			http.StatusOK,
			`{"modem-list":["/org/freedesktop/ModemManager1/Modem/0"]}`,
		},
		{
			"Error",
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, backendMock := ProvideTestApi(t)

			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			r := newRequest(logger)
			w := httptest.NewRecorder()

			if tt.expectStatusCode == http.StatusInternalServerError {
				backendMock.EXPECT().Exec("-L").Return(nil, errors.New("failed"))
			} else {
				backendMock.EXPECT().Exec("-L").Return([]byte(tt.expectBody), nil)
			}

			api.ModemList(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectBody, strings.Trim(string(body), " \n"))
		})
	}
}

func TestModemDetail(t *testing.T) {
	newRequest := func(modemId string, logger *slog.Logger) *http.Request {
		r := httptest.NewRequest("GET", "http://127.0.0.1:8743/modem/", nil)

		ctx := context.WithValue(r.Context(), "logger", logger)

		if len(modemId) > 0 {
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", modemId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
		}

		return r.WithContext(ctx)
	}

	tests := []struct {
		name             string
		modemId          string
		expectStatusCode int
		expectBody       string
	}{
		{
			"Success",
			"1",
			http.StatusOK,
			`{"modem":{"3gpp":{"5gnr"...}}}`,
		},
		{
			"Error",
			"1",
			http.StatusInternalServerError,
			`{"error":"internal server error"}`,
		},
		{
			"InvalidModemId",
			"/&.",
			http.StatusBadRequest,
			`{"error":"invalid ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, backendMock := ProvideTestApi(t)

			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			r := newRequest(tt.modemId, logger)
			w := httptest.NewRecorder()

			if tt.expectStatusCode == http.StatusInternalServerError {
				backendMock.EXPECT().Exec("-m", tt.modemId).Return(nil, errors.New("failed"))
			} else if tt.expectStatusCode == http.StatusBadRequest {
			} else {
				backendMock.EXPECT().Exec("-m", tt.modemId).Return([]byte(tt.expectBody), nil)
			}

			api.ModemDetail(w, r)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectBody, strings.Trim(string(body), " \n"))
		})
	}
}
