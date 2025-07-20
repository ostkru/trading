package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"portaldata-api/internal/api"
	"portaldata-api/internal/config"
	"portaldata-api/internal/database"
	"portaldata-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	router := gin.Default()

	// Включаю защиту от брутфорса до всех остальных middleware
	router.Use(api.BruteForceMiddleware())

	// Инициализация сервисов с БД
	offerService := services.NewOfferService(db)
	orderService := services.NewOrderService(db)
	warehouseService := services.NewWarehouseService(db)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db.DB)
	authService := api.NewAuthService(userService)

	// Инициализация обработчиков
	offerHandlers := api.NewOfferHandlers(offerService)
	orderHandlers := api.NewOrderHandlers(orderService)
	warehouseHandlers := api.NewWarehouseHandlers(warehouseService)
	productHandlers := api.NewProductHandlers(productService)

	// Публикация OpenAPI YAML
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("/var/www/go/openapi.yaml")
	})

	// Swagger UI редирект (можно заменить на локальный Swagger UI при необходимости)
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://petstore.swagger.io/?url=http://"+c.Request.Host+"/openapi.yaml")
	})

	// Группировка API
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(authService.Middleware())
	api.RegisterOfferRoutes(apiGroup, offerHandlers)
	api.RegisterOrderRoutes(apiGroup, orderHandlers)
	api.RegisterWarehouseRoutes(apiGroup, warehouseHandlers)
	api.RegisterProductRoutes(apiGroup, productHandlers)

	router.Static("/swagger", "/var/www/go/swagger")

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		if c.IsAborted() {
			return
		}

		c.Next()
	})

	go func() {
		log.Printf("Starting server on port 8095")
		if err := router.Run(":8095"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
