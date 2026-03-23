package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"test1/internal/models"
	"test1/internal/service"
)

type MemberHandler struct {
	Service *service.MemberService
}

func NewMemberHandler(service *service.MemberService) *MemberHandler {
	return &MemberHandler{Service: service}
}

func (h *MemberHandler) HandleMembers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetMembers(w, r)
	case http.MethodPost:
		h.CreateMember(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *MemberHandler) HandleMemberByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "members" {
		http.NotFound(w, r)
		return
	}

	memberID, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid member id")
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.DeleteMember(w, r, memberID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *MemberHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	members, err := h.Service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch members")
		return
	}

	writeJSON(w, http.StatusOK, members)
}

func (h *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMemberRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	member, err := h.Service.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, member)
}

func (h *MemberHandler) DeleteMember(w http.ResponseWriter, r *http.Request, memberID int) {
	err := h.Service.Delete(r.Context(), memberID)
	if err != nil {
		if err.Error() == "member not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "member deleted successfully",
	})
}

func (h *MemberHandler) GetMembersByBranch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 3 || parts[0] != "branches" || parts[2] != "members" {
		http.NotFound(w, r)
		return
	}

	branchID, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid branch id")
		return
	}

	members, err := h.Service.GetByBranch(r.Context(), branchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch members")
		return
	}

	writeJSON(w, http.StatusOK, members)
}
