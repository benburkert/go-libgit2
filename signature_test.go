package libgit2

import (
	"testing"
	"time"
)

func TestDefaultSignature(t *testing.T) {
	repo := mustInitTestRepo(t)

	now := time.Now().Truncate(time.Second)
	sig, err := repo.DefaultSignature()
	if err != nil {
		t.Fatal(err)
	}

	if sig.Name != "Default" {
		t.Errorf("want sig name %q, got %q", "Default", sig.Name)
	}

	if sig.Email != "default@example.com" {
		t.Errorf("want sig email %q, got %q", "default@example.com", sig.Email)
	}

	if !sig.When.Equal(now) {
		t.Errorf("want sig when %q, got %q", now, sig.When)
	}
}
