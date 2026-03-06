package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"project/common/auth"
	"project/common/config"
	"project/common/database"
	"project/common/logger"
	"project/services/user/api"
	"project/services/user/db"

	// _ "project/services/user/docs"
	"project/services/user/handler"
	"project/services/user/repository"
	"project/services/user/service"
)

// @title           GopherPaste User Service
// @version         1.0
// @description     用户中心服务，提供注册、登录功能
// @host            localhost:8081
// @BasePath        /api/v1
func main() {
	// 初始化基础设施
	log := logger.NewLogger("user")
	defer log.Sync()

	cfg := config.LoadConfig()
	conn := database.InitPostgres(cfg.Database, log)
	queries := db.New(conn)

	// 依赖注入
	userRepo := repository.NewUserRepository(queries)
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret)
	userSvc := service.NewUserService(userRepo, jwtManager, log)
	userHandler := handler.NewUserHandler(userSvc, log)

	// 路由
	r := gin.New()
	api.SetupRouter(r, userHandler, log)
	r.Run(":" + cfg.Server.Port)
}
