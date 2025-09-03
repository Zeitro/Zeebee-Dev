package model

import (
	"time"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Loan struct {
	ID           int64     `db:"id"`
	BorrowerID   int64     `db:"borrower_id"`
	LenderID     int64     `db:"lender_id"`
	Amount       float64   `db:"amount"`
	InterestRate float64   `db:"interest_rate"`
	Term         int       `db:"term"`
	Status       string    `db:"status"`
	Purpose      string    `db:"purpose"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type LoanModel struct {
	// db *sql.DB
	db *gorm.DB
}

func NewLoanModel(db *gorm.DB) *LoanModel {
	// db, err := sql.Open("postgres", dataSource)
	// if err != nil {
	// 	panic(err)
	// }
	// if err := db.Ping(); err != nil {
	// 	panic(err)
	// }
	return &LoanModel{db: db}
}

// func (m *LoanModel) Close() error {
// 	return m.db.Close()
// }
