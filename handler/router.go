package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "stlab.itechart-group.com/go/food_delivery/authentication_service/docs"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/pkg/logging"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/service"
)

type Handler struct {
	logger  logging.Logger
	service *service.Service
}

func NewHandler(logger logging.Logger, service *service.Service) *Handler {
	return &Handler{logger: logger, service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(
		h.CorsMiddleware,
	)

	userNoAuth := router.Group("/users")
	{
		userNoAuth.POST("/login", h.authUser)
		userNoAuth.POST("/customer", h.createCustomer)
		userNoAuth.POST("/restorePassword", h.restorePassword)
	}

	userAuth := router.Group("/users")
	userAuth.Use(h.userIdentity)
	{
		userAuth.GET("/:id", h.getUser)
		userAuth.GET("/", h.getUsers)
		userAuth.POST("/staff", h.createStaff)
		userAuth.PUT("/", h.updateUser)
		userAuth.DELETE("/:id", h.deleteUserByID)
	}
	return router
}
