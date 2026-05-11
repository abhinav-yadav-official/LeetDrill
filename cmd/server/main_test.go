package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoginPageUsesSplitIntroLayout(t *testing.T) {
	body := fmt.Sprintf(loginPage, "invalid email or password.", "/leetdrill/login", "/leetdrill/signup")

	for _, want := range []string{
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`Daily review flow for LeetCode practice.`,
		`Track recent submissions, spaced repetition, and difficult problems from one focused workspace.`,
		`invalid email or password.`,
		`action="/leetdrill/login"`,
		`href="/leetdrill/signup"`,
		`type="email"`,
		`type="password"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("login page missing %q:\n%s", want, body)
		}
	}
}

func TestSignupPageUsesSplitIntroLayout(t *testing.T) {
	body := fmt.Sprintf(signupPage, "create an account.", "/leetdrill/signup", "/leetdrill/login")

	for _, want := range []string{
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`Daily review flow for LeetCode practice.`,
		`Create account`,
		`action="/leetdrill/signup"`,
		`href="/leetdrill/login"`,
		`name="email"`,
		`name="password"`,
		`name="confirm_password"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("signup page missing %q:\n%s", want, body)
		}
	}
}

func TestValidateSignupForm(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		confirm     string
		wantEmail   string
		wantMessage string
	}{
		{
			name:      "valid normalized email",
			email:     "  USER@example.COM ",
			password:  "correct horse",
			confirm:   "correct horse",
			wantEmail: "user@example.com",
		},
		{
			name:        "email required",
			password:    "correct horse",
			confirm:     "correct horse",
			wantMessage: "email is required.",
		},
		{
			name:        "password length",
			email:       "user@example.com",
			password:    "short",
			confirm:     "short",
			wantMessage: "password must be at least 8 characters.",
		},
		{
			name:        "confirmation mismatch",
			email:       "user@example.com",
			password:    "correct horse",
			confirm:     "wrong horse",
			wantMessage: "passwords do not match.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, gotMessage := validateSignupForm(tt.email, tt.password, tt.confirm)
			if gotEmail != tt.wantEmail || gotMessage != tt.wantMessage {
				t.Fatalf("validateSignupForm() = (%q, %q), want (%q, %q)", gotEmail, gotMessage, tt.wantEmail, tt.wantMessage)
			}
		})
	}
}
