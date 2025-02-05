package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"gluttony/internal/database"
	"gluttony/internal/media"
	"gluttony/internal/recipe"
	"gluttony/internal/security"
	"gluttony/internal/share"
	"gluttony/internal/templating"
	"gluttony/internal/user"
	"gluttony/tools/reload"
	"gluttony/x/httpx"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := Run(groupCtx, group, logger); err != nil {
		logger.Error("failed to gracefully start gluttony", slog.String("error", err.Error()))
		return
	}

	if err := group.Wait(); err != nil {
		logger.Error("failed to gracefully shutdown goroutine", slog.String("error", err.Error()))
	}
}

func Run(ctx context.Context, group *errgroup.Group, logger *slog.Logger) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	// TODO: configurable
	workDir := filepath.Clean("/mnt/c/Users/HARDPC/Documents/gluttony-workdir")

	if err := os.MkdirAll(workDir, os.ModePerm); err != nil {
		return err
	}

	db, err := database.New(workDir)
	if err != nil {
		return fmt.Errorf("create db: %w", err)
	}

	assetsDir := os.DirFS(filepath.Join(wd, "assets"))
	workDirFS := afero.NewBasePathFs(afero.NewOsFs(), workDir)
	sessionStore := security.NewSessionStore()
	userService := user.NewService(db, sessionStore)
	mediaStore := media.NewStore(workDirFS)
	recipeService := recipe.NewService(db, mediaStore, workDir)

	reloader := reload.New()
	if err := reloader.Watch(ctx, reload.WatchConfig{
		Extensions: []string{".gohtml", ".html", ".css", ".js"},
		Directories: []string{
			filepath.Join("assets"),
			filepath.Clean(filepath.Join("internal/templating/templates")),
			filepath.Clean(filepath.Join("internal/user/templates")),
			filepath.Clean(filepath.Join("internal/recipe/templates")),
		},
	}); err != nil {
		return err
	}

	userTemplating := templating.New(
		os.DirFS(filepath.Join("internal/user/templates")),
	)
	recipeTemplating := templating.New(
		os.DirFS(filepath.Join("internal/recipe/templates")),
	)

	if err := userService.Create(ctx, "admin", "admin"); err != nil {
		return err
	}

	assetHandle := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.StripPrefix("/assets", http.FileServerFS(assetsDir)).ServeHTTP(w, r)
	}

	mediaHandle := func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/media", http.FileServerFS(os.DirFS(workDir))).ServeHTTP(w, r)
	}

	router := http.NewServeMux()
	middlewares := []httpx.MiddlewareFunc{
		security.NewAuthenticationMiddleware(sessionStore),
		share.ContextMiddleware,
		httpx.NewErrorMiddleware(logger),
	}

	router.HandleFunc("GET /reload", reloader.Handle)
	// Public
	router.HandleFunc("GET /assets/{pathname...}", assetHandle)
	router.HandleFunc("GET /media/{pathname...}", mediaHandle)

	userDeps := user.NewDeps(sessionStore, userTemplating, logger, userService)
	user.Routes(userDeps, router, middlewares...)

	recipeDeps := recipe.NewDeps(recipeService, logger, recipeTemplating, mediaStore)
	recipe.Routes(recipeDeps, router, middlewares...)

	httpServer := &http.Server{
		// TODO: cfg
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}

	group.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("start http server: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFn()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return nil
	})

	return nil
}
