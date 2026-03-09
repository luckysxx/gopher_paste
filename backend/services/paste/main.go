package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"project/common/auth"
	"project/common/config"
	"project/common/database"
	"project/common/logger"

	"project/services/paste/api"
	_ "project/services/paste/docs"
	"project/services/paste/handler"
	"project/services/paste/repository"
	"project/services/paste/service"
)

// @title           GopherPaste API
// @version         1.0
// @description     一个类似于 Pastebin 的代码分享服务
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// 初始化基础设施
	log := logger.NewLogger("paste")
	defer log.Sync()

	cfg := config.LoadConfig()
	conn := database.InitPostgres(cfg.Database, log)
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret)

	// 依赖注入
	pasteRepo := repository.NewPasteRepository(conn)
	pasteSvc := service.NewPasteService(pasteRepo, log)
	pasteHandler := handler.NewPasteHandler(pasteSvc, log)

	// 路由
	r := gin.New()
	api.SetupRouter(r, pasteHandler, jwtManager, log)
	r.Run(":" + cfg.Server.Port)
}
