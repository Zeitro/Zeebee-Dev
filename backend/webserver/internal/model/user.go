package model

import (
	"time"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
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
	// db *sql.DB
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {

	return &UserModel{db: db}
}

// func (m *UserModel) Close() error {
// 	return m.db.Close()
// }
