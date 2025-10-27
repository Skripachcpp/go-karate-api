package handler

// # HTTP-обработчики (API слой)

import (
	"accounting-core/internal/domain"
	"accounting-core/internal/repository"
	"encoding/json"
	"net/http"
)

type AccountHandler struct {
	repo *repository.PostgresRepo
}

func NewAccountHandler(repo *repository.PostgresRepo) *AccountHandler {
	return &AccountHandler{repo: repo}
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	if err := h.repo.TransferMoney(req.DebitAccountID, req.CreditAccountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer successful"})
}
