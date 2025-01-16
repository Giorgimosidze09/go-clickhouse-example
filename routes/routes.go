package routes

import (
	"go-clickhouse-example/config"
	"go-clickhouse-example/handlers"
	"go-clickhouse-example/middleware" // Import the middleware
	"go-clickhouse-example/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	cfg := config.LoadConfig()

	// Initialize services
	dbService := services.NewDBService(cfg.ClickHouse)
	dbService.CreateTable()
	natsService := services.NewNATSService(cfg.NATSURL, cfg.StreamName, cfg.SubjectName)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(dbService, natsService)
	authService := services.NewAuthService(dbService)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize the router
	router := gin.Default()

	// Public routes for user registration and login
	router.POST("/register", authHandler.RegisterUser)
	router.POST("/login", authHandler.LoginUser)

	// Protected routes (Require authentication and authorization)
	// Apply AuthMiddleware to secure the routes and RBACMiddleware for role-based access control
	router.POST("/items", middleware.AuthMiddleware(), middleware.RBACMiddleware("admin"), itemHandler.CreateItem)
	router.GET("/items", middleware.AuthMiddleware(), itemHandler.GetItems)

	router.GET("/items/:id", middleware.AuthMiddleware(), itemHandler.GetItem)
	router.PUT("/items/:id", middleware.AuthMiddleware(), middleware.RBACMiddleware("admin"), itemHandler.UpdateItem)
	router.DELETE("/items/:id", middleware.AuthMiddleware(), middleware.RBACMiddleware("admin"), itemHandler.DeleteItem)

	return router
}
