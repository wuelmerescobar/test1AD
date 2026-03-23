package handler

import (
	"net/http"

	"test1/internal/middleware"
)

type MetricsHandler struct {
	Store *middleware.MetricsStore
}

func NewMetricsHandler(store *middleware.MetricsStore) *MetricsHandler {
	return &MetricsHandler{Store: store}
}

func (h *MetricsHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	h.Store.Handler(w, r)
}
