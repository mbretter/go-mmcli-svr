package middleware

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

func LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := ctx.Value("logger").(*slog.Logger)

		routeCtx := ctx.Value(chi.RouteCtxKey).(*chi.Context)

		log.Info("Request", "method", routeCtx.RouteMethod, "headers", r.Header, "query", r.URL.Query(), "params", routeCtx.URLParams)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
