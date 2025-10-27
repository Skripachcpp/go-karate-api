package repository

import (
	"accounting-core/internal/domain"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

func (r *PostgresRepo) CreateTransaction(tx *domain.Transaction) error {
	// Проверяем, что дебетуемый счет имеет достаточно средств
	var balance float64
	err := r.db.QueryRow("SELECT balance FROM accounts WHERE id = $1", tx.DebitAccountID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < tx.Amount {
		return errors.New("insufficient funds")
	}

	// Выполняем проводку в транзакции
	dbTx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer dbTx.Rollback()

	// Списание с дебетового счета
	_, err = dbTx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2",
		tx.Amount, tx.DebitAccountID)
	if err != nil {
		return err
	}

	// Зачисление на кредитовый счет
	_, err = dbTx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		tx.Amount, tx.CreditAccountID)
	if err != nil {
		return err
	}

	// Сохраняем проводку
	_, err = dbTx.Exec(
		"INSERT INTO transactions (debit_account_id, credit_account_id, amount, description) VALUES ($1, $2, $3, $4)",
		tx.DebitAccountID, tx.CreditAccountID, tx.Amount, tx.Description)

	if err != nil {
		return err
	}

	return dbTx.Commit()
}
