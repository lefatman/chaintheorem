package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/mvp-repo/internal/app"
	"example.com/mvp-repo/internal/config"
)

const (
	serverConfigPath   = "config/server.json"
	gameplayConfigPath = "config/gameplay.json"
)

func main() {
	serverCfg, err := config.LoadServerConfig(serverConfigPath)
	if err != nil {
		log.Fatalf("load server config: %v", err)
	}
	_, err = config.LoadGameplayConfig(gameplayConfigPath)
	if err != nil {
		log.Fatalf("load gameplay config: %v", err)
	}

	application, err := app.New(serverCfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	httpServer := &http.Server{
		Addr:    serverCfg.HTTP.ListenAddr,
		Handler: healthMux,
	}

	wsMux := http.NewServeMux()
	wsMux.Handle(serverCfg.WS.Path, application.Gateway)
	wsServer := &http.Server{
		Addr:    serverCfg.WS.ListenAddr,
		Handler: wsMux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()
	go func() {
		if serverCfg.WS.TLS.Enabled {
			errCh <- wsServer.ListenAndServeTLS(serverCfg.WS.TLS.CertFile, serverCfg.WS.TLS.KeyFile)
			return
		}
		errCh <- wsServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown requested")
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdownCtx)
	_ = wsServer.Shutdown(shutdownCtx)
}
