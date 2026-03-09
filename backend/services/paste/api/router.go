package api

import (
	"project/common/auth"
	"project/common/middleware"
	"project/services/paste/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(r *gin.Engine, pasteHandler *handler.PasteHandler, jwtManager *auth.JWTManager, log *zap.Logger) {
	r.Use(middleware.GinLogger(log))
	r.Use(middleware.GinRecovery(log, true))

	v1 := r.Group("/api/v1")
	{
		me := v1.Group("/me")
		me.Use(middleware.JWTAuthMiddleware(jwtManager))
		{
			me.GET("/pastes", pasteHandler.ListMine)
		}

		pastes := v1.Group("/pastes")
		pastes.Use(middleware.JWTAuthMiddleware(jwtManager))
		{
			pastes.POST("", pasteHandler.Create)
			pastes.GET("/:id", pasteHandler.Get)
			pastes.PUT("/:id", pasteHandler.Update)
		}
	}
}
