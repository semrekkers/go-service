package apih

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type loggerCtxKey struct{}

func WithLoggerContext(parent context.Context, log *zap.Logger) context.Context {
	return context.WithValue(parent, loggerCtxKey{}, log)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger)
	if !ok {
		panic("apih: no Logger in the given context")
	}
	return log
}

func LoggerMiddleware(log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				WithLoggerContext(r.Context(), log),
			))
		})
	}
}
