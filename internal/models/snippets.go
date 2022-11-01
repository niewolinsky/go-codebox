package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Snippet struct {
	SnippetID int
	UserID    int
	Title     string
	Content   string
	Created   time.Time
	Expires   time.Time
	Public    bool
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, public bool) (int, error) {
	stmt := `INSERT INTO snippets (user_id, title, content, created, expires, public) VALUES (1, $1::varchar, $2::varchar, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '3 day', $3) RETURNING snippet_id`
	var snippet_id int

	err := m.DB.QueryRow(context.Background(), stmt, title, content, public).Scan(&snippet_id)
	if err != nil {
		return -1, err
	}

	return int(snippet_id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT snippet_id, user_id, title, content, created, expires, public FROM snippets WHERE snippet_id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	s := &Snippet{}
	err := row.Scan(&s.SnippetID, &s.UserID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT snippet_id, user_id, title, content, created, expires, public FROM snippets WHERE public = true LIMIT 5`

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sArr := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.SnippetID, &s.UserID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Public)
		if err != nil {
			return nil, err
		}

		sArr = append(sArr, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sArr, nil
}
