package main

import (
	"log"
	"os"

	"github.com/SignorMercurio/cncamp_homework/httpserver"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s [listen address]", os.Args[0])
	}
	setEnv()

	log.Fatal(httpserver.NewServer(os.Args[1]))
}

func setEnv() {
	os.Setenv("VERSION", "1.1.0")
}
