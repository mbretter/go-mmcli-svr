package middleware

import (
	"context"
	"github.com/mbretter/go-mmcli-svr/api"
	utils "github.com/mbretter/go-mmcli-svr/http"
	"log/slog"
	"net/http"
)

func HttpModemMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			modem := r.URL.Query().Get("modem")
			if modem != "" {
				log.Info("HttpModemMiddleware", "modem", modem)
			}

			modem, err := api.ValidatePathIndex(modem)
			if err != nil {
				utils.WriteError(w, r, http.StatusBadRequest, err)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "modem", modem)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
