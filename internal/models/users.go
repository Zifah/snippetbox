package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	email    string
	password []byte
	Create   time.Time
}

type UserModel struct {
	DB         *sql.DB
	insertStmt *sql.Stmt
}

func NewUserModel(db *sql.DB) (*UserModel, error) {
	stmt, err := db.Prepare(`INSERT INTO users(name, email, hashed_password, created) 
	VALUES(?, ?, ?, UTC_TIMESTAMP())`)
	if err != nil {
		return nil, err
	}

	return &UserModel{
		DB:         db,
		insertStmt: stmt,
	}, nil
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	result, err := m.insertStmt.Exec(name, email, passwordHash)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	user := User{}
	qry := `SELECT id, hashed_password from users where email = ?`
	err := m.DB.QueryRow(qry, email).Scan(&user.ID, &user.password)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, ErrInvalidCredentials
	}

	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(user.password, []byte(password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidCredentials
	}

	return user.ID, err
}

func (m *UserModel) Exists(id int) (exists bool, err error) {
	stmt := "SELECT EXISTS(SELECT true FROM users where id = ?)"
	err = m.DB.QueryRow(stmt, id).Scan(&exists)
	return
}
