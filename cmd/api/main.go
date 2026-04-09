package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	stdlog "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/theretech/retech-core-api/internal/application/usecase"
	"github.com/theretech/retech-core-api/internal/infrastructure/config"
	applogger "github.com/theretech/retech-core-api/internal/infrastructure/logger"
	"github.com/theretech/retech-core-api/internal/infrastructure/persistence"
	httptransport "github.com/theretech/retech-core-api/internal/interfaces/http"
	"github.com/theretech/retech-core-api/internal/version"
)

// rewriteDoubleDashConfig converte "--config arquivo" em "-config arquivo" para o flag padrão do Go.
func rewriteDoubleDashConfig() {
	if len(os.Args) < 2 {
		return
	}
	out := make([]string, 0, len(os.Args))
	out = append(out, os.Args[0])
	for i := 1; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == "--config" && i+1 < len(os.Args):
			out = append(out, "-config", os.Args[i+1])
			i++
		case strings.HasPrefix(os.Args[i], "--config="):
			out = append(out, "-config="+strings.TrimPrefix(os.Args[i], "--config="))
		default:
			out = append(out, os.Args[i])
		}
	}
	os.Args = out
}

// httpPublicBase monta a URL base para logs (ex.: ":8080" → http://0.0.0.0:8080).
func httpPublicBase(addr string) string {
	addr = strings.TrimSpace(addr)
	switch {
	case addr == "":
		return "http://localhost:8080"
	case strings.HasPrefix(addr, ":"):
		return "http://0.0.0.0" + addr
	default:
		return "http://" + addr
	}
}

// listenDisplay retorna (rótulo da porta, endereço de escuta) para logs no estilo do auth-api.
func listenDisplay(addr string) (portLabel, listenAddr string) {
	addr = strings.TrimSpace(addr)
	if strings.HasPrefix(addr, ":") {
		return strings.TrimPrefix(addr, ":"), "0.0.0.0" + addr
	}
	return addr, addr
}

func loadDotEnvFiles() error {
	rewriteDoubleDashConfig()

	var configPath string
	flag.StringVar(&configPath, "config", "", "caminho do arquivo .env (opcional)")
	flag.Parse()

	switch {
	case configPath != "":
		if err := godotenv.Load(configPath); err != nil {
			return fmt.Errorf("carregar %q: %w", configPath, err)
		}
		return nil
	default:
		if err := godotenv.Load(".env"); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return fmt.Errorf("carregar .env: %w", err)
		}
		return nil
	}
}

func main() {
	if err := loadDotEnvFiles(); err != nil {
		stdlog.Fatalf("❌ Erro ao carregar arquivo de ambiente: %v", err)
	}

	// Env obrigatórias validadas em config.Load; falha se ausentes mesmo após .env.
	cfg, err := config.Load()
	if err != nil {
		stdlog.Fatalf("❌ Erro ao carregar configuração: %v", err)
	}

	log := applogger.New(cfg.LogLevel)
	slog.SetDefault(log)

	db, err := persistence.Open(cfg.DatabaseURL)
	if err != nil {
		log.Error("❌ Erro ao conectar ao banco de dados", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("✅ Conectado ao banco de dados com sucesso!")

	log.Info("🔄 Sincronizando schema do banco (AutoMigrate)...")
	if err := persistence.AutoMigrate(db); err != nil {
		log.Error("❌ Falha ao sincronizar schema do banco", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("✅ Schema do banco sincronizado!")

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

	base := httpPublicBase(cfg.HTTPAddr)
	portLabel, listenAddr := listenDisplay(cfg.HTTPAddr)
	log.Info(fmt.Sprintf("🚀 Servidor iniciado na porta %s (escutando em %s)", portLabel, listenAddr))
	log.Info(fmt.Sprintf("📝 Gin mode: %s", cfg.GinMode))
	log.Info(fmt.Sprintf("📝 Nível de log: %s", cfg.LogLevel))
	log.Info(fmt.Sprintf("📦 Versão: %s | serviço: %s", version.Version, version.Service))
	log.Info(fmt.Sprintf("🔗 Health check: %s/health", base))
	log.Info(fmt.Sprintf("📚 Documentação: %s/docs", base))
	log.Info(fmt.Sprintf("📄 OpenAPI: %s/openapi.yaml", base))

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("❌ Erro no servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("🛑 Sinal de encerramento recebido; finalizando com segurança...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("❌ Encerramento gracioso falhou", slog.Any("error", err))
	}
	log.Info("✅ Servidor encerrado.")
}
