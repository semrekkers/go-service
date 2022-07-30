package dsql

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Null represents a nullable type of `T`.
type Null[T any] struct {
	Some  T
	Valid bool
}

// NewNull returns a new nullable type of `T`.
func NewNull[T any](value T) Null[T] {
	return Null[T]{
		Some:  value,
		Valid: true,
	}
}

// Set sets the value.
func (n *Null[T]) Set(value T) {
	n.Some, n.Valid = value, true
}

// Scan implements `sql.Scanner`.
func (n *Null[T]) Scan(value any) error {
	if value == nil {
		var zero T
		n.Some, n.Valid = zero, false
		return nil
	}
	n.Some, n.Valid = value.(T)
	if !n.Valid {
		return fmt.Errorf("dsql.Null: converting value type %T to %T is unsupported", value, n.Some)
	}
	return nil
}

// Value implements `sql/driver.Valuer`.
func (n Null[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Some, nil
}

// Ptr returns a pointer to the value if valid, otherwise nil.
func (n Null[T]) Ptr() *T {
	if !n.Valid {
		return nil
	}
	return &n.Some
}

var nullBytes = []byte("null")

// MarshalJSON implements `json.Marshaler`.
func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return nullBytes, nil
	}
	return json.Marshal(n.Some)
}

// UnmarshalJSON implements `json.Unmarshaler`
func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		var zero T
		n.Some, n.Valid = zero, false
		return nil
	}
	if err := json.Unmarshal(data, &n.Some); err != nil {
		return fmt.Errorf("dsql.Null: could not unmarshal type %T: %w", n.Some, err)
	}
	n.Valid = true
	return nil
}
