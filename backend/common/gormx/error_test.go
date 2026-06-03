package gormx

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

func TestIsUniqueViolation(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"plain", errors.New("boom"), false},
		{"direct 23505", &pgconn.PgError{Code: "23505"}, true},
		{"wrapped 23505", errors.WithStack(&pgconn.PgError{Code: "23505"}), true},
		{"double wrapped 23505", errors.Wrap(errors.WithStack(&pgconn.PgError{Code: "23505"}), "create"), true},
		{"foreign key 23503", &pgconn.PgError{Code: "23503"}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsUniqueViolation(c.err); got != c.want {
				t.Fatalf("IsUniqueViolation(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}
