package srs

import (
	"testing"
	"time"
)

// ---- DeriveRating ----

func TestDeriveRating(t *testing.T) {
	tests := []struct {
		name string
		in   Outcome
		want Rating
	}{
		{
			name: "first-try AC under expected time = Strong",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 1, TimeTakenSec: 10 * 60, Difficulty: DifficultyMedium},
			want: RatingStrong,
		},
		{
			name: "first-try AC at exactly expected time = Strong",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 1, TimeTakenSec: 25 * 60, Difficulty: DifficultyMedium},
			want: RatingStrong,
		},
		{
			name: "first-try AC over expected time = Normal",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 1, TimeTakenSec: 35 * 60, Difficulty: DifficultyMedium},
			want: RatingNormal,
		},
		{
			name: "two submissions reasonable time = Normal",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 2, TimeTakenSec: 20 * 60, Difficulty: DifficultyMedium},
			want: RatingNormal,
		},
		{
			name: "two submissions way over time = Struggled",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 2, TimeTakenSec: 60 * 60, Difficulty: DifficultyMedium},
			want: RatingStruggled,
		},
		{
			name: "three submissions = Struggled",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 3, TimeTakenSec: 15 * 60, Difficulty: DifficultyMedium},
			want: RatingStruggled,
		},
		{
			name: "five submissions still Struggled",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 5, TimeTakenSec: 30 * 60, Difficulty: DifficultyMedium},
			want: RatingStruggled,
		},
		{
			name: "six submissions = Failed (brute-forced via WAs)",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 6, TimeTakenSec: 30 * 60, Difficulty: DifficultyMedium},
			want: RatingFailed,
		},
		{
			name: "non-AC verdict = Failed regardless of subs",
			in:   Outcome{Verdict: VerdictWA, SubmissionCount: 1, TimeTakenSec: 5 * 60, Difficulty: DifficultyEasy},
			want: RatingFailed,
		},
		{
			name: "TLE = Failed",
			in:   Outcome{Verdict: VerdictTLE, SubmissionCount: 3, TimeTakenSec: 30 * 60, Difficulty: DifficultyHard},
			want: RatingFailed,
		},
		{
			name: "hard expected time is more generous",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 1, TimeTakenSec: 35 * 60, Difficulty: DifficultyHard},
			want: RatingStrong, // 35 min on a Hard is under the 40-min expected
		},
		{
			name: "easy has tighter expected time",
			in:   Outcome{Verdict: VerdictAC, SubmissionCount: 1, TimeTakenSec: 20 * 60, Difficulty: DifficultyEasy},
			want: RatingNormal, // 20 min on an Easy is over the 15-min expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveRating(tt.in)
			if got != tt.want {
				t.Errorf("DeriveRating(%+v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// ---- NextState: first attempt from "new" ----

func TestNextState_FirstAttempt(t *testing.T) {
	init := NewState()

	t.Run("first strong solve graduates to review with interval 3", func(t *testing.T) {
		got := NextState(init, RatingStrong)
		if got.IntervalDays != 3 {
			t.Errorf("interval = %d, want 3", got.IntervalDays)
		}
		if got.Streak != 1 {
			t.Errorf("streak = %d, want 1", got.Streak)
		}
		if got.Status != StatusReview {
			t.Errorf("status = %v, want %v", got.Status, StatusReview)
		}
		if got.EaseFactor <= 2.5 {
			t.Errorf("ease = %f, want > 2.5 (Strong bumps it up)", got.EaseFactor)
		}
	})

	t.Run("first failed attempt = interval 1, streak 0, fail counted", func(t *testing.T) {
		got := NextState(init, RatingFailed)
		if got.IntervalDays != 1 {
			t.Errorf("interval = %d, want 1", got.IntervalDays)
		}
		if got.Streak != 0 {
			t.Errorf("streak = %d, want 0", got.Streak)
		}
		if got.TotalFails != 1 {
			t.Errorf("fails = %d, want 1", got.TotalFails)
		}
		if got.EaseFactor >= 2.5 {
			t.Errorf("ease = %f, want < 2.5 (Failed drops it)", got.EaseFactor)
		}
	})
}

// ---- NextState: ease factor clamping ----

func TestNextState_EaseClamp(t *testing.T) {
	t.Run("ease never drops below 1.3", func(t *testing.T) {
		s := State{EaseFactor: 1.35, IntervalDays: 5, Status: StatusReview}
		// Two failures should not drive ease below 1.3.
		s = NextState(s, RatingFailed)
		s = NextState(s, RatingFailed)
		if s.EaseFactor < minEase {
			t.Errorf("ease = %f, want >= %f", s.EaseFactor, minEase)
		}
	})

	t.Run("ease never exceeds 3.0", func(t *testing.T) {
		s := State{EaseFactor: 2.9, IntervalDays: 30, Status: StatusReview}
		// Many strong solves shouldn't push ease above 3.0.
		for i := 0; i < 5; i++ {
			s = NextState(s, RatingStrong)
		}
		if s.EaseFactor > maxEase {
			t.Errorf("ease = %f, want <= %f", s.EaseFactor, maxEase)
		}
	})
}

// ---- NextState: interval cap ----

func TestNextState_IntervalCap(t *testing.T) {
	s := State{EaseFactor: 3.0, IntervalDays: 100, Status: StatusMastered}
	got := NextState(s, RatingStrong)
	if got.IntervalDays > maxIntervalDays {
		t.Errorf("interval = %d, want <= %d", got.IntervalDays, maxIntervalDays)
	}
}

// ---- NextState: failure resets interval ----

func TestNextState_FailureResetsInterval(t *testing.T) {
	s := State{EaseFactor: 2.5, IntervalDays: 30, Streak: 4, Status: StatusReview}
	got := NextState(s, RatingFailed)

	if got.IntervalDays != 1 {
		t.Errorf("interval after failure = %d, want 1", got.IntervalDays)
	}
	if got.Streak != 0 {
		t.Errorf("streak after failure = %d, want 0", got.Streak)
	}
	if got.TotalFails != 1 {
		t.Errorf("fails after failure = %d, want 1", got.TotalFails)
	}
}

// ---- NextState: struggled keeps interval growing but slowly ----

func TestNextState_Struggled(t *testing.T) {
	s := State{EaseFactor: 2.5, IntervalDays: 10, Streak: 3, Status: StatusReview}
	got := NextState(s, RatingStruggled)

	if got.Streak != 0 {
		t.Errorf("struggle should reset streak, got %d", got.Streak)
	}
	if got.IntervalDays >= 10 {
		t.Errorf("struggle should shrink/stall interval, got %d (was 10)", got.IntervalDays)
	}
	if got.IntervalDays < 1 {
		t.Errorf("interval must be at least 1, got %d", got.IntervalDays)
	}
}

// ---- Mastery promotion ----

func TestNextState_PromotesToMastered(t *testing.T) {
	// Drive a problem to a high interval through repeated strong solves.
	s := NewState()
	for i := 0; i < 8; i++ {
		s = NextState(s, RatingStrong)
	}
	if s.Status != StatusMastered {
		t.Errorf("after 8 strong solves status = %v, want %v (interval=%d, ease=%f)",
			s.Status, StatusMastered, s.IntervalDays, s.EaseFactor)
	}
}

// ---- Full trajectory test ----

func TestNextState_RealisticTrajectory(t *testing.T) {
	// Story: user attempts a Medium problem.
	//   day 0: fails it
	//   day 1: solves with 3 submissions (struggled)
	//   day 2: solves with 2 submissions, normal time (normal)
	//   day ~5: solves cleanly first try (strong)
	//   day ~14: solves cleanly first try (strong)
	// We're checking intervals grow sensibly and status progresses.

	s := NewState()
	type step struct {
		rating       Rating
		wantMinIvl   int
		wantMaxIvl   int
		wantStatus   Status
		wantNonZeroF int // expected fails
	}
	steps := []step{
		{RatingFailed, 1, 1, StatusReview, 1},
		{RatingStruggled, 1, 2, StatusReview, 1},
		{RatingNormal, 2, 5, StatusReview, 1},
		{RatingStrong, 4, 10, StatusReview, 1},
		{RatingStrong, 10, 30, StatusReview, 1},
	}

	for i, st := range steps {
		s = NextState(s, st.rating)
		if s.IntervalDays < st.wantMinIvl || s.IntervalDays > st.wantMaxIvl {
			t.Errorf("step %d (%v): interval = %d, want in [%d, %d]",
				i, st.rating, s.IntervalDays, st.wantMinIvl, st.wantMaxIvl)
		}
		if s.Status != st.wantStatus {
			t.Errorf("step %d (%v): status = %v, want %v",
				i, st.rating, s.Status, st.wantStatus)
		}
		if s.TotalFails != st.wantNonZeroF {
			t.Errorf("step %d (%v): fails = %d, want %d",
				i, st.rating, s.TotalFails, st.wantNonZeroF)
		}
	}
}

// ---- NextDueAt sanity ----

func TestNextDueAt(t *testing.T) {
	base := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	s := State{IntervalDays: 7}
	due := NextDueAt(s, base)
	expected := base.Add(7 * 24 * time.Hour)
	if !due.Equal(expected) {
		t.Errorf("NextDueAt = %v, want %v", due, expected)
	}
}
