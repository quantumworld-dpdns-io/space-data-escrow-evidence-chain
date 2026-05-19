package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/api"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/config"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/middleware"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/repo/memory"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/telemetry"
)

func main() {
	cfg := config.Load()
	svc := service.New(memory.NewEvidenceRepo(), memory.NewCustodyRepo(), memory.NewAuditRepo())
	r := api.NewRouter(svc, cfg.APIKey, map[string]string{
		"version":    cfg.Version,
		"commit":     cfg.Commit,
		"build_date": cfg.BuildDate,
	})
	reg := telemetry.NewRegistry()

	handler := middleware.Chain(
		r.Handler(),
		middleware.RequestID(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.Timeout(cfg.ReadTimeout),
		telemetry.HTTPMetrics(reg),
	)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: handler, ReadTimeout: cfg.ReadTimeout, WriteTimeout: cfg.WriteTimeout}
	go func() {
		log.Printf("starting API server on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimout)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
