package store

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestDuplicateAttemptIsDetectedFromNoReturningRow(t *testing.T) {
	if !errors.Is(classifyInsertAttemptScanErr(pgx.ErrNoRows), ErrDuplicateAttempt) {
		t.Fatalf("expected pgx.ErrNoRows to classify as ErrDuplicateAttempt")
	}
}
