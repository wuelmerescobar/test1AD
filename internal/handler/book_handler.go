package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"test1/internal/models"
	"test1/internal/service"
)

type BookHandler struct {
	Service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{Service: service}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetBooks(w, r)
	case http.MethodPost:
		h.CreateBook(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "books" {
		http.NotFound(w, r)
		return
	}

	bookID, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.UpdateBook(w, r, bookID)
	case http.MethodDelete:
		h.DeleteBook(w, r, bookID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch books")
		return
	}

	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	book, err := h.Service.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request, bookID int) {
	var req models.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	book, err := h.Service.Update(r.Context(), bookID, req)
	if err != nil {
		if err.Error() == "book not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request, bookID int) {
	err := h.Service.Delete(r.Context(), bookID)
	if err != nil {
		if err.Error() == "book not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "book deleted successfully",
	})
}

func (h *BookHandler) GetBooksByBranch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 3 || parts[0] != "branches" || parts[2] != "books" {
		http.NotFound(w, r)
		return
	}

	branchID, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid branch id")
		return
	}

	books, err := h.Service.GetAvailableByBranch(r.Context(), branchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch books")
		return
	}

	writeJSON(w, http.StatusOK, books)
}
