package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vaha95/golang_pet/api/shortenerpb"
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/config/db"
	"github.com/Vaha95/golang_pet/internal/grpcserver"
	"github.com/Vaha95/golang_pet/internal/handler"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/service/audit"
	"github.com/Vaha95/golang_pet/internal/service/build"
	deleteurl "github.com/Vaha95/golang_pet/internal/service/delete_url"
	"github.com/labstack/echo-contrib/v5/pprof"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	mv "github.com/Vaha95/golang_pet/internal/infrastructure/middleware"
)

func main() {
	build := build.Create()
	build.Print()

	mainConfig := config.GetMainConfig()

	dbService, err := db.InitDb(mainConfig)
	if err != nil && !errors.Is(err, db.ErrorDBDsnEmpty) {
		log.Fatal(
			fmt.Errorf("can`t init db: %w", err).Error(),
		)
	}
	isDBAllowed := false
	if dbService != nil {
		defer dbService.Close()

		isDBAllowed = dbService.Ping() == nil
	}

	cfg := config.GetConfig(dbService, isDBAllowed, mainConfig)

	l, err := getLogger()
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t init logger: %w", err).Error(),
		)
	}

	e := echo.New()

	pprof.Register(e)

	auditWorker := audit.CreateAuditWorker(cfg.Config)
	auditWorker.Run()

	e.GET(`/:id`, handler.GetURLHandler(cfg, l, auditWorker.AuditCh))
	e.POST(`/`, handler.GetSaveURLHandler(cfg, l, auditWorker.AuditCh))
	e.POST(`/api/shorten`, handler.GetSaveURLShortenHandler(cfg, l, auditWorker.AuditCh))
	e.GET(`/ping`, handler.GetPingDBHandler(dbService, l))
	e.POST(`/api/shorten/batch`, handler.GetSaveURLBatchHandler(cfg, l))
	e.GET(`/api/user/urls`, handler.GetURLByUserHandler(cfg, l))
	e.GET(`/api/imternal/stats`, handler.GetStatsHandler(cfg, l))

	deleteCh := make(chan DTO.DeleteBatch)
	listener := deleteurl.GetDeleteUrlListener(cfg, deleteCh, l)
	go listener()
	e.DELETE(`/api/user/urls`, handler.GetDeleteURLHandler(cfg, deleteCh, l))

	err = mv.AddMiddlewares(cfg, e, l)
	if err != nil {
		log.Fatal(
			fmt.Errorf("can`t start Web server: %w", err).Error(),
		)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	grpcSvc := startGRPC(cfg, auditWorker.AuditCh, l)

	srvCh := make(chan *http.Server, 1)
	go func() {
		srvCh <- serveStart(e, l, cfg.Config)
	}()
	srv := <-srvCh
	<-ctx.Done()
	if grpcSvc != nil {
		grpcSvc.GracefulStop()
	}
	srv.Shutdown(ctx)
}

func startGRPC(cfg config.StorageConfig, auditCh chan DTO.BaseAuditItem, l *zap.SugaredLogger) *grpc.Server {
	grpcListener, err := net.Listen("tcp", cfg.Config.GRPCPort)
	if err != nil {
		l.Errorf("can`t listen on gRPC port %s: %v", cfg.Config.GRPCPort, err)
		return nil
	}

	grpcSvc := grpc.NewServer()
	shortenerpb.RegisterShortenerServiceServer(grpcSvc, grpcserver.NewServer(cfg, auditCh))
	go func() {
		if err := grpcSvc.Serve(grpcListener); err != nil {
			l.Errorf("gRPC server error: %v", err)
		}
	}()
	l.Infof("gRPC server listening on %s", cfg.Config.GRPCPort)

	return grpcSvc
}

func serveStart(e *echo.Echo, l *zap.SugaredLogger, cfg config.Config) *http.Server {
	if cfg.EnableHttps {
		return serveStartTLS(e, l, cfg)
	}

	return serveStartDefault(e, cfg)
}

func serveStartDefault(e *echo.Echo, cfg config.Config) *http.Server {
	srv := &http.Server{
		Addr:    cfg.ListenHost,
		Handler: e,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(
				fmt.Errorf("can`t start Web server: %w", err).Error(),
			)
		}
	}()

	return srv
}

func serveStartTLS(e *echo.Echo, l *zap.SugaredLogger, cfg config.Config) *http.Server {
	srv := &http.Server{
		Addr:    cfg.TLSAddress,
		Handler: e,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil && err != http.ErrServerClosed {
			l.Fatal(err.Error())
		}
	}()

	return srv
}

func getLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("%w", mv.ErrorZapLoggerInitialize)
	}
	defer logger.Sync()

	return logger.Sugar(), nil
}
