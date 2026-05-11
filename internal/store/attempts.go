package store

import (
	"context"
	"encoding/json"
	"fmt"

	"leetdrill/internal/models"
)

// InsertAttempt writes a row to attempts. Caller may set LeetcodeSubmissionID
// to enable extension/sync dedup via the partial unique index.
func InsertAttempt(ctx context.Context, db DBTX, a models.Attempt) (int64, error) {
	tags, err := json.Marshal(a.MistakeTags)
	if err != nil {
		return 0, fmt.Errorf("marshal mistake tags: %w", err)
	}
	const q = `
INSERT INTO attempts (
    user_id, problem_id, started_at, completed_at,
    verdict, submission_count_in_session, time_taken_sec,
    runtime_ms, memory_kb, language, code,
    derived_rating, journal, mistake_tags, leetcode_submission_id
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13, $14::jsonb, NULLIF($15,'')
)
ON CONFLICT (user_id, leetcode_submission_id)
    WHERE leetcode_submission_id IS NOT NULL
    DO NOTHING
RETURNING id`
	var id int64
	row := db.QueryRow(ctx, q,
		a.UserID, a.ProblemID, a.StartedAt, a.CompletedAt,
		a.Verdict, a.SubmissionCountInSession, a.TimeTakenSec,
		a.RuntimeMs, a.MemoryKB, a.Language, a.Code,
		a.DerivedRating, a.Journal, tags, a.LeetcodeSubmissionID,
	)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("insert attempt: %w", err)
	}
	return id, nil
}
