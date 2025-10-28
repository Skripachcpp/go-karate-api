package handler

import (
	"accounting-core/internal/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock репозитория для тестов
type MockRepo struct{}

func (m *MockRepo) InitDB() error   { return nil }
func (m *MockRepo) SeedData() error { return nil }

func (m *MockRepo) GetAccounts() ([]domain.Account, error) {
	return []domain.Account{
		{ID: 1, Number: "1001", Balance: 1000},
		{ID: 2, Number: "1002", Balance: 500},
	}, nil
}

func (m *MockRepo) TransferMoney(from, to int64, amount float64) error {
	if amount <= 0 {
		return &domain.BusinessError{Message: "Amount must be positive"}
	}
	if from == to {
		return &domain.BusinessError{Message: "Cannot transfer to same account"}
	}
	return nil
}

func (m *MockRepo) GetTransactions() ([]domain.Transaction, error) {
	return []domain.Transaction{
		{ID: 1, DebitAccountID: 1, CreditAccountID: 2, Amount: 100},
	}, nil
}

func TestAccountHandler_GetAccounts(t *testing.T) {
	handler := NewAccountHandler(&MockRepo{})

	req := httptest.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()

	handler.GetAccounts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var accounts []domain.Account
	json.NewDecoder(w.Body).Decode(&accounts)

	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}
}

func TestAccountHandler_Transfer(t *testing.T) {
	handler := NewAccountHandler(&MockRepo{})

	t.Run("Valid transfer", func(t *testing.T) {
		transfer := domain.Transaction{
			DebitAccountID:  1,
			CreditAccountID: 2,
			Amount:          100,
		}
		body, _ := json.Marshal(transfer)

		req := httptest.NewRequest("POST", "/transaction", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Transfer(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Invalid amount", func(t *testing.T) {
		transfer := domain.Transaction{Amount: -100}
		body, _ := json.Marshal(transfer)

		req := httptest.NewRequest("POST", "/transaction", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Transfer(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}
