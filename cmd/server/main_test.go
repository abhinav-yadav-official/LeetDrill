package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoginPageUsesSplitIntroLayout(t *testing.T) {
	body := fmt.Sprintf(loginPage, "invalid email or password.", "/leetdrill/login")

	for _, want := range []string{
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`Daily review flow for LeetCode practice.`,
		`Track recent submissions, spaced repetition, and difficult problems from one focused workspace.`,
		`invalid email or password.`,
		`action="/leetdrill/login"`,
		`type="email"`,
		`type="password"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("login page missing %q:\n%s", want, body)
		}
	}
}
