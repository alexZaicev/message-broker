package main

import (
	"context"
	"injector/internal/adapters/proxy"
	"injector/internal/drivers/codes"
	"injector/internal/drivers/logging"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func run() int {
	server, err := proxy.NewInjectorProxyServer(8080)
	if err != nil {
		slog.Error("failed to create proxy server", logging.WithError(err))
		return codes.Failure
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err = server.Serve(); err != nil {
		slog.Error("failed to start proxy server", logging.WithError(err))
		return codes.Failure
	}

	slog.Info("proxy server started")

	<-ctx.Done()
	stop()

	slog.Info("stopping proxy server")
	server.Stop()
	slog.Info("proxy server stopped")

	return codes.Success
}

func main() {
	logging.SetupJSONLogger()
	os.Exit(run())
}
