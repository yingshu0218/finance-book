package app

import (
	"ledger/internal/book"
	"ledger/internal/config"
	"ledger/internal/entry"
	"ledger/internal/stats"
	"ledger/internal/storage"
)

var (
	cfg             *config.Config
	bookManager     *book.Manager
	entryService    *entry.Service
	statsCalculator *stats.Calculator
)

func Init() error {
	var err error
	cfg, err = config.Load()
	if err != nil {
		return err
	}

	// Ensure all configured book databases are initialized
	for _, b := range cfg.Books {
		db, err := storage.OpenDB(b.Path)
		if err != nil {
			return err
		}
		if err := storage.InitSchema(db); err != nil {
			db.Close()
			return err
		}
		db.Close()
	}

	bookManager = book.NewManager(cfg)
	entryService = entry.NewService(bookManager)
	statsCalculator = stats.NewCalculator(bookManager)

	return nil
}

func GetDefaultBook() string {
	return bookManager.GetDefaultBook()
}

func ListBooks() []config.BookInfo {
	return bookManager.ListBooks()
}

func CreateBook(name string) error {
	return bookManager.CreateBook(name)
}

func DeleteBook(name string) error {
	return bookManager.DeleteBook(name)
}

func AddEntry(bookName string, amount float64, category, date, note string) error {
	return entryService.Add(bookName, amount, category, date, note)
}

func ListEntries(bookName, month, category string) ([]entry.Entry, error) {
	return entryService.List(bookName, month, category)
}

func GetEntry(bookName string, id int) (*entry.Entry, error) {
	return entryService.Get(bookName, id)
}

func UpdateEntry(bookName string, id int, amount float64, category, date, note string) error {
	return entryService.Update(bookName, id, amount, category, date, note)
}

func DeleteEntry(bookName string, id int) error {
	return entryService.Delete(bookName, id)
}

func GetBalance(bookName, month string) (*stats.BalanceResult, error) {
	return statsCalculator.GetBalance(bookName, month)
}
