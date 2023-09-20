package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	s := Snippet{}

	err := m.DB.
		QueryRow(`SELECT id, title, content, created, expires from snippets	WHERE expires > UTC_TIMESTAMP() and id=?`, id).
		Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if nerr := normalizeDbError(err); nerr != nil {
		return nil, nerr
	}

	return &s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {

	rows, err := m.DB.Query(`SELECT id, title, content, created, expires 
	FROM snippets where expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoRecord
	}

	if nerr := normalizeDbError(err); nerr != nil {
		return nil, nerr
	}

	defer rows.Close()
	snippets := []*Snippet{}

	for rows.Next() {
		s := Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func normalizeDbError(e error) error {
	if errors.Is(e, sql.ErrNoRows) {
		return ErrNoRecord
	}

	return e
}
