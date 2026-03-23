package handler

import (
	"encoding/json"
	"net/http"

	"test1/internal/models"
	"test1/internal/service"
)

type BranchHandler struct {
	Service *service.BranchService
}

func NewBranchHandler(service *service.BranchService) *BranchHandler {
	return &BranchHandler{Service: service}
}

func (h *BranchHandler) HandleBranches(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetBranches(w, r)
	case http.MethodPost:
		h.CreateBranch(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *BranchHandler) GetBranches(w http.ResponseWriter, r *http.Request) {
	branches, err := h.Service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch branches")
		return
	}

	writeJSON(w, http.StatusOK, branches)
}

func (h *BranchHandler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBranchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	branch, err := h.Service.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, branch)
}
