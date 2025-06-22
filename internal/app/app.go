package app

import (
	"context"
	"fmt"
	"gluttony/internal/config"
	"gluttony/internal/handlers"
	"gluttony/internal/recipe/bleve"
	"gluttony/internal/service/recipe"
	"gluttony/internal/service/user"
	"gluttony/internal/user/postgres"
	"gluttony/migrations"
	"gluttony/pkg/database"
	"gluttony/pkg/log"
	"gluttony/pkg/media"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
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

func New(cfg *config.Config) (*App, error) {
	logger, err := NewLogger(cfg.Environment, cfg.Log.Level, cfg.Log.Path)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	if err := os.MkdirAll(cfg.WorkDirectoryPath, 0750); err != nil {
		return nil, fmt.Errorf("create root working directory: %w", err)
	}

	db, err := database.New(context.Background(), cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	if err := database.Migrate(db, migrations.Migrations); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	rootDir, err := os.OpenRoot(cfg.WorkDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("open root directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(cfg.WorkDirectoryPath, "media"), 0750); err != nil {
		return nil, fmt.Errorf("create media directory: %w", err)
	}

	mediaDir, err := rootDir.OpenRoot("media")
	if err != nil {
		return nil, fmt.Errorf("open media directory: %w", err)
	}

	assetsFS, err := GetAssets(cfg.Environment)
	if err != nil {
		return nil, fmt.Errorf("get assets: %w", err)
	}

	sessionStore := session.NewStoreMemory()
	sessionService := session.NewService(sessionStore)
	userStore := postgres.NewStore(db)

	userService := user.NewService(userStore, sessionService)
	mediaService := media.New(mediaDir)
	recipeSearchIndex, err := bleve.New(cfg.WorkDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("create recipe search index: %w", err)
	}

	recipeService, err := recipe.NewService(db, mediaService, recipeSearchIndex, logger)
	if err != nil {
		return nil, fmt.Errorf("create recipe service: %w", err)
	}

	mux := router.NewRouter()

	middlewares := []router.Middleware{
		handlers.ErrorMiddleware(logger),
		handlers.AuthenticationMiddleware(sessionService),
	}
	if cfg.Environment == config.EnvDevelopment && cfg.Impersonate != "" {
		middlewares = append(middlewares, handlers.ImpersonateMiddleware(cfg.Impersonate, userService, sessionService))
	}

	mux.Use(middlewares...)
	MountRoutes(mux, cfg.Environment, assetsFS, mediaDir.FS())
	MountWebRoutes(mux, sessionService, userService, recipeService)

	const defaultTimeout = 15 * time.Second
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
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

func NewLogger(mode config.Environment, level slog.Level, filePath string) (*slog.Logger, error) {
	if mode == config.EnvProduction {
		logger, err := log.NewProd(level, filePath)
		if err != nil {
			return nil, fmt.Errorf("create prod logger: %w", err)
		}

		return logger, nil
	}

	return log.NewDev(level), nil
}
