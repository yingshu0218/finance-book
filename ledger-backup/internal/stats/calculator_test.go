package stats

import (
	"os"
	"testing"

	"ledger/internal/book"
	"ledger/internal/config"
	"ledger/internal/storage"
)

func TestGetBalance(t *testing.T) {
	// Use a temporary config directory for testing
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Ensure default book database is initialized
	db, err := storage.OpenDB(cfg.Books[0].Path)
	if err != nil {
		t.Fatalf("failed to open default db: %v", err)
	}
	if err := storage.InitSchema(db); err != nil {
		t.Fatalf("failed to init schema: %v", err)
	}
	db.Close()

	bm := book.NewManager(cfg)
	calc := NewCalculator(bm)

	result, err := calc.GetBalance("default", "")
	if err != nil {
		t.Fatalf("failed to get balance: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
