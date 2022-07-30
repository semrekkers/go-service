package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
)

type Trace struct {
	ID   string
	Meta any
}

type traceCtxKey struct{}

func WithTraceContext(parent context.Context, trace *Trace) context.Context {
	return context.WithValue(parent, traceCtxKey{}, trace)
}

func TraceFromContext(ctx context.Context) *Trace {
	trace, ok := ctx.Value(traceCtxKey{}).(*Trace)
	if !ok {
		panic("tracing: no Trace in the given context")
	}
	return trace
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(
			WithTraceContext(r.Context(), &Trace{ID: newTraceID()}),
		))
	})
}

func newTraceID() string {
	var buf [16]byte
	if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf[:])
}
