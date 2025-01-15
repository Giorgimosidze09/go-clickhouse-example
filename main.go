package main

import (
	"log"

	"go-clickhouse-example/config"
	_ "go-clickhouse-example/docs"
	"go-clickhouse-example/routes"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Your API Title
// @version 1.0
// @description Your API description.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	cfg := config.LoadConfig()

	router := routes.SetupRouter()

	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Swagger JSON for custom handling
	router.GET("/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	log.Printf("Starting server on %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
