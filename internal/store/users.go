package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// EnsureSingleUser returns the id of the single-user row, creating it if
// missing. Used when SINGLE_USER=true. Email is optional.
func EnsureSingleUser(ctx context.Context, db DBTX, email, leetcodeUsername string) (int64, error) {
	const lookup = `SELECT id FROM users ORDER BY id LIMIT 1`
	var id int64
	err := db.QueryRow(ctx, lookup).Scan(&id)
	if err == nil {
		// Refresh leetcode username if provided.
		if leetcodeUsername != "" {
			if _, err := db.Exec(ctx,
				`UPDATE users SET leetcode_username = $2 WHERE id = $1 AND COALESCE(leetcode_username,'') <> $2`,
				id, leetcodeUsername,
			); err != nil {
				return 0, fmt.Errorf("update leetcode username: %w", err)
			}
		}
		return id, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("lookup user: %w", err)
	}

	const insert = `
INSERT INTO users (email, leetcode_username)
VALUES (NULLIF($1,''), NULLIF($2,''))
RETURNING id`
	if err := db.QueryRow(ctx, insert, email, leetcodeUsername).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	return id, nil
}

type SyncUser struct {
	ID               int64
	LeetcodeUsername string
}

func ListUsersForRecentSync(ctx context.Context, db DBTX) ([]SyncUser, error) {
	const q = `
SELECT id, leetcode_username
FROM users
WHERE COALESCE(leetcode_username, '') <> ''
ORDER BY id`
	rows, err := db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list users for recent sync: %w", err)
	}
	defer rows.Close()
	var out []SyncUser
	for rows.Next() {
		var u SyncUser
		if err := rows.Scan(&u.ID, &u.LeetcodeUsername); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func GetVacationUntil(ctx context.Context, db DBTX, userID int64) (*time.Time, error) {
	var until *time.Time
	err := db.QueryRow(ctx, `SELECT vacation_until FROM users WHERE id = $1`, userID).Scan(&until)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return until, err
}

func SetVacationUntil(ctx context.Context, db DBTX, userID int64, until *time.Time) error {
	_, err := db.Exec(ctx, `UPDATE users SET vacation_until = $2 WHERE id = $1`, userID, until)
	return err
}

// EnsureGoogleUser returns a verified user linked to the Google subject.
func EnsureGoogleUser(ctx context.Context, db DBTX, googleSub, email string) (int64, error) {
	googleSub = strings.TrimSpace(googleSub)
	email = strings.ToLower(strings.TrimSpace(email))
	if googleSub == "" || email == "" {
		return 0, errors.New("store: google subject and email are required")
	}

	const q = `
WITH updated_google AS (
	UPDATE users
	SET email_verified_at = COALESCE(email_verified_at, now())
	WHERE google_sub = $1
	RETURNING id
), inserted_or_linked AS (
	INSERT INTO users (email, google_sub, email_verified_at)
	SELECT $2, $1, now()
	WHERE NOT EXISTS (SELECT 1 FROM updated_google)
	ON CONFLICT (email) DO UPDATE
	SET google_sub = EXCLUDED.google_sub,
		email_verified_at = COALESCE(users.email_verified_at, now())
	WHERE users.google_sub IS NULL OR users.google_sub = EXCLUDED.google_sub
	RETURNING id
)
SELECT id FROM updated_google
UNION ALL
SELECT id FROM inserted_or_linked
LIMIT 1`
	var id int64
	err := db.QueryRow(ctx, q, googleSub, email).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.New("store: email already linked to another google account")
	}
	if err != nil {
		return 0, fmt.Errorf("ensure google user: %w", err)
	}
	return id, nil
}
