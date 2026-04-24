package middleware

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func generateToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func newContextWithToken(t *testing.T, token string) context.Context {
	t.Helper()
	md := metadata.New(nil)
	md.Set("authorization", "Bearer "+token)
	return metadata.NewServerContext(context.Background(), md)
}

func TestServerAuth_ValidToken(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id":  float64(123),
		"username": "testuser",
		"role":     "admin",
	}
	token := generateToken(t, claims)
	ctx := newContextWithToken(t, token)

	mw := ServerAuth(testSecret)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		info, ok := GetAuthInfo(ctx)
		if !ok {
			t.Fatal("expected AuthInfo in context")
		}
		if info.UserID != 123 {
			t.Fatalf("expected UserID 123, got %d", info.UserID)
		}
		if info.Username != "testuser" {
			t.Fatalf("expected username testuser, got %s", info.Username)
		}
		if info.Role != "admin" {
			t.Fatalf("expected role admin, got %s", info.Role)
		}
		return "ok", nil
	}

	resp, err := mw(handler)(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected resp 'ok', got %v", resp)
	}
}

func TestServerAuth_EmptyToken_Rejected(t *testing.T) {
	ctx := newContextWithToken(t, "")

	mw := ServerAuth(testSecret)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	_, err := mw(handler)(ctx, nil)
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestServerAuth_EmptyToken_Allowed(t *testing.T) {
	ctx := newContextWithToken(t, "")

	mw := ServerAuth(testSecret, WithAllowEmptyToken())
	called := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		called = true
		return "ok", nil
	}

	resp, err := mw(handler)(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("handler should have been called")
	}
	if resp != "ok" {
		t.Fatalf("expected resp 'ok', got %v", resp)
	}
}

func TestServerAuth_InvalidToken(t *testing.T) {
	ctx := newContextWithToken(t, "invalid.token.here")

	mw := ServerAuth(testSecret)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	_, err := mw(handler)(ctx, nil)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestServerAuth_SigningMethodCheck(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id":  float64(123),
		"username": "testuser",
		"role":     "admin",
	}
	// Use none signing method
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	ctx := newContextWithToken(t, s)

	mw := ServerAuth(testSecret, WithSigningMethodCheck())
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	_, err := mw(handler)(ctx, nil)
	if err == nil {
		t.Fatal("expected error for none signing method")
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name     string
		auth     string
		expected string
	}{
		{"with bearer prefix", "Bearer mytoken", "mytoken"},
		{"without prefix", "mytoken", "mytoken"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := metadata.New(nil)
			md.Set("authorization", tt.auth)
			ctx := metadata.NewServerContext(context.Background(), md)
			got := ExtractToken(ctx)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetAuthInfo(t *testing.T) {
	ctx := context.Background()

	// Without auth info
	_, ok := GetAuthInfo(ctx)
	if ok {
		t.Fatal("expected no auth info")
	}

	// With auth info
	info := &AuthInfo{UserID: 123, Username: "test", Role: "admin"}
	ctx = context.WithValue(ctx, authKey{}, info)
	got, ok := GetAuthInfo(ctx)
	if !ok {
		t.Fatal("expected auth info")
	}
	if got.UserID != 123 {
		t.Fatalf("expected UserID 123, got %d", got.UserID)
	}
}