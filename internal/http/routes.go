package handlers

import (
	"effective-project/internal/http/handlers"
	"effective-project/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	logger *slog.Logger,
	userService service.UserService,
	paymentService service.PaymentService,
	subscriptionService service.SubscriptionService,
	serviceService service.ServiceService,
	categoryService service.CategoryService,
) {
	userHandler := handlers.NewUserHandler(userService, logger)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, logger)
	serviceHandler := handlers.NewServiceHandler(serviceService, logger)
	paymentHandler := handlers.NewPaymentHandlers(paymentService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)

	userHandler.RegisterRoutes(router)
	subscriptionHandler.RegisterRoutes(router)
	serviceHandler.RegisterRoutes(router)
	paymentHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
}
