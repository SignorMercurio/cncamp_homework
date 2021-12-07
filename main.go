package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/SignorMercurio/cncamp_homework/httpserver"
	"github.com/SignorMercurio/cncamp_homework/logger"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	sugar := zap.S()

	if len(os.Args) != 2 {
		sugar.Fatalf("Usage: %s [listen address]", os.Args[0])
	}
	os.Setenv("VERSION", "1.5.2")

	addr := os.Args[1]
	sugar.Debugw("Creating new server...", "address", addr)
	srv := httpserver.NewServer(addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			sugar.Error(err)
		}
	}()

	// Graceful termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	sugar.Debug("Received SIGTERM, gracefully terminating...")
	// from ConfigMap httpserver-config
	graceTimeout, err := strconv.Atoi(os.Getenv("GRACE_TIMEOUT"))
	if err != nil {
		graceTimeout = 30
		sugar.Warnw(
			"Failed to read GRACE_TIMEOUT from env, using default",
			"default",
			graceTimeout,
		)
	}
	sugar.Debugw(
		"Set graceful termination timeout successfully",
		"timeout",
		graceTimeout,
	)
	timeout := time.Second * time.Duration(graceTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	srv.Shutdown(ctx)
	sugar.Info("Shutting down...")
}
