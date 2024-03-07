package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"gluttony/internal/database"
	"gluttony/internal/recipe"
	"gluttony/internal/util/filepathutil"
	"gluttony/internal/util/httputil"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	root := filepath.Join(xdg.DataHome, "gluttony")
	if !filepathutil.IsFileExist(root) {
		if err := os.Mkdir(root, os.ModePerm); err != nil {
			panic(err)
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	logger.Info("Starting gluttony", slog.String("workdir", root))

	// TODO(AK) 05/03/2024: add config (koanf)
	databaseCfg := database.Config{
		Host:     "localhost",
		Database: "dev",
		Password: "developer",
		Username: "developer",
	}

	if err := database.Migrate(databaseCfg); err != nil {
		panic(err)
	}

	pool, err := database.ConnectPostgres(ctx, databaseCfg)
	if err != nil {
		panic(err)
	}

	recipeStore, err := recipe.NewPostgresStore(pool)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	if err := mountRecipeHandler(mux, recipeStore); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    ":6001",
		Handler: h2c.NewHandler(httputil.AllowAllCORSMiddleware(mux), &http2.Server{}),
		// TODO(AK) 05/03/2024: timeouts
	}

	go func() {
		wg.Add(1)
		defer func() {
			wg.Done()
		}()

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Http server", slog.String("err", err.Error()))
		}
	}()

	slog.Info("Gluttony started")
	<-shutdownChan
	// Begin shutdown
	cancel()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Warn("Graceful shut down of http server", slog.String("err", err.Error()))
	}

	wg.Wait()
}

func mountRecipeHandler(mux *http.ServeMux, recipeStore recipe.Store) error {
	path, handler, err := recipe.NewConnectHandler(recipeStore)
	if err != nil {
		return fmt.Errorf("mount recipe connect handler: %w", err)
	}

	mux.Handle(path, handler)

	return nil
}
