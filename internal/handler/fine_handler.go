package handler

import (
	"net/http"

	"test1/internal/service"
)

type FineHandler struct {
	Service *service.FineService
}

func NewFineHandler(service *service.FineService) *FineHandler {
	return &FineHandler{Service: service}
}

func (h *FineHandler) HandleFines(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetFines(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *FineHandler) GetFines(w http.ResponseWriter, r *http.Request) {
	fines, err := h.Service.GetAllReports(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch fines")
		return
	}

	writeJSON(w, http.StatusOK, fines)
}
