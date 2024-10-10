package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/handlers"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

func main() {
	config, err := pkg.LoadConfig("/app/.envs/.local/")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	tokenMaker, err := pkg.NewPaseto(config.TOKEN_SYMMETRY_KEY)
	if err != nil {
		log.Fatalf("failed to create token maker: %v", err)
	}

	store := mysql.NewStore(config, tokenMaker)

	err = store.Open()
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}

	server := handlers.NewHttpServer(tokenMaker, config)

	server.SetDependencies(store)

	log.Println("Starting server at port: ", config.HTTP_PORT)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-quit

	log.Println("Shutdown...")

	if err := server.Close(); err != nil {
		log.Fatalf("failed to close server: %v", err)
	}

	if err := store.Close(); err != nil {
		log.Fatalf("failed to close store: %v", err)
	}

	log.Println("Server exiting")
}
