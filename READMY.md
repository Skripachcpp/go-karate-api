# Подключиться к БД и выполнить файлы
psql -h localhost -U postgres -d accounting -f migrations/001_create_tables.sql

# Посмотреть все таблицы в базе
docker exec -it accounting-core-postgres-1 psql -U postgres -d postgres -c "\dt"

# Подключиться к PostgreSQL интерактивно
docker exec -it accounting-core-postgres-1 psql -U postgres -d postgres

# пример запроса
/transfer
{
  "amount": 100,
  "credit_account_id": 2,
  "debit_account_id": 1
}

# Запуск всех тестов
go test ./...

# Запуск тестов с verbose
go test -v ./...

# Запуск тестов конкретного пакета
go test -v ./internal/repository
go test -v ./internal/handler

# Запуск с покрытием
go test -cover ./...

# запустить проект
docker-compose up -d 
go run cmd/server/main.go
