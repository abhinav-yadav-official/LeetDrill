// Package sync imports and reconciles LeetCode history into LeetDrill state.
package sync

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"leetdrill/internal/leetcode"
	"leetdrill/internal/models"
	"leetdrill/internal/store"
	"leetdrill/internal/vault"
)

type leetcodeClient interface {
	GetUser(ctx context.Context, username string) (*leetcode.MatchedUser, error)
	RecentACSubmissions(ctx context.Context, username string, limit int) ([]leetcode.RecentSubmission, error)
	ListSolvedProblems(ctx context.Context, creds leetcode.Credentials, skip, limit int) ([]leetcode.ProblemListItem, int, error)
	ListSubmissions(ctx context.Context, creds leetcode.Credentials, offset, limit int, lastKey, questionSlug string) (*leetcode.SubmissionListPage, error)
}

type ColdStartImporter struct {
	Store              *store.Store
	Vault              *vault.Vault
	Client             leetcodeClient
	MaxSubmissionPages int
	SubmissionPageSize int
}

type ColdStartResult struct {
	Username          string `json:"username"`
	PublicSolved      int    `json:"public_solved"`
	RecentImported    int    `json:"recent_imported"`
	AuthedImported    int    `json:"authed_imported"`
	SolvedImported    int    `json:"solved_imported"`
	DuplicatesSkipped int    `json:"duplicates_skipped"`
	UnknownSkipped    int    `json:"unknown_skipped"`
	AuthedSkipped     bool   `json:"authed_skipped"`
}

func (i *ColdStartImporter) Run(ctx context.Context, userID int64, username string) (ColdStartResult, error) {
	if i.Store == nil {
		return ColdStartResult{}, errors.New("cold-start: store required")
	}
	if i.Client == nil {
		i.Client = leetcode.New()
	}
	username = strings.TrimSpace(username)
	if username == "" {
		c, err := store.GetLeetcodeCookies(ctx, i.Store.DB(), userID)
		if err != nil && !errors.Is(err, store.ErrNotFound) {
			return ColdStartResult{}, err
		}
		username = strings.TrimSpace(c.Username)
	}

	result := ColdStartResult{Username: username}
	now := time.Now().UTC()
	if username != "" {
		user, err := i.Client.GetUser(ctx, username)
		if err != nil {
			return result, err
		}
		result.Username = user.Username
		result.PublicSolved = solvedCount(user)
		if err := setLeetcodeUsername(ctx, i.Store.DB(), userID, user.Username); err != nil {
			return result, err
		}

		recent, err := i.Client.RecentACSubmissions(ctx, user.Username, 20)
		if err != nil {
			return result, err
		}
		for idx, sub := range reverseRecent(recent) {
			imported, err := i.importAccepted(ctx, userID, sub.TitleSlug, sub.ID, parseOrNow(sub.Timestamp, now), coldStartDueAt(now, idx))
			if err != nil {
				return result, err
			}
			switch imported {
			case importInserted:
				result.RecentImported++
			case importDuplicate:
				result.DuplicatesSkipped++
			case importUnknownProblem:
				result.UnknownSkipped++
			}
		}
	}

	authedResult, err := i.importAuthedAccepted(ctx, userID, now)
	if err != nil {
		return result, err
	}
	result.AuthedImported = authedResult.imported
	result.SolvedImported = authedResult.solved
	result.DuplicatesSkipped += authedResult.duplicates
	result.UnknownSkipped += authedResult.unknown
	result.AuthedSkipped = authedResult.skipped
	if result.AuthedSkipped && username == "" {
		return result, errors.New("cold-start: leetcode username or synced cookies required")
	}
	if err := coldStartEmptyImportError(result); err != nil {
		return result, err
	}
	return result, nil
}

func coldStartEmptyImportError(result ColdStartResult) error {
	if result.PublicSolved > 0 &&
		result.RecentImported == 0 &&
		result.AuthedImported == 0 &&
		result.SolvedImported == 0 &&
		result.DuplicatesSkipped == 0 &&
		result.AuthedSkipped {
		return fmt.Errorf(
			"cold-start: public profile has %d solved problems but no public recent AC slugs; sync cookies for full history import",
			result.PublicSolved,
		)
	}
	return nil
}

type importStatus int

const (
	importInserted importStatus = iota
	importDuplicate
	importUnknownProblem
)

func (i *ColdStartImporter) importAccepted(ctx context.Context, userID int64, slug, submissionID string, completedAt, dueAt time.Time) (importStatus, error) {
	p, err := store.GetProblemBySlug(ctx, i.Store.DB(), slug)
	if errors.Is(err, store.ErrNotFound) {
		return importUnknownProblem, nil
	}
	if err != nil {
		return importUnknownProblem, err
	}
	_, err = i.Store.Apply(ctx, store.ApplyInput{
		UserID:          userID,
		ProblemID:       p.ID,
		Difficulty:      p.Difficulty,
		CompletedAt:     completedAt,
		Verdict:         leetcode.VerdictCodeAC,
		SubmissionCount: 1,
		TimeTakenSec:    0,
		LeetcodeSubmID:  submissionID,
		Now:             time.Now().UTC(),
	})
	if errors.Is(err, store.ErrDuplicateAttempt) {
		return importDuplicate, nil
	}
	if err != nil {
		return importUnknownProblem, err
	}
	if err := setNextDueAt(ctx, i.Store.DB(), userID, p.ID, dueAt); err != nil {
		return importUnknownProblem, err
	}
	return importInserted, nil
}

type authedImportResult struct {
	imported   int
	solved     int
	duplicates int
	unknown    int
	skipped    bool
}

func (i *ColdStartImporter) importAuthedAccepted(ctx context.Context, userID int64, now time.Time) (authedImportResult, error) {
	if i.Vault == nil {
		return authedImportResult{skipped: true}, nil
	}
	c, err := store.GetLeetcodeCookies(ctx, i.Store.DB(), userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return authedImportResult{skipped: true}, nil
		}
		return authedImportResult{}, err
	}
	if !c.Valid || len(c.SessionEnc) == 0 || len(c.CSRFEnc) == 0 {
		return authedImportResult{skipped: true}, nil
	}
	session, err := i.Vault.Open(c.SessionEnc)
	if err != nil {
		_ = store.MarkCookiesInvalid(ctx, i.Store.DB(), userID)
		return authedImportResult{skipped: true}, nil
	}
	csrf, err := i.Vault.Open(c.CSRFEnc)
	if err != nil {
		_ = store.MarkCookiesInvalid(ctx, i.Store.DB(), userID)
		return authedImportResult{skipped: true}, nil
	}

	pageSize := i.SubmissionPageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}
	maxPages := i.MaxSubmissionPages
	if maxPages <= 0 || maxPages > 100 {
		maxPages = 20
	}

	creds := leetcode.Credentials{Session: string(session), CSRF: string(csrf)}
	var out authedImportResult
	solved, err := i.importSolvedProblems(ctx, userID, creds, now)
	if errors.Is(err, leetcode.ErrAuthExpired) {
		_ = store.MarkCookiesInvalid(ctx, i.Store.DB(), userID)
		return out, nil
	}
	if err != nil {
		return out, err
	}
	out.solved = solved

	lastKey := ""
	spreadIdx := 20
	for page := 0; page < maxPages; page++ {
		p, err := i.Client.ListSubmissions(ctx, creds, page*pageSize, pageSize, lastKey, "")
		if errors.Is(err, leetcode.ErrAuthExpired) {
			_ = store.MarkCookiesInvalid(ctx, i.Store.DB(), userID)
			return out, nil
		}
		if err != nil {
			return out, err
		}
		for _, sub := range p.Submissions {
			if leetcode.NormalizeVerdict(sub.StatusDisplay) != leetcode.VerdictCodeAC {
				continue
			}
			status, err := i.importAccepted(ctx, userID, sub.TitleSlug, sub.ID, parseOrNow(sub.Timestamp, now), coldStartDueAt(now, spreadIdx))
			if err != nil {
				return out, err
			}
			switch status {
			case importInserted:
				out.imported++
				spreadIdx++
			case importDuplicate:
				out.duplicates++
			case importUnknownProblem:
				out.unknown++
			}
		}
		if !p.HasNext || p.LastKey == "" {
			break
		}
		lastKey = p.LastKey
		time.Sleep(1100 * time.Millisecond)
	}
	return out, nil
}

func (i *ColdStartImporter) importSolvedProblems(ctx context.Context, userID int64, creds leetcode.Credentials, now time.Time) (int, error) {
	const pageSize = 50
	imported := 0
	for skip := 0; ; skip += pageSize {
		items, total, err := i.Client.ListSolvedProblems(ctx, creds, skip, pageSize)
		if err != nil {
			return imported, err
		}
		for _, item := range items {
			p, err := store.GetProblemBySlug(ctx, i.Store.DB(), item.TitleSlug)
			if errors.Is(err, store.ErrNotFound) {
				continue
			}
			if err != nil {
				return imported, err
			}
			due := coldStartSolvedDueAt(now, imported)
			inserted, err := store.InsertSolvedUserProblemIfMissing(ctx, i.Store.DB(), models.UserProblem{
				UserID:        userID,
				ProblemID:     p.ID,
				EaseFactor:    2.5,
				IntervalDays:  7,
				NextDueAt:     &due,
				TotalAttempts: 1,
				CleanSolves:   1,
				Status:        models.StatusReview,
			})
			if err != nil {
				return imported, err
			}
			if inserted {
				imported++
			}
		}
		if len(items) == 0 || skip+len(items) >= total {
			break
		}
		time.Sleep(1100 * time.Millisecond)
	}
	return imported, nil
}

func solvedCount(u *leetcode.MatchedUser) int {
	for _, c := range u.SubmitStats.ACSubmissionNum {
		if strings.EqualFold(c.Difficulty, "All") {
			return c.Count
		}
	}
	total := 0
	for _, c := range u.SubmitStats.ACSubmissionNum {
		if !strings.EqualFold(c.Difficulty, "All") {
			total += c.Count
		}
	}
	return total
}

func reverseRecent(in []leetcode.RecentSubmission) []leetcode.RecentSubmission {
	out := make([]leetcode.RecentSubmission, len(in))
	for i := range in {
		out[len(in)-1-i] = in[i]
	}
	return out
}

func parseOrNow(raw string, now time.Time) time.Time {
	if t, ok := parseUnixSeconds(raw); ok {
		return t
	}
	return now
}

func parseUnixSeconds(raw string) (time.Time, bool) {
	sec, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || sec <= 0 {
		return time.Time{}, false
	}
	return time.Unix(sec, 0).UTC(), true
}

func coldStartDueAt(now time.Time, idx int) time.Time {
	if idx < 0 {
		idx = 0
	}
	offset := idx * 13 / 19
	if offset > 13 {
		offset = 13
	}
	return now.Add(time.Duration(offset+1) * 24 * time.Hour)
}

func coldStartSolvedDueAt(now time.Time, idx int) time.Time {
	if idx < 0 {
		idx = 0
	}
	offset := idx % 30
	return now.Add(time.Duration(offset+1) * 24 * time.Hour)
}

func setLeetcodeUsername(ctx context.Context, db store.DBTX, userID int64, username string) error {
	_, err := db.Exec(ctx, `UPDATE users SET leetcode_username = $2 WHERE id = $1`, userID, username)
	if err != nil {
		return fmt.Errorf("set leetcode username: %w", err)
	}
	return nil
}

func setNextDueAt(ctx context.Context, db store.DBTX, userID, problemID int64, dueAt time.Time) error {
	_, err := db.Exec(ctx, `UPDATE user_problems SET next_due_at = $3 WHERE user_id = $1 AND problem_id = $2`, userID, problemID, dueAt)
	if err != nil {
		return fmt.Errorf("set next due: %w", err)
	}
	return nil
}
