package stats

import (
	"ledger/internal/book"
	"ledger/internal/storage"
)

type BalanceResult struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type Calculator struct {
	bookManager *book.Manager
}

func NewCalculator(bm *book.Manager) *Calculator {
	return &Calculator{bookManager: bm}
}

func (c *Calculator) GetBalance(bookName, month string) (*BalanceResult, error) {
	path, err := c.bookManager.GetBookPath(bookName)
	if err != nil {
		return nil, err
	}

	db, err := storage.OpenDB(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT amount FROM entries"
	params := []interface{}{}

	if month != "" {
		query += " WHERE date LIKE ?"
		params = append(params, month+"-%")
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var income, expense float64
	for rows.Next() {
		var amount float64
		if err := rows.Scan(&amount); err != nil {
			return nil, err
		}

		if amount > 0 {
			income += amount
		} else {
			expense += amount
		}
	}

	return &BalanceResult{
		Income:  income,
		Expense: -expense,
		Balance: income + expense,
	}, nil
}
