package main

import (
    "accounting-core/internal/handler"
    "accounting-core/internal/repository"
    "database/sql"
    "log"
    "net/http"

    _ "github.com/lib/pq"
)

func main() {
    // Подключение к БД
    connStr := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Cannot connect to DB:", err)
    }

    // Инициализация репозитория и базы
    repo := repository.NewPostgresRepo(db)
    if err := repo.InitDB(); err != nil {
        log.Fatal("Failed to init DB:", err)
    }
    if err := repo.SeedData(); err != nil {
        log.Fatal("Failed to seed data:", err)
    }

    // Инициализация хендлеров
    accountHandler := handler.NewAccountHandler(repo)

    // Маршруты
    http.HandleFunc("/accounts", accountHandler.GetAccounts)
    http.HandleFunc("/transfer", accountHandler.Transfer)

    log.Println("Server starting on :8080...")
    log.Println("API endpoints:")
    log.Println("  GET  http://localhost:8080/accounts")
    log.Println("  POST http://localhost:8080/transfer")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
