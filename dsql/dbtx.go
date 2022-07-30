package dsql

import (
	"context"
	"database/sql"
)

// DBTX represents an abstract database context.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type dbtxContextKey struct{}

func WithDBTXContext(parent context.Context, dbtx DBTX) context.Context {
	return context.WithValue(parent, dbtxContextKey{}, dbtx)
}

func DBTXFromContext(ctx context.Context) DBTX {
	dbtx, ok := ctx.Value(dbtxContextKey{}).(DBTX)
	if !ok {
		panic("dsql: no DBTX in the given context")
	}
	return dbtx
}
