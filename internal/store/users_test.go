package store

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestEnsureGoogleUserFirstLinksByGoogleSub(t *testing.T) {
	db := &captureGoogleQueryRowDB{row: fakeRow{id: 12}}

	id, err := EnsureGoogleUser(context.Background(), db, "google-sub", "USER@example.COM")
	if err != nil {
		t.Fatalf("EnsureGoogleUser() error = %v", err)
	}
	if id != 12 {
		t.Fatalf("id = %d, want 12", id)
	}
	if !strings.Contains(db.sql, "google_sub = $1") {
		t.Fatalf("query must prefer existing google_sub:\n%s", db.sql)
	}
	if !strings.Contains(db.sql, "ON CONFLICT (email)") {
		t.Fatalf("query must link existing email rows:\n%s", db.sql)
	}
	if got, want := db.args[1], "user@example.com"; got != want {
		t.Fatalf("email arg = %q, want %q", got, want)
	}
}

type captureGoogleQueryRowDB struct {
	sql  string
	args []any
	row  pgx.Row
}

func (db *captureGoogleQueryRowDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	panic("Exec not implemented")
}

func (db *captureGoogleQueryRowDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("Query not implemented")
}

func (db *captureGoogleQueryRowDB) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	db.sql = sql
	db.args = args
	return db.row
}

type fakeRow struct {
	id  int64
	err error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	*(dest[0].(*int64)) = r.id
	return nil
}
