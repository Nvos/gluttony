package app

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
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

	rootFS := afero.NewBasePathFs(afero.NewOsFs(), cfg.WorkDirectoryPath)
	if err := os.MkdirAll(cfg.WorkDirectoryPath, 0750); err != nil {
		return nil, fmt.Errorf("create root working directory: %w", err)
	}

	directories, err := NewDirectories(cfg.Mode, rootFS)
	if err != nil {
		return nil, fmt.Errorf("create directories: %w", err)
	}

	// TODO: move to config
	dbCfg := database.Config{
		Name:     "gluttony",
		User:     "postgres",
		Host:     "localhost",
		Port:     "5432",
		Password: "postgres",
	}

	db, err := database.New(context.Background(), dbCfg)
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	if err := database.Migrate(db, migrations.Migrations); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	sessionStore := session.NewStoreMemory()
	sessionService := session.NewService(sessionStore)
	userStore := postgres.NewStore(db)

	userService := user.NewService(userStore, sessionService)
	mediaStore := media.NewStore(directories.Media)
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
	MountRoutes(mux, cfg.Mode, liveReload, directories)
	MountWebRoutes(mux, sessionService, userService, recipeService)

	const (
		defaultHeaderTimeout = time.Second * 10
	)
	httpServer := &http.Server{
		ReadHeaderTimeout: defaultHeaderTimeout,
		Addr:              fmt.Sprintf("%s:%d", cfg.WebHost, cfg.WebPort),
		Handler:           mux,
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
