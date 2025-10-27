package repository

// # Работа с базой данных

import (
	"accounting-core/internal/domain"
	"database/sql"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) InitDB() error {
	_, err := r.db.Exec(`
        CREATE TABLE IF NOT EXISTS accounts (
            id SERIAL PRIMARY KEY,
            number VARCHAR(50) UNIQUE NOT NULL,
            balance DECIMAL(15,2) DEFAULT 0
        )
    `)
	return err
}

func (r *PostgresRepo) SeedData() error {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count)
	if count == 0 {
		r.db.Exec("INSERT INTO accounts (number, balance) VALUES ('1001', 1000.00)")
		r.db.Exec("INSERT INTO accounts (number, balance) VALUES ('1002', 500.00)")
	}
	return nil
}

func (r *PostgresRepo) GetAccounts() ([]domain.Account, error) {
	rows, err := r.db.Query("SELECT id, number, balance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var acc domain.Account
		err := rows.Scan(&acc.ID, &acc.Number, &acc.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *PostgresRepo) TransferMoney(from, to int64, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
