package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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

	timeout := time.Second * 15 // TODO: configurable
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Shutting down...")
}
