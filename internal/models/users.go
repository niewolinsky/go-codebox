package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID         int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, $4)`

	_, err = m.DB.Exec(context.Background(), stmt, name, email, string(hashedPassword), time.Now())
	if err != nil {
		// rewrite database, use pgx directly not through db.sql, this does not checks for returned errors by psql
		return ErrDuplicateEmail
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT user_id, hashed_password FROM users WHERE email = $1"
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE user_id = $1)"
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&exists)

	return exists, err
}

func (m *UserModel) EmailTaken(email string) bool {
	stmt := `SELECT email FROM snippets WHERE email = $1`

	row := m.DB.QueryRow(context.Background(), stmt, email)

	u := &User{}
	err := row.Scan(&u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}

	return true
}
