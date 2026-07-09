package book

import (
	"os"
	"testing"

	"ledger/internal/config"
)

func TestCreateBook(t *testing.T) {
	// Use a temporary config directory for testing
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	m := NewManager(cfg)

	err = m.CreateBook("test-book")
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	err = m.CreateBook("test-book")
	if err != ErrBookExists {
		t.Fatal("expected ErrBookExists")
	}

	// Cleanup
	m.DeleteBook("test-book")
}
