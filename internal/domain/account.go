package domain

// # Бизнес-сущности (модели данных)

type Account struct {
	ID      int64   `json:"id"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	DebitAccountID  int64   `json:"debit_account_id"`
	CreditAccountID int64   `json:"credit_account_id"`
	Amount          float64 `json:"amount"`
}
