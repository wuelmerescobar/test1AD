package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"test1/internal/config"
	"test1/internal/handler"
	"test1/internal/middleware"
	"test1/internal/repository"
	"test1/internal/service"
)

func NewRouter(db *sql.DB, cfg *config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	metricsStore := middleware.NewMetricsStore()

	healthHandler := handler.NewHealthHandler()
	metricsHandler := handler.NewMetricsHandler(metricsStore)

	branchRepo := repository.NewBranchRepository(db)
	branchService := service.NewBranchService(branchRepo)
	branchHandler := handler.NewBranchHandler(branchService)

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	bookCopyRepo := repository.NewBookCopyRepository(db)
	bookCopyService := service.NewBookCopyService(bookCopyRepo)
	bookCopyHandler := handler.NewBookCopyHandler(bookCopyService)

	memberRepo := repository.NewMemberRepository(db)
	memberService := service.NewMemberService(memberRepo)
	memberHandler := handler.NewMemberHandler(memberService)

	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/metrics", metricsHandler.Metrics)

	mux.HandleFunc("/branches", branchHandler.HandleBranches)
	mux.HandleFunc("/branches/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 6 && r.URL.Path[len(r.URL.Path)-6:] == "/books" {
			bookHandler.GetBooksByBranch(w, r)
			return
		}
		if len(r.URL.Path) >= 8 && r.URL.Path[len(r.URL.Path)-8:] == "/members" {
			memberHandler.GetMembersByBranch(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/books", bookHandler.HandleBooks)
	mux.HandleFunc("/books/", bookHandler.HandleBookByID)
	mux.HandleFunc("/book-copies", bookCopyHandler.HandleBookCopies)

	mux.HandleFunc("/members", memberHandler.HandleMembers)
	mux.HandleFunc("/members/", memberHandler.HandleMemberByID)

	var handlerChain http.Handler = mux
	handlerChain = middleware.RateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst)(handlerChain)
	handlerChain = middleware.CORS(cfg.AllowedOrigins)(handlerChain)
	handlerChain = middleware.Gzip(handlerChain)
	handlerChain = metricsStore.Middleware(handlerChain)
	handlerChain = middleware.Logging(logger)(handlerChain)
	handlerChain = middleware.Recover(logger)(handlerChain)

	return handlerChain
}
