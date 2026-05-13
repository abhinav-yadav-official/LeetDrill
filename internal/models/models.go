// Package models contains the domain types shared between store, http, and
// other internal packages. Kept thin — no DB tags, no JSON tags unless
// needed at boundaries.
package models

import "time"

type Difficulty string

const (
	DifficultyEasy   Difficulty = "Easy"
	DifficultyMedium Difficulty = "Medium"
	DifficultyHard   Difficulty = "Hard"
)

type Status string

const (
	StatusNew      Status = "new"
	StatusReview   Status = "review"
	StatusMastered Status = "mastered"
)

type Tag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Problem struct {
	ID                 int64
	LeetcodeSlug       string
	LeetcodeQuestionID string
	LeetcodeFrontendID string
	Title              string
	Difficulty         Difficulty
	URL                string
	ContentHTML        string
	TopicTags          []Tag
	ACRate             float64
	PaidOnly           bool
	SyncedAt           time.Time
}

type Pattern struct {
	ID          int64
	Slug        string
	Name        string
	Description string
}

type UserProblem struct {
	UserID          int64
	ProblemID       int64
	EaseFactor      float64
	IntervalDays    int
	NextDueAt       *time.Time
	LastAttemptedAt *time.Time
	TotalAttempts   int
	CleanSolves     int
	TotalFails      int
	Streak          int
	Status          Status
}

type Attempt struct {
	ID                       int64
	UserID                   int64
	ProblemID                int64
	StartedAt                *time.Time
	CompletedAt              time.Time
	Verdict                  string
	SubmissionCountInSession int
	TimeTakenSec             int
	RuntimeMs                *int
	MemoryKB                 *int
	Language                 string
	Code                     string
	DerivedRating            string
	Journal                  string
	MistakeTags              []string
	LeetcodeSubmissionID     string
}
