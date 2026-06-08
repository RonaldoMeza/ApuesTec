package auth

import (
	"testing"
)

func TestHashToken(t *testing.T) {
	token := "test-refresh-token"
	hash1 := hashToken(token)
	hash2 := hashToken(token)
	if hash1 != hash2 {
		t.Errorf("hashToken should be deterministic, got %s != %s", hash1, hash2)
	}
	if hash1 == "" {
		t.Error("hashToken should not return empty string")
	}
}

func TestGenerateRefreshTokenPair(t *testing.T) {
	raw, hash, err := generateRefreshTokenPair()
	if err != nil {
		t.Fatalf("generateRefreshTokenPair failed: %v", err)
	}
	if raw == "" {
		t.Error("raw token should not be empty")
	}
	if hash == "" {
		t.Error("hash should not be empty")
	}
	if raw == hash {
		t.Error("raw token should differ from hash")
	}
	if len(raw) != 64 {
		t.Errorf("raw token should be 64 hex chars, got %d", len(raw))
	}
	if hash != hashToken(raw) {
		t.Error("hash should match hashToken(raw)")
	}
}

func TestHashTokenConsistency(t *testing.T) {
	tokens := []string{"", "a", "abc", "hello-world-token-123"}
	for _, tok := range tokens {
		h1 := hashToken(tok)
		h2 := hashToken(tok)
		if h1 != h2 {
			t.Errorf("hashToken inconsistent for %q: %s != %s", tok, h1, h2)
		}
	}
}
