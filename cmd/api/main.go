package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theretech/retech-core-api/internal/application/usecase"
	"github.com/theretech/retech-core-api/internal/infrastructure/config"
	applogger "github.com/theretech/retech-core-api/internal/infrastructure/logger"
	"github.com/theretech/retech-core-api/internal/infrastructure/persistence"
	httptransport "github.com/theretech/retech-core-api/internal/interfaces/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", slog.Any("error", err))
		os.Exit(1)
	}

	log := applogger.New(cfg.LogLevel)
	slog.SetDefault(log)

	db, err := persistence.Open(cfg.DatabaseURL)
	if err != nil {
		log.Error("database open failed", slog.Any("error", err))
		os.Exit(1)
	}
	if err := persistence.AutoMigrate(db); err != nil {
		log.Error("auto migrate failed", slog.Any("error", err))
		os.Exit(1)
	}

	tenantRepo := persistence.NewTenantRepository(db)
	appRepo := persistence.NewApplicationRepository(db)
	taRepo := persistence.NewTenantApplicationRepository(db)
	utRepo := persistence.NewUserTenantRepository(db)

	tenantUC := usecase.NewTenant(tenantRepo)
	appUC := usecase.NewApplication(appRepo, taRepo)
	relUC := usecase.NewRelationship(tenantRepo, appRepo, taRepo, utRepo)

	router := httptransport.NewRouter(httptransport.RouterDeps{
		DB:             db,
		Log:            log,
		TenantRepo:     tenantRepo,
		TenantUC:       tenantUC,
		ApplicationUC:  appUC,
		RelationshipUC: relUC,
		GinMode:        cfg.GinMode,
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("http server listening", slog.String("addr", cfg.HTTPAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", slog.Any("error", err))
	}
	log.Info("server stopped")
}
