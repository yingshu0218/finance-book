package book

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ledger/internal/config"
	"ledger/internal/storage"
)

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) GetDefaultBook() string {
	return m.cfg.DefaultBook
}

func (m *Manager) ListBooks() []config.BookInfo {
	return m.cfg.Books
}

func (m *Manager) CreateBook(name string) error {
	for _, b := range m.cfg.Books {
		if b.Name == name {
			return ErrBookExists
		}
	}

	bookPath := filepath.Join(config.GetBooksDir(), name+".db")

	db, err := storage.OpenDB(bookPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := storage.InitSchema(db); err != nil {
		return err
	}

	m.cfg.Books = append(m.cfg.Books, config.BookInfo{
		Name:      name,
		Path:      bookPath,
		CreatedAt: time.Now().Format("2006-01-02"),
	})

	return m.cfg.Save()
}

func (m *Manager) DeleteBook(name string) error {
	for i, b := range m.cfg.Books {
		if b.Name == name {
			if err := os.Remove(b.Path); err != nil && !os.IsNotExist(err) {
				return err
			}

			m.cfg.Books = append(m.cfg.Books[:i], m.cfg.Books[i+1:]...)

			if m.cfg.DefaultBook == name && len(m.cfg.Books) > 0 {
				m.cfg.DefaultBook = m.cfg.Books[0].Name
			}

			return m.cfg.Save()
		}
	}
	return ErrBookNotFound
}

func (m *Manager) GetBookPath(name string) (string, error) {
	for _, b := range m.cfg.Books {
		if b.Name == name {
			return b.Path, nil
		}
	}
	return "", ErrBookNotFound
}

var (
	ErrBookExists   = fmt.Errorf("book already exists")
	ErrBookNotFound = fmt.Errorf("book not found")
)
