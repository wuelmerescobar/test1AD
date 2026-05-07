package handler

import (
	"net/http"

	"test1/internal/service"
)

type LoanHandler struct {
	Service *service.LoanService
}

func NewLoanHandler(service *service.LoanService) *LoanHandler {
	return &LoanHandler{Service: service}
}

func (h *LoanHandler) HandleLoans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetLoans(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *LoanHandler) GetLoans(w http.ResponseWriter, r *http.Request) {
	loans, err := h.Service.GetAllReports(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch loans")
		return
	}

	writeJSON(w, http.StatusOK, loans)
}
