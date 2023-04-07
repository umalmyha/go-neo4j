package main

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/umalmyha/go-neo4j/internal/config"
	"github.com/umalmyha/go-neo4j/internal/handler"
	"github.com/umalmyha/go-neo4j/internal/storage"
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

	emailStorage := storage.NewNeo4jEmailStorage(drv)
	emailHandler := handler.NewEmailHandler(emailStorage)

	e := echo.New()
	e.POST("/users", emailHandler.CreateUser)
	e.GET("/users/:id", emailHandler.FindByID)
	e.POST("/emails", emailHandler.CreateEmail)

	if err := e.Start(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.Fatal(err)
	}
}
