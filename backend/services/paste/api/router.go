package api

import (
	"project/common/middleware"
	"project/services/paste/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(r *gin.Engine, pasteHandler *handler.PasteHandler, log *zap.Logger) {
	r.Use(middleware.GinLogger(log))
	r.Use(middleware.GinRecovery(log, true))

	v1 := r.Group("/api/v1")
	{
		pastes := v1.Group("/pastes")
		{
			pastes.POST("", pasteHandler.Create)
			pastes.GET("/:id", pasteHandler.Get)
		}
	}
}
