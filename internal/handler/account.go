package handler

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

// GetAccounts godoc
// @Summary Get all accounts
// @Description Get list of all bank accounts
// @Tags accounts
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Account
// @Failure 500 {object} map[string]string
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// Transfer godoc
// @Summary Transfer money between accounts
// @Description Make a money transfer from one account to another
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param transaction body domain.Transaction true "Transfer data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transfer [post]
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
