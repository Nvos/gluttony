package main

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/jackc/pgx/v5/stdlib"
	"gluttony/internal/auth"
	"gluttony/internal/config"
	"gluttony/internal/database/sqldb"
	"gluttony/internal/database/transaction"
	"gluttony/internal/logger"
	"gluttony/internal/recipe"
	"gluttony/internal/user"
	"gluttony/internal/util/connectutil"
	"gluttony/internal/util/filepathutil"
	"gluttony/internal/util/httputil"
	"gluttony/migrations"
	"gluttony/seeds"
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
	wg := &sync.WaitGroup{}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	cleanup, err := Main(ctx, wg)
	if err != nil {
		panic(fmt.Sprintf("failed to start gluttony: %v", err))
	}

	<-shutdownChan
	// Begin shutdown
	cancel()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := cleanup(shutdownCtx); err != nil {
		panic(fmt.Sprintf("failed to gracefully shutdown gluttony: %v", err))
	}

	wg.Wait()
}

func Main(ctx context.Context, wg *sync.WaitGroup) (func(ctx context.Context) error, error) {
	dataDir, configDir, err := initializeUsedDirectories()
	if err != nil {
		return nil, fmt.Errorf("initialize app directories: %w", err)
	}

	log, _, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	log.Info("Starting gluttony", slog.String("dataDir", dataDir), slog.String("configDir", configDir))
	cfg, err := config.LoadConfig(configDir)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	pool, err := sqldb.ConnectPostgres(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("create postgres connection: %w", err)
	}

	conn := stdlib.OpenDBFromPool(pool)
	if err := sqldb.Migrate(ctx, conn, migrations.FS); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	if err := sqldb.Seed(ctx, conn, seeds.FS); err != nil {
		return nil, fmt.Errorf("seed database: %w", err)
	}

	if err := conn.Close(); err != nil {
		return nil, fmt.Errorf("close migrate db conn: %w", err)
	}

	recipeStore, err := recipe.NewPostgresStore(pool)
	if err != nil {
		return nil, fmt.Errorf("create recipe postgres store: %w", err)
	}

	sessionStore := auth.NewMemoryStorage[user.Session]()
	sessionManager := auth.NewSessionManager[user.Session](sessionStore)
	beginner := transaction.NewPgxBeginner(pool)

	userStore, err := user.NewPostgresStore(pool)
	if err != nil {
		return nil, fmt.Errorf("create user postgres store: %w", err)
	}

	userService, err := user.NewService(userStore, sessionManager)
	if err != nil {
		return nil, fmt.Errorf("create user service: %w", err)
	}

	recipeService, err := recipe.NewService(beginner, recipeStore)
	if err != nil {
		return nil, fmt.Errorf("create recipe service: %w", err)
	}

	connectInterceptors := connect.WithInterceptors(connectutil.ErrorInterceptor(log))

	mux := http.NewServeMux()
	if err := mountRecipeHandler(mux, recipeService, connectInterceptors); err != nil {
		return nil, fmt.Errorf("mount recipe http handlers: %w", err)
	}

	if err := mountUserHandler(mux, userService, connectInterceptors); err != nil {
		return nil, fmt.Errorf("mount user http handlers: %w", err)
	}

	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", cfg.Server.Port),
		Handler: h2c.NewHandler(httputil.ComposeMiddlewares(
			mux,
			httputil.AllowAllCORSMiddleware,
			auth.SessionHttpMiddleware(sessionManager),
		), &http2.Server{}),
		// TODO(AK) 05/03/2024: timeouts
	}

	go func() {
		wg.Add(1)
		defer func() {
			wg.Done()
		}()

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Http server", slog.String("err", err.Error()))
		}
	}()

	log.Info("Gluttony started")
	return func(timeoutCtx context.Context) error {
		var errs []error
		if err := server.Shutdown(timeoutCtx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown http server: %w", err))
		}

		return errors.Join(errs...)
	}, nil
}

func mountRecipeHandler(mux *http.ServeMux, recipeStore *recipe.Service, opts ...connect.HandlerOption) error {
	path, handler, err := recipe.NewConnectHandler(recipeStore, opts...)
	if err != nil {
		return fmt.Errorf("mount recipe connect handler: %w", err)
	}

	mux.Handle(path, handler)

	return nil
}

func mountUserHandler(mux *http.ServeMux, service *user.Service, opts ...connect.HandlerOption) error {
	path, handler, err := user.NewConnectHandler(service, opts...)
	if err != nil {
		return fmt.Errorf("mount user connect handler: %w", err)
	}

	mux.Handle(path, handler)

	return nil
}

func initializeUsedDirectories() (string, string, error) {
	dataDir := filepath.Join(xdg.DataHome, "gluttony")
	if !filepathutil.IsFileExist(dataDir) {
		if err := os.Mkdir(dataDir, os.ModePerm); err != nil {
			return "", "", fmt.Errorf("create data directory: %w", err)
		}
	}

	configDir := filepath.Join(xdg.ConfigHome, "gluttony")
	if !filepathutil.IsFileExist(configDir) {
		if err := os.Mkdir(configDir, os.ModePerm); err != nil {
			return "", "", fmt.Errorf("create config directory: %w", err)
		}
	}

	return dataDir, configDir, nil
}
