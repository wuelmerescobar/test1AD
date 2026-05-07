package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"test1/internal/models"
	"test1/internal/service"
)

type BookCopyHandler struct {
	Service *service.BookCopyService
}

func NewBookCopyHandler(service *service.BookCopyService) *BookCopyHandler {
	return &BookCopyHandler{Service: service}
}

func (h *BookCopyHandler) HandleBookCopies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateBookCopy(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *BookCopyHandler) CreateBookCopy(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBookCopyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	bookCopy, err := h.Service.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, bookCopy)
}

func (h *BookCopyHandler) GetCopiesByBook(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "books" || parts[2] != "copies" {
		http.NotFound(w, r)
		return
	}

	bookID, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	copies, err := h.Service.GetByBook(r.Context(), bookID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, copies)
}
