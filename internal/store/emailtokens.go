package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// ErrTokenInvalid is returned when a token is not found, expired, or already used.
var ErrTokenInvalid = errors.New("store: token invalid or expired")

// EmailTokenKind enumerates the two token types.
type EmailTokenKind string

const (
	EmailTokenVerify EmailTokenKind = "verify"
	EmailTokenReset  EmailTokenKind = "reset"
)

// CreateEmailToken invalidates any existing unused tokens of the same
// user+kind, then inserts a new hashed token.
func CreateEmailToken(ctx context.Context, db DBTX, userID int64, kind EmailTokenKind, tokenHash []byte, expiresAt time.Time) error {
	_, err := db.Exec(ctx,
		`UPDATE email_tokens SET used_at = now()
		 WHERE user_id = $1 AND kind = $2 AND used_at IS NULL`,
		userID, string(kind),
	)
	if err != nil {
		return fmt.Errorf("invalidate old email tokens: %w", err)
	}
	_, err = db.Exec(ctx,
		`INSERT INTO email_tokens (user_id, kind, token_hash, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		userID, string(kind), tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert email token: %w", err)
	}
	return nil
}

// ConsumeEmailToken finds a valid (unused, unexpired) token by hash and kind,
// marks it used, and returns the associated user_id.
// Returns ErrTokenInvalid if not found, expired, or already used.
func ConsumeEmailToken(ctx context.Context, db DBTX, tokenHash []byte, kind EmailTokenKind) (int64, error) {
	const q = `
UPDATE email_tokens
SET used_at = now()
WHERE token_hash = $1
  AND kind = $2
  AND used_at IS NULL
  AND expires_at > now()
RETURNING user_id`

	var userID int64
	err := db.QueryRow(ctx, q, tokenHash, string(kind)).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrTokenInvalid
	}
	if err != nil {
		return 0, fmt.Errorf("consume email token: %w", err)
	}
	return userID, nil
}

// MarkEmailVerified sets email_verified_at = now() for the given user.
func MarkEmailVerified(ctx context.Context, db DBTX, userID int64) error {
	_, err := db.Exec(ctx,
		`UPDATE users SET email_verified_at = now() WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("mark email verified: %w", err)
	}
	return nil
}

// SetPasswordHash updates the password_hash for a user.
func SetPasswordHash(ctx context.Context, db DBTX, userID int64, hash string) error {
	_, err := db.Exec(ctx,
		`UPDATE users SET password_hash = $2 WHERE id = $1`,
		userID, hash,
	)
	if err != nil {
		return fmt.Errorf("set password hash: %w", err)
	}
	return nil
}

// GetUserIDByEmail returns the user's id and email_verified_at by email.
// Returns ErrNotFound if no user exists with that email.
func GetUserIDByEmail(ctx context.Context, db DBTX, email string) (id int64, verifiedAt *time.Time, err error) {
	const q = `SELECT id, email_verified_at FROM users WHERE email = $1`
	err = db.QueryRow(ctx, q, email).Scan(&id, &verifiedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil, ErrNotFound
	}
	if err != nil {
		return 0, nil, fmt.Errorf("get user by email: %w", err)
	}
	return id, verifiedAt, nil
}

// GetEmailVerifiedAt returns the email_verified_at for a user id.
func GetEmailVerifiedAt(ctx context.Context, db DBTX, userID int64) (*time.Time, error) {
	var t *time.Time
	err := db.QueryRow(ctx, `SELECT email_verified_at FROM users WHERE id = $1`, userID).Scan(&t)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get email_verified_at: %w", err)
	}
	return t, nil
}
