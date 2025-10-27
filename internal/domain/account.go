package domain

import "time"

type Account struct {
	ID      int64     `json:"id"`
	Number  string    `json:"number"`
	Balance float64   `json:"balance"`
	Created time.Time `json:"created_at"`
}

type Transaction struct {
	ID              int64     `json:"id"`
	DebitAccountID  int64     `json:"debit_account_id"`  // Счет дебета
	CreditAccountID int64     `json:"credit_account_id"` // Счет кредита
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}
