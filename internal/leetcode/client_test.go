package leetcode

import (
	"strings"
	"testing"
)

func TestSubmissionListQueryAvoidsRemovedSubmissionCommentField(t *testing.T) {
	if strings.Contains(querySubmissionList, "submissionComment") {
		t.Fatalf("querySubmissionList includes removed LeetCode field submissionComment")
	}
}

func TestSolvedProblemListQueryFiltersAcceptedProblems(t *testing.T) {
	if !strings.Contains(querySolvedProblemList, `status: AC`) {
		t.Fatalf("querySolvedProblemList must filter status AC")
	}
}
