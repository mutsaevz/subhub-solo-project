package main

import (
	"effective-project/internal/cache"
	"effective-project/internal/config"
	handlers "effective-project/internal/http"
	"effective-project/internal/models"
	"effective-project/internal/redis"
	"effective-project/internal/repository"
	"effective-project/internal/service"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// инициализация логгера
	logger := config.InitLogger()

	// gin router
	router := gin.New()
	router.Use(gin.Recovery())

	// healthcheck
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	db := config.SetUpDatabaseConnection(logger)
	if db == nil {
		logger.Error("database is nil")
		return
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Subscription{},
		&models.Service{},
		&models.Payment{},
		&models.Category{},
		&models.Order{},
	); err != nil {
		logger.Error("failed to migrate database", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("migrations completed")

	// redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redisClient := redis.New(redisAddr)

	ttl := 5 * time.Minute

	userCache := cache.NewUserRedisCache(redisClient)
	paymentCache := cache.NewPaymentRedisCache(redisClient)
	subscriptionCache := cache.NewSubscriptionRedisCache(redisClient)
	serviceCache := cache.NewServiceCache(redisClient, ttl)
	categoryCache := cache.NewCategoryRedisCache(redisClient)

	// repositories
	userRepo := repository.NewUserRepository(db, logger)
	subscriptionRepo := repository.NewSubscriptionRepository(db, logger)
	serviceRepo := repository.NewServiceRepository(db, logger)
	paymentRepo := repository.NewPaymentRepository(db, logger)
	categoryRepo := repository.NewCategoryRepository(db, logger)

	// services
	userService := service.NewUserService(
		userRepo,
		userCache,
		logger,
	)

	subscriptionService := service.NewSubscriptionService(
		subscriptionRepo,
		serviceRepo,
		paymentRepo,
		subscriptionCache,
		logger,
	)

	serviceService := service.NewServiceService(
		serviceRepo,
		serviceCache,
		logger,
	)

	paymentService := service.NewPaymentService(
		paymentRepo,
		paymentCache,
		logger,
	)

	categoryService := service.NewCategoryService(
		categoryRepo,
		categoryCache,
		logger,
	)

	api := router.Group("")
	// handlers / routes
	handlers.RegisterRoutes(
		api,
		logger,
		userService,
		paymentService,
		subscriptionService,
		serviceService,
		categoryService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("application started successfully", "port", port)

	if err := router.Run(":" + port); err != nil {
		logger.Error("ошибка запуска сервера", slog.Any("error", err))
	}
}
