package store

import (
	"context"
	"errors"
	"fmt"

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
