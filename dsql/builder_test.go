package dsql_test

import (
	"testing"

	"github.com/semrekkers/go-service/dsql"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	var (
		b     dsql.Builder
		table = "users"
	)

	query, args := b.Write("SELECT * ").
		Writef("FROM %s ", table).
		Writep("WHERE id = {}", 432).
		Done()

	assert.Equal(t, `SELECT * FROM users WHERE id = $1`, query)
	assert.Equal(t, []any{432}, args)
}

func TestBuilderEnumParam(t *testing.T) {
	var (
		b = dsql.Builder{
			ParamFmt: dsql.EnumParameter,
		}
		table = "users"
	)

	query, args := b.Write("SELECT * ").
		Writef("FROM %s ", table).
		Writep("WHERE id = {}", 432).
		Done()

	assert.Equal(t, `SELECT * FROM users WHERE id = ?`, query)
	assert.Equal(t, []any{432}, args)
}

func TestBuilderComplex(t *testing.T) {
	var (
		b          dsql.Builder
		table      = "users"
		predicates = []dsql.NamedValue{
			{"role LIKE {}", "admin%"},
			{"primary_group_id = {}", 1},
			{"is_active", nil /* no need */},
			{"country = {}", "NL"},
			{"buildings.place IN ({})", []string{"Rotterdam", "Amsterdam"}},
			{"(manager_id = {} OR manager_id = {})", dsql.Values{765, 92}},
		}
	)

	b.Write("SELECT * ").
		Writef("FROM %s ", table).
		Write("LEFT JOIN buildings ON users.last_presence_id = buildings.id ")
	if len(predicates) > 0 {
		b.Write("WHERE ").
			Writev("%s", " AND ", predicates...)
	}
	query, args := b.Done()

	assert.Equal(t, `SELECT * FROM users LEFT JOIN buildings ON users.last_presence_id = buildings.id WHERE role LIKE $1 AND primary_group_id = $2 AND is_active AND country = $3 AND buildings.place IN ($4) AND (manager_id = $5 OR manager_id = $6)`, query)
	assert.Equal(t, []any{"admin%", 1, "NL", []string{"Rotterdam", "Amsterdam"}, 765, 92}, args)
}

func TestBuilderValuesLeftovers(t *testing.T) {
	var (
		b          dsql.Builder
		table      = "users"
		predicates = []dsql.NamedValue{
			{"manager_id = {}", dsql.Values{765, 92 /* <- leftover */}},
		}
	)

	b.Write("SELECT * ").
		Writef("FROM %s ", table)
	if len(predicates) > 0 {
		b.Write("WHERE ").
			Writev("%s", " AND ", predicates...)
	}
	query, args := b.Done()

	assert.Equal(t, `SELECT * FROM users WHERE manager_id = $1`, query)
	assert.Equal(t, []any{765}, args)
}
