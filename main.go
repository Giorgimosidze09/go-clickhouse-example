// package main

// import (
// 	"log"

// 	"go-clickhouse-example/config"
// 	_ "go-clickhouse-example/docs"
// 	"go-clickhouse-example/routes"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/cors"
// 	files "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // @title Your API Title
// // @version 1.0
// // @description Your API description.
// // @securityDefinitions.apikey BearerAuth
// // @in header
// // @name Authorization
// func main() {

// 	cfg := config.LoadConfig()

// 	router = routes.SetupRouter()

// 	// Swagger setup
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

// 	// Swagger JSON for custom handling
// 	router.GET("/swagger.json", func(c *gin.Context) {
// 		c.File("./docs/swagger.json")
// 	})

// 	log.Printf("Starting server on %s", cfg.ServerPort)
// 	if err := router.Run(cfg.ServerPort); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}

// }

package main

import (
	"log"
	"net/http"

	"go-clickhouse-example/config"
	_ "go-clickhouse-example/docs"
	"go-clickhouse-example/routes"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
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
	// Load configuration (e.g., server port)
	cfg := config.LoadConfig()

	// Create a new Gin router
	router := routes.SetupRouter()

	// Apply custom CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend URL
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Swagger setup (if you are using Swagger for API docs)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler)) // Swagger route
	router.GET("/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json") // Swagger JSON route
	})

	// Start the server with CORS handler
	log.Printf("Starting server on %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
