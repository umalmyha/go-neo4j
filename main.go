package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/umalmyha/go-neo4j/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	drv, err := neo4j.NewDriverWithContext(cfg.Neo4jURL, neo4j.BasicAuth(cfg.Neo4jUsername, cfg.Neo4jPassword, ""))
	if err != nil {
		log.Fatal(err)
	}
	defer drv.Close(context.Background())

	app := fiber.New()

	app.Post("/emails", func(ctx *fiber.Ctx) error {

	})

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
