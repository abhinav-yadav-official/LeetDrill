package store

import "testing"

func TestSessionNeedsProblemBuild(t *testing.T) {
	tests := []struct {
		name string
		s    Session
		want bool
	}{
		{
			name: "empty incomplete session rebuilds after ingest",
			s:    Session{},
			want: true,
		},
		{
			name: "session with problems is reused",
			s:    Session{ProblemIDs: []int64{1, 2, 3}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sessionNeedsProblemBuild(tt.s); got != tt.want {
				t.Fatalf("sessionNeedsProblemBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
