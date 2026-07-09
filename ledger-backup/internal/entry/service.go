package entry

import (
	"time"

	"ledger/internal/book"
	"ledger/internal/storage"
)

type Service struct {
	bookManager *book.Manager
}

func NewService(bm *book.Manager) *Service {
	return &Service{bookManager: bm}
}

func (s *Service) getRepo(bookName string) (*Repository, error) {
	path, err := s.bookManager.GetBookPath(bookName)
	if err != nil {
		return nil, err
	}

	db, err := storage.OpenDB(path)
	if err != nil {
		return nil, err
	}

	return NewRepository(db), nil
}

func (s *Service) Add(bookName string, amount float64, category, date, note string) error {
	repo, err := s.getRepo(bookName)
	if err != nil {
		return err
	}
	defer repo.Close()

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	entry := &Entry{
		Amount:   amount,
		Category: category,
		Date:     date,
		Note:     note,
	}

	return repo.Create(entry)
}

func (s *Service) List(bookName, month, category string) ([]Entry, error) {
	repo, err := s.getRepo(bookName)
	if err != nil {
		return nil, err
	}
	defer repo.Close()

	return repo.List(month, category)
}

func (s *Service) Get(bookName string, id int) (*Entry, error) {
	repo, err := s.getRepo(bookName)
	if err != nil {
		return nil, err
	}
	defer repo.Close()

	return repo.GetByID(id)
}

func (s *Service) Update(bookName string, id int, amount float64, category, date, note string) error {
	repo, err := s.getRepo(bookName)
	if err != nil {
		return err
	}
	defer repo.Close()

	entry, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	entry.Amount = amount
	entry.Category = category
	entry.Date = date
	entry.Note = note

	return repo.Update(entry)
}

func (s *Service) Delete(bookName string, id int) error {
	repo, err := s.getRepo(bookName)
	if err != nil {
		return err
	}
	defer repo.Close()

	return repo.Delete(id)
}
