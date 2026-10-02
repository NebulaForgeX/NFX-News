// Package enumx provides a generic slice type for scanning and writing
// PostgreSQL enum arrays (e.g. "social".group_invite_policy[]) through
// database/sql / GORM without depending on a driver-specific array type.
package enumx

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Array is a string-backed enum slice that maps to a PostgreSQL array column.
//
// It implements sql.Scanner and driver.Valuer using the Postgres array text
// representation ("{a,b,c}"), which every supported driver exposes via
// database/sql, so it works for enum arrays without a custom driver type.
//
// A nil/empty Array is written as the empty array literal "{}".
type Array[T ~string] []T

// Value renders the slice as a PostgreSQL array literal ("{a,b,c}").
//
// Type Parameters:
//   - T: string-backed enum type.
//
// Returns:
//   - driver.Value: the array literal string.
//   - error: always nil.
func (a Array[T]) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	elems := make([]string, len(a))
	for i, v := range a {
		elems[i] = quoteArrayElement(string(v))
	}
	return "{" + strings.Join(elems, ",") + "}", nil
}

// Scan parses a PostgreSQL array value ("{a,b,c}") into the slice.
//
// Supported sources are string, []byte and nil (nil yields an empty slice).
//
// Type Parameters:
//   - T: string-backed enum type.
//
// Parameters:
//   - src: raw value from the database driver.
//
// Returns:
//   - error: non-nil when src is an unsupported type or malformed literal.
func (a *Array[T]) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}

	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("enumx: cannot scan %T into Array", src)
	}

	parsed, err := parseArrayLiteral(raw)
	if err != nil {
		return err
	}
	out := make(Array[T], len(parsed))
	for i, s := range parsed {
		out[i] = T(s)
	}
	*a = out
	return nil
}

// quoteArrayElement quotes an element only when required by the Postgres array
// syntax (empty, contains a comma, brace, quote, backslash or whitespace).
func quoteArrayElement(s string) string {
	if s == "" {
		return `""`
	}
	if strings.ContainsAny(s, `,{}"\ `) {
		escaped := strings.ReplaceAll(s, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return s
}

// parseArrayLiteral parses a single-dimension Postgres array literal into its
// element strings. It understands quoting and backslash escapes.
func parseArrayLiteral(raw string) ([]string, error) {
	s := strings.TrimSpace(raw)
	if s == "" || s == "{}" {
		return nil, nil
	}
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return nil, fmt.Errorf("enumx: malformed array literal %q", raw)
	}
	body := s[1 : len(s)-1]

	var (
		out     []string
		cur     strings.Builder
		quoted  bool
		escaped bool
		hasElem bool
	)
	for i := 0; i < len(body); i++ {
		c := body[i]
		switch {
		case escaped:
			cur.WriteByte(c)
			escaped = false
			hasElem = true
		case c == '\\':
			escaped = true
			hasElem = true
		case c == '"':
			quoted = !quoted
			hasElem = true
		case c == ',' && !quoted:
			out = append(out, cur.String())
			cur.Reset()
			hasElem = false
		default:
			cur.WriteByte(c)
			hasElem = true
		}
	}
	if quoted || escaped {
		return nil, fmt.Errorf("enumx: malformed array literal %q", raw)
	}
	if hasElem || len(out) > 0 {
		out = append(out, cur.String())
	}
	return out, nil
}
