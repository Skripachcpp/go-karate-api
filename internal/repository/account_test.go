package repository

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgresRepo_TransferMoney(t *testing.T) {
	// Подключаемся к тестовой БД
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	// Инициализируем тестовые данные
	repo.InitDB()
	db.Exec("DELETE FROM accounts") // Очищаем перед тестом
	db.Exec("INSERT INTO accounts (id, number, balance) VALUES (1, 'test1', 1000), (2, 'test2', 500)")

	t.Run("Successful transfer", func(t *testing.T) {
		err := repo.TransferMoney(1, 2, 100)
		if err != nil {
			t.Errorf("Transfer failed: %v", err)
		}

		// Проверяем балансы
		var balance1, balance2 float64
		db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&balance1)
		db.QueryRow("SELECT balance FROM accounts WHERE id = 2").Scan(&balance2)

		if balance1 != 900 {
			t.Errorf("Expected balance 900, got %f", balance1)
		}
		if balance2 != 600 {
			t.Errorf("Expected balance 600, got %f", balance2)
		}
	})

	// TODO: разобраться в чем ошибка
	t.Run("Insufficient funds", func(t *testing.T) {
		err := repo.TransferMoney(1, 2, 2000) // Слишком большая сумма
		if err == nil {
			t.Error("Expected error for insufficient funds, got nil")
		}
	})
}

func TestPostgresRepo_GetAccounts(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}``
	defer db.Close()

	repo := NewPostgresRepo(db)
	repo.InitDB()
	db.Exec("DELETE FROM accounts")
	db.Exec("INSERT INTO accounts (number, balance) VALUES ('test1', 1000)")

	accounts, err := repo.GetAccounts()
	if err != nil {
		t.Errorf("GetAccounts failed: %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(accounts))
	}
}
