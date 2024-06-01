package middleware

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpModemMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		modem            string
		expectStatusCode int
	}{
		{
			"Success",
			"1",
			http.StatusOK,
		},
		{
			"Success Empty",
			"",
			http.StatusOK,
		},
		{
			"Error invalid modem id",
			".*'",
			http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buff bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buff, nil))

			req := httptest.NewRequest("GET", "http://localhost/?modem="+tt.modem, nil)
			ctx := context.WithValue(req.Context(), "logger", logger)

			chiCtx := chi.NewRouteContext()
			ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			mw := HttpModemMiddleware(logger)

			// Run the middleware function
			mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				modem, ok := r.Context().Value("modem").(string)

				if tt.expectStatusCode == http.StatusOK {
					if tt.modem != "" {
						assert.Contains(t, buff.String(), "level=INFO msg=HttpModemMiddleware modem="+tt.modem)
					}

					assert.Equal(t, modem, tt.modem)
					assert.True(t, ok)
				} else {
					assert.False(t, ok)
				}
			})).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectStatusCode, rr.Code)
		})
	}
}
