package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"portaldata-api/internal/pkg/config"
	"portaldata-api/internal/pkg/database"

	offer "portaldata-api/internal/modules/offer"
	order "portaldata-api/internal/modules/order"
	products "portaldata-api/internal/modules/products"
	user "portaldata-api/internal/modules/user"
	warehouse "portaldata-api/internal/modules/warehouse"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	// Middlewares
	// router.Use(middleware.BruteForceMiddleware())

	userService := user.NewService(db.DB)
	authService := user.NewAuthService(userService)
	authMiddleware := authService.AuthMiddleware()

	productsService := products.NewService(db)
	productsHandlers := products.NewHandlers(productsService)

	offerService := offer.NewService(db)
	offerHandlers := offer.NewHandlers(offerService)

	orderService := order.NewService(db)
	orderHandlers := order.NewHandlers(orderService)

	warehouseService := warehouse.NewService(db)
	warehouseHandlers := warehouse.NewHandlers(warehouseService)

	// Основной endpoint для проверки доступности
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"data": gin.H{
				"message":  "API ПорталДанных.РФ доступен",
				"version":  "v1",
				"database": "MySQL",
				"status":   "running",
			},
		})
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты (без авторизации)
	publicGroup := router.Group("/api/v1")

	// Публичные офферы
	offer.RegisterPublicRoutes(publicGroup, offerHandlers)

	// Защищенные маршруты (с авторизацией)
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(authMiddleware)
	products.RegisterRoutes(apiGroup, productsHandlers)
	offer.RegisterRoutes(apiGroup, offerHandlers)
	order.RegisterRoutes(apiGroup, orderHandlers)
	warehouse.RegisterRoutes(apiGroup, warehouseHandlers)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Сервер останавливается...")
}
