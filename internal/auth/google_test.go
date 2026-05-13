package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoogleUserInfoDecodesVerifiedEmail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer access-token" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sub":"google-sub","email":"USER@example.COM","email_verified":true}`))
	}))
	defer srv.Close()

	oauth := &GoogleOAuth{UserInfoEndpoint: srv.URL, HTTPClient: srv.Client()}

	user, err := oauth.userInfo(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("userInfo() error = %v", err)
	}
	if !user.EmailVerified {
		t.Fatalf("EmailVerified = false, want true")
	}
}
