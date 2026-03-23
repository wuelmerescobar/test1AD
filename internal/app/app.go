package app

import (
	"database/sql"
	"log/slog"
	"os"

	"test1/internal/config"
)

type App struct {
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func New(cfg *config.Config, db *sql.DB) *App {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &App{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}
