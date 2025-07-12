package app

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/internal/handlers"
	"gluttony/internal/i18n"
	"gluttony/internal/recipe/bleve"
	"gluttony/internal/service/recipe"
	"gluttony/internal/service/user"
	"gluttony/internal/user/postgres"
	"gluttony/x/httpx"
	"gluttony/x/image"
	"gluttony/x/log"
	"gluttony/x/session"
	"gluttony/x/sqlx"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger

	recipeService *recipe.Service
	userService   *user.Service

	httpServer *http.Server
}

func New(cfg *config.Config, sec *config.Secret) (*App, error) {
	logger, err := NewLogger(cfg.Mode, cfg.Logger.Level, cfg.Logger.Path)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	if err := os.MkdirAll(cfg.WorkDir, 0750); err != nil {
		return nil, fmt.Errorf("create root working directory: %w", err)
	}

	db, err := sqlx.New(context.Background(), cfg.Database, sec.Database)
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	rootDir, err := os.OpenRoot(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("open root directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(cfg.WorkDir, "media"), 0750); err != nil {
		return nil, fmt.Errorf("create media directory: %w", err)
	}

	mediaDir, err := rootDir.OpenRoot("media")
	if err != nil {
		return nil, fmt.Errorf("open media directory: %w", err)
	}

	assetsFS, err := GetAssets(cfg.Mode)
	if err != nil {
		return nil, fmt.Errorf("get assets: %w", err)
	}

	sessionStore := session.NewStoreMemory()
	sessionService := session.NewService(sessionStore)
	userStore := postgres.NewStore(db)

	userService := user.NewService(userStore, sessionService)
	mediaService := image.NewService(mediaDir)
	recipeSearchIndex, err := bleve.New(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("create recipe search index: %w", err)
	}

	recipeService, err := recipe.NewService(db, mediaService, recipeSearchIndex, logger)
	if err != nil {
		return nil, fmt.Errorf("create recipe service: %w", err)
	}

	mux := httpx.NewRouter()
	i18nManager := i18n.NewI18n()
	middlewares := []httpx.Middleware{
		handlers.ErrorMiddleware(logger),
		handlers.I18nMiddleware(i18nManager),
		handlers.AuthenticationMiddleware(sessionService),
	}
	if cfg.Mode == config.ModeDev && cfg.Impersonate != "" {
		middlewares = append(middlewares, handlers.ImpersonateMiddleware(cfg.Impersonate, userService, sessionService))
	}

	mux.Use(middlewares...)
	MountRoutes(mux, cfg.Mode, assetsFS, mediaDir.FS())
	MountWebRoutes(mux, cfg, sessionService, userService, recipeService)

	const defaultTimeout = 15 * time.Second
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           mux,
		ReadTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       time.Minute,
	}

	return &App{
		recipeService: recipeService,
		userService:   userService,
		httpServer:    httpServer,
		cfg:           cfg,
		logger:        logger,
	}, nil
}

func NewLogger(mode config.Mode, level slog.Level, filePath string) (*slog.Logger, error) {
	if mode == config.ModeProd {
		logger, err := log.NewProd(level, filePath)
		if err != nil {
			return nil, fmt.Errorf("create prod logger: %w", err)
		}

		return logger, nil
	}

	return log.NewDev(level), nil
}
