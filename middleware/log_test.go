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

func TestLogRoute(t *testing.T) {

	var buff bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buff, nil))

	req := httptest.NewRequest("GET", "http://localhost/?modem=1", nil)
	req.Header.Set("User-Agent", "test")
	ctx := context.WithValue(req.Context(), "logger", logger)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("foo", "bar")
	chiCtx.RouteMethod = "GET"
	ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)

	req = req.WithContext(ctx)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	rr := httptest.NewRecorder()

	// Run the middleware function
	LogRoute(nextHandler).ServeHTTP(rr, req)

	assert.Contains(t, buff.String(), `level=INFO msg=Request method=GET headers=map[User-Agent:[test]] query=map[modem:[1]] params="{Keys:[foo] Values:[bar]}"`)
}
