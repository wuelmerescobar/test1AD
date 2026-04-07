package handler

import (
	"net/http"
	"strconv"
	"strings"

	"test1/internal/service"
)

type StaffUserHandler struct {
	Service *service.StaffUserService
}

func NewStaffUserHandler(service *service.StaffUserService) *StaffUserHandler {
	return &StaffUserHandler{Service: service}
}

func (h *StaffUserHandler) GetStaffByBranch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 3 || parts[0] != "branches" || parts[2] != "staff" {
		http.NotFound(w, r)
		return
	}

	branchID, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid branch id")
		return
	}

	staff, err := h.Service.GetByBranch(r.Context(), branchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch staff members")
		return
	}

	writeJSON(w, http.StatusOK, staff)
}
