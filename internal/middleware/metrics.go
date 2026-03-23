package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type MetricsStore struct {
	mu               sync.Mutex
	TotalRequests    int                `json:"total_requests"`
	RequestsByPath   map[string]int     `json:"requests_by_path"`
	RequestsByMethod map[string]int     `json:"requests_by_method"`
	TotalErrors      int                `json:"total_errors"`
	AverageLatencyMS float64            `json:"average_latency_ms"`
	LatencyByPathMS  map[string]float64 `json:"latency_by_path_ms"`
	totalLatency     time.Duration
	latencyByPath    map[string]time.Duration
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		RequestsByPath:   make(map[string]int),
		RequestsByMethod: make(map[string]int),
		LatencyByPathMS:  make(map[string]float64),
		latencyByPath:    make(map[string]time.Duration),
	}
}

func (m *MetricsStore) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &metricsStatusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(rec, r)
		latency := time.Since(start)

		m.mu.Lock()
		defer m.mu.Unlock()

		m.TotalRequests++
		m.RequestsByPath[r.URL.Path]++
		m.RequestsByMethod[r.Method]++
		m.totalLatency += latency
		m.latencyByPath[r.URL.Path] += latency

		if rec.statusCode >= 400 {
			m.TotalErrors++
		}

		m.AverageLatencyMS = float64(m.totalLatency.Milliseconds()) / float64(m.TotalRequests)

		for path, total := range m.latencyByPath {
			count := m.RequestsByPath[path]
			if count > 0 {
				m.LatencyByPathMS[path] = float64(total.Milliseconds()) / float64(count)
			}
		}
	})
}

func (m *MetricsStore) Handler(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(m)
}

type metricsStatusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *metricsStatusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
