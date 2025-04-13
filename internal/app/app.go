package app

import (
	"context"
	"fmt"
	"gluttony/internal/handlers"
	"gluttony/internal/recipe/bleve"
	"gluttony/internal/service/recipe"
	"gluttony/internal/service/user"
	"gluttony/internal/user/postgres"
	"gluttony/migrations"
	"gluttony/pkg/database"
	"gluttony/pkg/html"
	"gluttony/pkg/livereload"
	"gluttony/pkg/log"
	"gluttony/pkg/media"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
	"gluttony/web/templates"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	cfg    Config
	logger *slog.Logger

	recipeService *recipe.Service
	userService   *user.Service

	liveReload *livereload.LiveReload
	httpServer *http.Server
}

func New(cfg Config) (*App, error) {
	logger, err := NewLogger(cfg.Mode, cfg.LogLevel, cfg.LogFilePath)
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

	assetsFS, err := GetAssets(cfg.Mode)
	if err != nil {
		return nil, fmt.Errorf("get assets: %w", err)
	}

	sessionStore := session.NewStoreMemory()
	sessionService := session.NewService(sessionStore)
	userStore := postgres.NewStore(db)

	userService := user.NewService(userStore, sessionService)
	mediaStore := media.NewStore(mediaDir)
	recipeSearchIndex, err := bleve.New(cfg.WorkDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("create recipe search index: %w", err)
	}

	recipeService, err := recipe.NewService(db, mediaStore, recipeSearchIndex, logger)
	if err != nil {
		return nil, fmt.Errorf("create recipe service: %w", err)
	}

	var liveReload *livereload.LiveReload
	if cfg.Mode == Dev {
		liveReload = livereload.New(logger)
	}

	renderer, err := html.NewRenderer(GetTemplates(cfg.Mode), html.RendererOptions{
		IsReloadEnabled: cfg.Mode == Dev,
	})
	if err != nil {
		return nil, fmt.Errorf("create html renderer: %w", err)
	}
	mux := router.NewRouter(renderer)

	middlewares := []router.Middleware{
		handlers.AuthenticationMiddleware(sessionService),
		handlers.ErrorMiddleware(logger),
	}

	mux.Use(middlewares...)
	MountRoutes(mux, cfg.Mode, liveReload, assetsFS, mediaDir.FS())
	MountWebRoutes(mux, sessionService, userService, recipeService)

	const defaultTimeout = 15 * time.Second
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.WebHost, cfg.WebPort),
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
		liveReload:    liveReload,
		logger:        logger,
	}, nil
}

func GetTemplates(mode Mode) fs.FS {
	if mode == Prod {
		return templates.Embedded
	}

	return os.DirFS("web/templates")
}

func NewLogger(mode Mode, level slog.Level, filePath string) (*slog.Logger, error) {
	if mode == Prod {
		logger, err := log.NewProd(level, filePath)
		if err != nil {
			return nil, fmt.Errorf("create prod logger: %w", err)
		}

		return logger, nil
	}

	return log.NewDev(level), nil
}
