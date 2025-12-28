package api

import (
	"project/common/middleware"
	"project/services/user/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(r *gin.Engine, userHandler *handler.UserHandler, log *zap.Logger) {
	r.Use(middleware.GinLogger(log))
	r.Use(middleware.GinRecovery(log, true))

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Register)
			users.GET("/:id", userHandler.Login)
		}
	}
}
