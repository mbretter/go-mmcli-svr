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

func TestLogger(t *testing.T) {
	const requestIdField = "requestId"
	const requestId = "657483"

	var buff bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buff, nil))

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	ctx := context.WithValue(req.Context(), requestIdField, requestId)

	chiCtx := chi.NewRouteContext()
	ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mw := HttpLoggerMiddleware(logger, requestIdField)

	// Run the middleware function
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log2, ok := r.Context().Value("logger").(*slog.Logger)
		log2.Info("test")

		assert.Contains(t, buff.String(), "level=INFO msg=test uid="+requestId)
		assert.True(t, ok)
		assert.NotEqual(t, logger, log2)
	})).ServeHTTP(rr, req)
}
