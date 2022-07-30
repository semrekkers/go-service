package dsql

import (
	"fmt"
	"strconv"
	"strings"
)

// A ParameterFormatter formats a parameter to a specific SQL dialect.
type ParameterFormatter func(sb *strings.Builder, i int)

// A Builder builds a query string using Write methods. An empty Builder is
// ready to use. Do not copy a non-empty Builder, use Clone() instead.
type Builder struct {
	ParamFmt ParameterFormatter

	s    strings.Builder
	args []any
}

// Returns a new Builder.
func NewBuilder() *Builder {
	return &Builder{
		ParamFmt: PositionedParameter,
	}
}

// Write appends s to the Builder's buffer. Returns the receiver Builder.
func (b *Builder) Write(s string) *Builder {
	b.s.WriteString(s)
	return b
}

// Writef formats according to a fmt format specifier and appends to the Builder's buffer.
// Returns the receiver Builder.
func (b *Builder) Writef(format string, a ...any) *Builder {
	fmt.Fprintf(&b.s, format, a...)
	return b
}

const paramTemplate = "{}"

// Writep substitutes every parameter denoted by "{}" in s to a formatted parameter, and appends it along
// with the positioned argument from a, to the Builder's buffer. The parameter format can be configured using
// Builder.ParameterFormatter. Returns the receiver Builder.
func (b *Builder) Writep(s string, p ...any) *Builder {
	if b.ParamFmt == nil {
		b.ParamFmt = PositionedParameter // default
	}
	b.s.Grow(len(s))
	inline := 0
	for {
		i := strings.Index(s, paramTemplate)
		if i == -1 {
			break
		}
		if inline > 0 {
			// This parameter value was already appended to b.args (inlined)
			b.s.WriteString(s[:i])
			b.ParamFmt(&b.s, len(b.args)-inline)
			s = s[i+len(paramTemplate):] // only advance s
			inline--
			continue
		}
		if len(p) < 1 {
			panic("dsql: insufficient parameters")
		}
		b.s.WriteString(s[:i])
		b.ParamFmt(&b.s, len(b.args))
		if values, ok := p[0].(Values); ok {
			b.args = append(b.args, values...)
			// Next parameters represent the now inlined values
			inline += len(values) - 1
		} else {
			b.args = append(b.args, p[0])
		}
		p, s = p[1:], s[i+len(paramTemplate):] // advance
	}
	b.s.WriteString(s)
	b.args = b.args[:len(b.args)-inline] // slice off any unused inline arguments
	return b
}

// NamedValue represents a named value.
type NamedValue struct {
	Name  string
	Value any
}

// Values represents multiple values. This can be helpful when you want to pass multiple
// values inside a NamedValue for example.
type Values []any

// Writev substitutes every NamedValue denoted by "{}" to a formatted parameter, and
// appends it along with the positioned Value from a to the Builder's buffer. The parameter format can be
// configured using Builder.ParameterFormatter. Returns the receiver Builder.
func (b *Builder) Writev(format string, sep string, a ...NamedValue) *Builder {
	if len(a) > 0 {
		b.Writep(strings.Replace(format, "%s", a[0].Name, 1), a[0].Value)
		for _, arg := range a[1:] {
			b.s.WriteString(sep)
			b.Writep(strings.Replace(format, "%s", arg.Name, 1), arg.Value)
		}
	}
	return b
}

// Clone returns a clone of the receiver Builder.
func (b *Builder) Clone() *Builder {
	c := &Builder{
		ParamFmt: b.ParamFmt,
		args:     append([]any(nil), b.args...),
	}
	c.s.WriteString(b.s.String())
	return c
}

// String returns the query string.
func (b *Builder) String() string {
	return b.s.String()
}

// Done finishes the builder and returns the query string with it's arguments.
func (b *Builder) Done() (string, []any) {
	return b.s.String(), b.args
}

// PositionedParameter formats the parameter as "$n" (e.g. PostgreSQL style).
// This is the default format.
func PositionedParameter(sb *strings.Builder, i int) {
	sb.WriteByte('$')
	sb.WriteString(strconv.Itoa(i + 1)) // positions start at 1, not 0
}

// EnumParameter or the enumerated parameter formatter, formats the parameter as "?" (e.g. MySQL/SQLite style).
func EnumParameter(sb *strings.Builder, _ int) {
	sb.WriteByte('?')
}
