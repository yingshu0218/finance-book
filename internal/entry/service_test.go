package entry

import (
	"os"
	"testing"

	"ledger/internal/book"
	"ledger/internal/config"
	"ledger/internal/storage"
)

func TestAddEntry(t *testing.T) {
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
	svc := NewService(bm)

	err = svc.Add("default", -50, "餐饮招待", "", "午餐")
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	entries, err := svc.List("default", "", "")
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Amount != -50 {
		t.Fatalf("expected amount -50, got %f", entries[0].Amount)
	}
}
