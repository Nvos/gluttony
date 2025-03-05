package app

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"gluttony/internal/config"
	"gluttony/internal/database"
	"gluttony/internal/html"
	"gluttony/internal/livereload"
	"gluttony/internal/log"
	"gluttony/internal/media"
	"gluttony/internal/recipe"
	"gluttony/internal/security"
	"gluttony/internal/user"
	"gluttony/internal/web"
	"gluttony/internal/web/templates"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	cfg    config.Config
	logger *slog.Logger

	recipeService *recipe.Service
	userService   *user.Service

	liveReload *livereload.LiveReload
	httpServer *http.Server
}

func NewApp(cfg config.Config) (*App, error) {
	logFile, err := os.Create(cfg.LogFilePath)
	if err != nil {
		return nil, fmt.Errorf("create log file: %w", err)
	}

	logger := log.New(cfg.Mode, cfg.LogLevel, logFile)

	rootFS := afero.NewBasePathFs(afero.NewOsFs(), cfg.WorkDirectoryPath)
	if err := os.MkdirAll(cfg.WorkDirectoryPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create root working directory: %w", err)
	}

	directories, err := NewDirectories(cfg.Mode, rootFS)
	if err != nil {
		return nil, fmt.Errorf("create directories: %w", err)
	}

	// TODO: move to config
	dbURL := "postgres://postgres:postgres@localhost:5432/gluttony?sslmode=disable"
	db, err := database.NewPostgres(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	sessionStore := security.NewSessionStore()
	userService := user.NewService(db, sessionStore)
	mediaStore := media.NewStore(directories.Media)
	recipeSearchIndex, err := recipe.NewSearchIndex(cfg.WorkDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("create recipe search index: %w", err)
	}

	recipeService, err := recipe.NewService(db, mediaStore, recipeSearchIndex, logger)
	if err != nil {
		return nil, fmt.Errorf("create recipe service: %w", err)
	}

	liveReload := livereload.New(logger)
	renderer, err := html.NewRenderer(templates.GetTemplates(cfg.Mode), html.RendererOptions{
		IsReloadEnabled: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create html renderer: %w", err)
	}
	mux := web.NewRouter(renderer)

	middlewares := []web.Middleware{
		web.AuthenticationMiddleware(sessionStore),
		web.ErrorMiddleware(logger),
	}

	mux.Use(middlewares...)
	MountRoutes(mux, cfg.Mode, liveReload, directories)
	MountWebRoutes(mux, logger, sessionStore, userService, recipeService)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.WebHost, cfg.WebPort),
		Handler: mux,
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
