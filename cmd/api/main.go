// @title Trading API
// @version 1.0
// @description API для торговой платформы PortalData
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8095
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/gin-gonic/gin"

	"portaldata-api/internal/pkg/config"
	"portaldata-api/internal/pkg/database"
	"portaldata-api/internal/pkg/middleware"

	categories "portaldata-api/internal/modules/categories"
	characteristics "portaldata-api/internal/modules/characteristics"
	offer "portaldata-api/internal/modules/offer"
	order "portaldata-api/internal/modules/order"
	products "portaldata-api/internal/modules/products"
	ratelimit "portaldata-api/internal/modules/ratelimit"
	search "portaldata-api/internal/modules/search"
	tariff "portaldata-api/internal/modules/tariff"
	user "portaldata-api/internal/modules/user"
	warehouse "portaldata-api/internal/modules/warehouse"

	"github.com/gin-contrib/cors"
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

	// Инициализация Redis клиента для кэширования
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	router := gin.Default()

	// Middlewares
	// router.Use(middleware.BruteForceMiddleware())

	userService := user.NewService(db.DB)
userHandlers := user.NewHandlers(userService)
	authService := user.NewAuthService(userService)
	authMiddleware := authService.AuthMiddleware()

	productsService := products.NewService(db)
	productsHandlers := products.NewHandlers(productsService)

	categoriesService := categories.NewService(db, redisClient)
	categoriesHandlers := categories.NewHandlers(categoriesService)

	characteristicsService := characteristics.NewService(db)
	characteristicsHandlers := characteristics.NewHandlers(characteristicsService)

	offerService := offer.NewService(db)
	offerHandlers := offer.NewHandlers(offerService)

	orderService := order.NewService(db)
	orderHandlers := order.NewHandlers(orderService)

	warehouseService := warehouse.NewService(db)
	warehouseHandlers := warehouse.NewHandlers(warehouseService)

	tariffService := tariff.NewService(db.DB)
	tariffHandlers := tariff.NewHandler(tariffService)

	// Инициализация поиска
	searchService := search.NewService("http://localhost:9200", db.DB, redisClient)
	searchHandlers := search.NewHandlers(searchService)

	// Инициализация синхронизации OpenSearch для offers
	opensearchSync := offer.NewOpenSearchSync(searchService)
	offerService.SetOpenSearchSync(opensearchSync)

	// Инициализация Redis Rate Limiting
	redisRateLimitService := ratelimit.NewRedisRateLimitService(
		"127.0.0.1:6379", // Redis адрес
		"",               // Redis пароль
		0,                // Redis база данных
		db.DB,            // База данных для получения лимитов тарифов
	)

	// Основной endpoint для проверки доступности
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"message":  "API ПорталДанных.РФ доступен",
				"version":  "v1",
				"database": "MySQL",
				"status":   "running",
			},
		})
	})

	// Страница для просмотра API в браузере
	router.GET("/browser", func(c *gin.Context) {
		c.File("browser_view.html")
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API маршруты (с rate limiting)
	apiGroup := router.Group("/api")

	// Применяем rate limiting middleware к API маршрутам ПЕРВЫМ (для быстрой блокировки)
	if redisRateLimitService == nil {
		log.Printf("⚠️  Предупреждение: Redis Rate Limiting недоступен, используется MySQL-based rate limiting")
		apiGroup.Use(ratelimit.RateLimitMiddleware(ratelimit.NewService(db.DB)))
	} else {
		log.Printf("✅ Redis Rate Limiting подключен успешно")
		apiGroup.Use(ratelimit.RedisRateLimitMiddleware(redisRateLimitService))
	}

	// Публичные офферы (без авторизации) - регистрируем ДО authMiddleware
	offer.RegisterPublicRoutes(apiGroup, offerHandlers)
	// Регистрируем роуты пользователей (публичные и защищенные)
	user.RegisterRoutes(apiGroup, userHandlers)
	// Регистрируем роуты поиска (публичные)
	search.RegisterRoutes(apiGroup, searchHandlers)
	// Регистрируем роуты категорий (публичные)
	categories.RegisterRoutes(apiGroup, categoriesHandlers)
	// Регистрируем роуты характеристик (публичные)
	characteristics.RegisterRoutes(apiGroup, characteristicsHandlers)
	// Применяем auth middleware к защищенным маршрутам (с кэшированием в Redis)
	apiGroup.Use(authMiddleware)
	
	// Применяем INN verification middleware для POST/PUT/DELETE запросов
	apiGroup.Use(middleware.INNVerificationMiddleware(db.DB))

	// Защищенные маршруты (с авторизацией)
	products.RegisterRoutes(apiGroup, productsHandlers)
	offer.RegisterRoutes(apiGroup, offerHandlers)
	order.RegisterRoutes(apiGroup, orderHandlers)
	warehouse.RegisterRoutes(apiGroup, warehouseHandlers)
	tariff.RegisterRoutes(apiGroup, tariffHandlers)

	// Redis Rate Limiting API маршруты (публичные для мониторинга)
	if redisRateLimitService != nil {
		redisHandler := ratelimit.NewRedisHandler(redisRateLimitService)
		rateLimitGroup := router.Group("/rate-limit")
		ratelimit.SetupRedisRoutes(rateLimitGroup, redisHandler)
		log.Printf("✅ Redis Rate Limiting API маршруты зарегистрированы")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"Content-Length", "X-RateLimit-Limit-Minute", "X-RateLimit-Limit-Day", "X-RateLimit-Remaining-Minute", "X-RateLimit-Remaining-Day"},
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
