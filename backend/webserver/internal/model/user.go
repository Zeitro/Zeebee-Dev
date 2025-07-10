package model

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserModel struct {
	db *sql.DB
}

func NewUserModel(dataSource string) *UserModel {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		panic(err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &UserModel{db: db}
}

func (m *UserModel) Close() error {
	return m.db.Close()
}
