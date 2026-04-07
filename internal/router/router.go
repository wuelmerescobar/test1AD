package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"test1/internal/config"
	"test1/internal/handler"
	"test1/internal/mailer"
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

	accountRepo := repository.NewAccountRepository(db)
	staffRepo := repository.NewStaffUserRepository(db)
	staffUserService := service.NewStaffUserService(staffRepo)
	staffUserHandler := handler.NewStaffUserHandler(staffUserService)

	appMailer := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.MailFrom)
	authService := service.NewAuthService(accountRepo, staffRepo, appMailer, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/metrics", metricsHandler.Metrics)

	mux.HandleFunc("/auth/register-staff", authHandler.RegisterStaff)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.Handle("/auth/me", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(authHandler.Me)))

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
		if len(r.URL.Path) >= 6 && r.URL.Path[len(r.URL.Path)-6:] == "/staff" {
			staffUserHandler.GetStaffByBranch(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/books", bookHandler.HandleBooks)

	protectedBooksByID := middleware.Auth(cfg.JWTSecret)(
		middleware.RequireRole("admin", "librarian")(http.HandlerFunc(bookHandler.HandleBookByID)),
	)
	mux.Handle("/books/", protectedBooksByID)

	protectedCopies := middleware.Auth(cfg.JWTSecret)(
		middleware.RequireRole("admin", "librarian")(http.HandlerFunc(bookCopyHandler.HandleBookCopies)),
	)
	mux.Handle("/book-copies", protectedCopies)

	mux.HandleFunc("/members", memberHandler.HandleMembers)

	protectedMembersByID := middleware.Auth(cfg.JWTSecret)(
		middleware.RequireRole("admin")(http.HandlerFunc(memberHandler.HandleMemberByID)),
	)
	mux.Handle("/members/", protectedMembersByID)

	var handlerChain http.Handler = mux
	handlerChain = middleware.RateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst)(handlerChain)
	handlerChain = middleware.CORS(cfg.AllowedOrigins)(handlerChain)
	handlerChain = middleware.Gzip(handlerChain)
	handlerChain = metricsStore.Middleware(handlerChain)
	handlerChain = middleware.Logging(logger)(handlerChain)
	handlerChain = middleware.Recover(logger)(handlerChain)

	return handlerChain
}
