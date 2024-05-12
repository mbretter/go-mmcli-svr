package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

func HttpLoggerMiddleware(log *slog.Logger, requestIDKey any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			log2 := log // copy original logger

			ctx := r.Context()
			uid, ok := ctx.Value(requestIDKey).(string)

			if ok == true && len(uid) > 0 {
				log2 = log.With("uid", uid)
			}

			ctx = context.WithValue(ctx, "logger", log2)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
