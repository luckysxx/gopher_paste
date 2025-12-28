package database

import (
	"database/sql"

	"project/common/config"

	"go.uber.org/zap"
)

// InitPostgres 初始化 PostgreSQL 数据库连接
func InitPostgres(cfg config.DatabaseConfig, log *zap.Logger) *sql.DB {
	conn, err := sql.Open(cfg.Driver, cfg.Source)
	if err != nil {
		log.Fatal("无法连接到 PostgreSQL 数据库", zap.Error(err))
		return nil
	}

	log.Info("成功连接到 PostgreSQL 数据库")
	return conn
}
