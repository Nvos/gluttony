package app

import (
	"fmt"
	"github.com/spf13/afero"
	"gluttony/internal/config"
	"gluttony/internal/database"
	"gluttony/internal/livereload"
	"gluttony/internal/media"
	"gluttony/internal/recipe"
	"gluttony/internal/security"
	"gluttony/internal/user"
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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	rootFS := afero.NewBasePathFs(afero.NewOsFs(), cfg.WorkDirectoryPath)
	if err := os.MkdirAll(cfg.WorkDirectoryPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create root working directory: %w", err)
	}

	directories, err := NewDirectories(cfg.Mode, rootFS)
	if err != nil {
		return nil, fmt.Errorf("create directories: %w", err)
	}

	db, err := database.NewSqlite(cfg.WorkDirectoryPath)
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

	recipeService, err := recipe.NewService(db, mediaStore, recipeSearchIndex)
	if err != nil {
		return nil, fmt.Errorf("create recipe service: %w", err)
	}

	liveReload := livereload.New(logger)

	mux := http.NewServeMux()
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
