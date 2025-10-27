package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Подключаемся к БД
	connStr := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицу если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			number VARCHAR(50) UNIQUE NOT NULL,
			balance DECIMAL(15,2) DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Добавляем тестовый счет если нет счетов
	var count int
	db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO accounts (number, balance) VALUES ('1001', 1000.00)")
		db.Exec("INSERT INTO accounts (number, balance) VALUES ('1002', 500.00)")
	}

	fmt.Println("База данных готова!")

	// Запускаем HTTP сервер
	http.HandleFunc("/accounts", getAccounts)
	http.HandleFunc("/transfer", transferMoney)

	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, number, balance FROM accounts")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var accounts []map[string]interface{}
	for rows.Next() {
		var id int
		var number string
		var balance float64
		rows.Scan(&id, &number, &balance)
		accounts = append(accounts, map[string]interface{}{
			"id":      id,
			"number":  number,
			"balance": balance,
		})
	}

	json.NewEncoder(w).Encode(accounts)
}

func transferMoney(w http.ResponseWriter, r *http.Request) {
	var data struct {
		From   int     `json:"from"`
		To     int     `json:"to"`
		Amount float64 `json:"amount"`
	}

	json.NewDecoder(r.Body).Decode(&data)

	// Простая проверка
	if data.Amount <= 0 {
		http.Error(w, "Сумма должна быть положительной", 400)
		return
	}

	// Выполняем перевод в транзакции
	tx, _ := db.Begin()

	// Списание
	_, err1 := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", data.Amount, data.From)
	// Зачисление
	_, err2 := tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", data.Amount, data.To)

	if err1 != nil || err2 != nil {
		tx.Rollback()
		http.Error(w, "Ошибка перевода", 500)
		return
	}

	tx.Commit()
	fmt.Fprintf(w, "Перевод выполнен успешно!")
}
