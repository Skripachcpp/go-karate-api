package main

import (
	"accounting-core/internal/handler"
	"accounting-core/internal/repository"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Accounting Core API
// @version 1.0

// @description Бухгалтерское ядро банка на Go
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
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

	// Swagger UI
	// Добавь этот код ПЕРЕД строкой http.HandleFunc("/swagger/", ...)

	// Прямой endpoint для swagger.json
	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Остальной код остается как есть
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server starting on :8080...")
	log.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
