package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"leetdrill/internal/models"
)

const leetcodeURLPrefix = "https://leetcode.com/problems/"

// UpsertProblem inserts or updates a row keyed by leetcode_slug. Returns the
// problem id, which is stable across upserts.
func UpsertProblem(ctx context.Context, db DBTX, p models.Problem) (int64, error) {
	tags, err := json.Marshal(p.TopicTags)
	if err != nil {
		return 0, fmt.Errorf("marshal topic tags: %w", err)
	}
	if p.URL == "" {
		p.URL = leetcodeURLPrefix + p.LeetcodeSlug + "/"
	}
	const q = `
INSERT INTO problems (
    leetcode_slug, leetcode_question_id, leetcode_frontend_id,
    title, difficulty, url, content_html, topic_tags, ac_rate, paid_only, synced_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9, $10, now())
ON CONFLICT (leetcode_slug) DO UPDATE SET
    leetcode_question_id = EXCLUDED.leetcode_question_id,
    leetcode_frontend_id = EXCLUDED.leetcode_frontend_id,
    title                = EXCLUDED.title,
    difficulty           = EXCLUDED.difficulty,
    url                  = EXCLUDED.url,
    content_html         = COALESCE(NULLIF(EXCLUDED.content_html, ''), problems.content_html),
    topic_tags           = EXCLUDED.topic_tags,
    ac_rate              = EXCLUDED.ac_rate,
    paid_only            = EXCLUDED.paid_only,
    synced_at            = now()
RETURNING id`
	var id int64
	row := db.QueryRow(ctx, q,
		p.LeetcodeSlug, p.LeetcodeQuestionID, p.LeetcodeFrontendID,
		p.Title, string(p.Difficulty), p.URL, p.ContentHTML, tags, p.ACRate, p.PaidOnly,
	)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("upsert problem %s: %w", p.LeetcodeSlug, err)
	}
	return id, nil
}

// UpsertPattern inserts or updates a pattern row keyed by slug. Returns id.
func UpsertPattern(ctx context.Context, db DBTX, slug, name string) (int64, error) {
	const q = `
INSERT INTO patterns (slug, name)
VALUES ($1, $2)
ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
RETURNING id`
	var id int64
	row := db.QueryRow(ctx, q, slug, name)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("upsert pattern %s: %w", slug, err)
	}
	return id, nil
}

// LinkProblemPattern is idempotent.
func LinkProblemPattern(ctx context.Context, db DBTX, problemID, patternID int64) error {
	const q = `
INSERT INTO problem_patterns (problem_id, pattern_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING`
	_, err := db.Exec(ctx, q, problemID, patternID)
	if err != nil {
		return fmt.Errorf("link problem %d ↔ pattern %d: %w", problemID, patternID, err)
	}
	return nil
}

// GetProblemBySlug returns a problem row.
func GetProblemBySlug(ctx context.Context, db DBTX, slug string) (*models.Problem, error) {
	const q = `
SELECT id, leetcode_slug, COALESCE(leetcode_question_id,''), COALESCE(leetcode_frontend_id,''),
       title, difficulty, url, COALESCE(content_html,''), topic_tags,
       COALESCE(ac_rate, 0), paid_only, synced_at
FROM problems WHERE leetcode_slug = $1`
	var p models.Problem
	var tagsBytes []byte
	var diff string
	err := db.QueryRow(ctx, q, slug).Scan(
		&p.ID, &p.LeetcodeSlug, &p.LeetcodeQuestionID, &p.LeetcodeFrontendID,
		&p.Title, &diff, &p.URL, &p.ContentHTML, &tagsBytes,
		&p.ACRate, &p.PaidOnly, &p.SyncedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get problem %s: %w", slug, err)
	}
	p.Difficulty = models.Difficulty(diff)
	if len(tagsBytes) > 0 {
		_ = json.Unmarshal(tagsBytes, &p.TopicTags)
	}
	return &p, nil
}
