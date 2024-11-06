package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/EmilioCliff/crocheted-ecommerce/backend/docs/statik"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/handlers"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

// @title           Crocheted Ecommerce API
// @version         1.0
// @description     This is a sample server for the ecommerce website.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Mozilla Public License Version 2.0
// @license.url   https://www.mozilla.org/en-US/MPL/2.0/

// @host      localhost:3030
// @BasePath  /api/v1

// @securityDefinitions.basic  BearerToken

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	config, err := pkg.LoadConfig("../../.envs/.local/")
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
