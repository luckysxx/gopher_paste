package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"project/common/cache"
	"project/common/config"
	"project/common/database"
	"project/common/logger"

	"project/services/paste/api"
	"project/services/paste/db"
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
	queries := db.New(conn)
	rdb := cache.InitRedis(cfg.Redis, log)
	cacheInstance := cache.NewCache(rdb)

	// 依赖注入
	pasteRepo := repository.NewPasteRepository(queries)
	pasteSvc := service.NewPasteService(pasteRepo, cacheInstance, log)
	pasteHandler := handler.NewPasteHandler(pasteSvc, log)

	// 路由
	r := gin.New()
	api.SetupRouter(r, pasteHandler, log)
	r.Run(":" + cfg.Server.Port)
}
