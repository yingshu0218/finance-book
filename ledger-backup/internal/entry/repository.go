package entry

import (
	"database/sql"
	"fmt"
	"time"
)

type Entry struct {
	ID        int     `json:"id"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	Date      string  `json:"date"`
	Note      string  `json:"note"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(e *Entry) error {
	e.CreatedAt = time.Now().Format(time.RFC3339)
	e.UpdatedAt = e.CreatedAt

	_, err := r.db.Exec(
		"INSERT INTO entries (amount, category, date, note, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		e.Amount, e.Category, e.Date, e.Note, e.CreatedAt, e.UpdatedAt,
	)
	return err
}

func (r *Repository) List(month, category string) ([]Entry, error) {
	query := "SELECT id, amount, category, date, note, created_at, updated_at FROM entries"
	params := []interface{}{}

	if month != "" {
		query += " WHERE date LIKE ?"
		params = append(params, month+"-%")
	}

	if category != "" {
		if len(params) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " category = ?"
		params = append(params, category)
	}

	query += " ORDER BY date DESC, id DESC"

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Date, &e.Note, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (r *Repository) GetByID(id int) (*Entry, error) {
	var e Entry
	err := r.db.QueryRow(
		"SELECT id, amount, category, date, note, created_at, updated_at FROM entries WHERE id = ?",
		id,
	).Scan(&e.ID, &e.Amount, &e.Category, &e.Date, &e.Note, &e.CreatedAt, &e.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrEntryNotFound
	}
	return &e, err
}

func (r *Repository) Update(e *Entry) error {
	e.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err := r.db.Exec(
		"UPDATE entries SET amount = ?, category = ?, date = ?, note = ?, updated_at = ? WHERE id = ?",
		e.Amount, e.Category, e.Date, e.Note, e.UpdatedAt, e.ID,
	)
	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM entries WHERE id = ?", id)
	return err
}

func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

var ErrEntryNotFound = fmt.Errorf("entry not found")
