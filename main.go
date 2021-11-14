package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/SignorMercurio/cncamp_homework/httpserver"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s [listen address]", os.Args[0])
	}
	os.Setenv("VERSION", "1.2.0")

	srv := httpserver.NewServer(os.Args[1])
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Graceful termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// from ConfigMap httpserver-config
	graceTimeout, err := strconv.Atoi(os.Getenv("GRACE_TIMEOUT"))
	if err != nil {
		log.Println("Failed to read GRACE_TIMEOUT from env, default to 30s")
		graceTimeout = 30
	}
	timeout := time.Second * time.Duration(graceTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Shutting down...")
}
