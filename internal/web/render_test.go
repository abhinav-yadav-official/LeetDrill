package web

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRendererLoadsCorePagesAndPartials(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	for _, name := range []string{
		"dashboard",
		"problems",
		"problem_detail",
		"patterns",
		"stats",
		"settings",
		"session_today",
	} {
		if _, ok := r.pages[name]; !ok {
			t.Fatalf("page %q not loaded", name)
		}
	}

	for _, name := range []string{"session_card", "problem_row"} {
		if _, ok := r.partials[name]; !ok {
			t.Fatalf("partial %q not loaded", name)
		}
	}
}

func TestRendererPageIncludesHTMXShell(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	rec := httptest.NewRecorder()
	r.Page(rec, "dashboard", PageData{
		Title:   "Dashboard",
		UserID:  7,
		NavItem: "dashboard",
		Data:    map[string]string{"Now": "ok"},
	})

	body := rec.Body.String()
	for _, want := range []string{
		`<script src="https://unpkg.com/htmx.org`,
		`href="/session/today"`,
		`LeetDrill`,
		`Today`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("rendered dashboard missing %q:\n%s", want, body)
		}
	}
}
