package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gravitum_test_task/internal/config"
	"github.com/sonikq/gravitum_test_task/internal/handler/user_management"
	"github.com/sonikq/gravitum_test_task/internal/server/middleware"
	"github.com/sonikq/gravitum_test_task/internal/service"
	"github.com/sonikq/gravitum_test_task/pkg/logger"
	"net/http"
)

type Handler struct {
	UserManagement *user_management.Handler
}

type Option struct {
	Conf    config.Config
	Logger  *logger.Logger
	Service *service.Service
}

func NewRouter(option Option) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestResponseLogger(option.Logger))

	h := Handler{UserManagement: user_management.New(&user_management.HandlerConfig{
		Config:  option.Conf,
		Logger:  option.Logger,
		Service: option.Service,
	}),
	}

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "I am alive!",
		})
	})

	// Authorized routes
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", h.UserManagement.CreateUser)
		userGroup.GET("/:id", h.UserManagement.GetUser)
		userGroup.PUT("/:id", h.UserManagement.UpdateUser)
		userGroup.DELETE("/:id", h.UserManagement.DeleteUser)

	}

	return router
}
