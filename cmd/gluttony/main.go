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
	"gluttony/internal/i18n"
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe"
	"gluttony/internal/static"
	"gluttony/internal/x/connectx"
	"gluttony/internal/x/filepathx"
	"gluttony/internal/x/httpx"
	"gluttony/internal/x/loggerx"
	"gluttony/migrations"
	"gluttony/seeds"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"io/fs"
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

	logger, err := loggerx.NewLogger(slog.LevelDebug)
	if err != nil {
		panic(fmt.Sprintf("create logger: %v", err))
	}

	if err := Run(groupCtx, group, logger); err != nil {
		logger.Error("failed to gracefully start gluttony", slog.String("error", err.Error()))

		os.Exit(1)
	}

	if err := group.Wait(); err != nil {
		logger.Error("failed to gracefully shutdown goroutine", slog.String("error", err.Error()))

		os.Exit(1)
	}
}

func Run(ctx context.Context, group *errgroup.Group, logger *slog.Logger) error {
	workDirectories, err := devDirectories()
	if err != nil {
		return fmt.Errorf("initialize work directories: %w", err)
	}

	logger.Info(
		"Starting gluttony",
		slog.String("dataDir", workDirectories.DataDir),
		slog.String("configDir", workDirectories.ConfigDir),
	)

	cfg, err := config.LoadConfig(workDirectories.ConfigFS)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dbCfg := sqldb.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		Options:  cfg.Database.Options,
	}
	pool, err := sqldb.ConnectPostgres(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("create postgres connection: %w", err)
	}

	conn := stdlib.OpenDBFromPool(pool)
	isDbRunning, err := sqldb.IsDBRunning(ctx, conn)
	if err != nil {
		return fmt.Errorf("check if db is running: %w", err)
	}

	if !isDbRunning {
		return fmt.Errorf("db is not running")
	}

	if err := sqldb.Migrate(ctx, conn, migrations.FS); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if err := sqldb.Seed(ctx, conn, seeds.FS); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("close migrate db conn: %w", err)
	}

	recipeStore := recipe.NewStorePostgres(pool)

	sessionStore := auth.NewMemoryStorage()
	sessionManager := auth.NewSessionManager(sessionStore)
	beginner := transaction.NewPgxBeginner(pool)

	userStore, err := auth.NewUserPostgresStore(pool)
	if err != nil {
		return fmt.Errorf("create user postgres store: %w", err)
	}

	userService, err := auth.NewService(userStore, sessionManager)
	if err != nil {
		return fmt.Errorf("create user service: %w", err)
	}

	recipeService := recipe.NewService(beginner, recipeStore)

	ingredientStore := ingredient.NewStorePostgres(pool)

	ingredientService := ingredient.NewService(ingredientStore)

	connectInterceptors := connect.WithInterceptors(connectx.ErrorInterceptor(logger))

	mux := http.NewServeMux()
	if err := mountRecipeHandler(mux, recipeService, connectInterceptors); err != nil {
		return fmt.Errorf("mount recipe http handlers: %w", err)
	}

	if err := mountUserHandler(mux, userService, connectInterceptors); err != nil {
		return fmt.Errorf("mount user http handlers: %w", err)
	}

	if err := mountIngredientHandler(mux, ingredientService, connectInterceptors); err != nil {
		return fmt.Errorf("mount ingredient http handlers: %w", err)
	}

	mux.Handle("/static/", static.FileServeHandler("/static/", workDirectories.DataFS))
	mux.Handle("/storage/", static.UploadHandler(workDirectories.DataDir))

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: h2c.NewHandler(httpx.ComposeMiddlewares(
			mux,
			httpx.AllowAllCORSMiddleware,
			i18n.LocaleInjectionMiddleware(),
			auth.SessionHttpMiddleware(sessionManager),
		), &http2.Server{}),
		// TODO(AK) 05/03/2024: timeouts
	}

	group.Go(func() error {
		logger.Info("Http server listening on", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		logger.Info("Graceful shutdown started")

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownRelease()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}

		pool.Close()

		logger.Info("Graceful shutdown completed")
		return nil
	})

	logger.Info("Gluttony started")
	return nil
}

func mountRecipeHandler(mux *http.ServeMux, recipeStore *recipe.Service, opts ...connect.HandlerOption) error {
	path, handler, err := recipe.NewConnectHandler(recipeStore, opts...)
	if err != nil {
		return fmt.Errorf("mount recipe connect handler: %w", err)
	}

	mux.Handle(path, handler)

	return nil
}

func mountUserHandler(mux *http.ServeMux, service *auth.Service, opts ...connect.HandlerOption) error {
	path, handler, err := auth.NewConnectHandler(service, opts...)
	if err != nil {
		return fmt.Errorf("mount user connect handler: %w", err)
	}

	mux.Handle(path, handler)

	return nil
}

func mountIngredientHandler(mux *http.ServeMux, service *ingredient.Service, opts ...connect.HandlerOption) error {
	path, handler, err := ingredient.NewConnectHandler(service, opts...)
	if err != nil {
		return fmt.Errorf("mount ingredient connect handler: %w", err)
	}

	mux.Handle(path, handler)

	return nil
}

type WorkDirectories struct {
	DataDir string
	DataFS  fs.FS

	ConfigDir string
	ConfigFS  fs.FS
}

func prodDirectories() (WorkDirectories, error) {
	dataDir := filepath.Join(xdg.DataHome, "gluttony")
	if !filepathx.IsFileExist(dataDir) {
		if err := os.Mkdir(dataDir, os.ModePerm); err != nil {
			return WorkDirectories{}, fmt.Errorf("create data directory: %w", err)
		}
	}

	configDir := filepath.Join(xdg.ConfigHome, "gluttony")
	if !filepathx.IsFileExist(configDir) {
		if err := os.Mkdir(configDir, os.ModePerm); err != nil {
			return WorkDirectories{}, fmt.Errorf("create config directory: %w", err)
		}
	}

	return WorkDirectories{
		DataDir:   dataDir,
		DataFS:    os.DirFS(dataDir),
		ConfigDir: configDir,
		ConfigFS:  os.DirFS(configDir),
	}, nil
}

func devDirectories() (WorkDirectories, error) {
	wd, err := os.Getwd()
	if err != nil {
		return WorkDirectories{}, err
	}

	dataDir := filepath.Join(wd, "workdir")
	if !filepathx.IsFileExist(dataDir) {
		if err := os.Mkdir(dataDir, os.ModePerm); err != nil {
			return WorkDirectories{}, fmt.Errorf("create data directory: %w", err)
		}
	}

	return WorkDirectories{
		DataDir:   dataDir,
		DataFS:    os.DirFS(dataDir),
		ConfigDir: wd,
		ConfigFS:  os.DirFS(wd),
	}, nil
}
