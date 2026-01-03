package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Driver string
	Source string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

// LoadConfig 从环境变量加载配置（后期可改为配置文件）
func LoadConfig() *Config {
	// 💡 优化：支持从单独的环境变量组合数据库连接串
	// 这样可以配合 Kubernetes ConfigMap/Secret 使用
	dbSource := getEnv("DB_SOURCE", "")
	if dbSource == "" {
		// 如果没有 DB_SOURCE，就从单独的环境变量中组合
		// 格式：postgres://user:password@host:port/dbname?sslmode=disable
		pgUser := getEnv("POSTGRES_USER", "postgres")
		pgPassword := getEnv("POSTGRES_PASSWORD", "123456")
		pgHost := getEnv("POSTGRES_HOST", "postgres")
		pgPort := getEnv("POSTGRES_PORT", "5432")
		pgDB := getEnv("POSTGRES_DB", "gopher_paste")
		pgSSLMode := getEnv("POSTGRES_SSLMODE", "disable")
		
		// 💡 修复：正确的 PostgreSQL 连接串格式
		// postgres://用户名:密码@主机:端口/数据库名?参数
		dbSource = "postgres://" + pgUser + ":" + pgPassword + "@" + 
		           pgHost + ":" + pgPort + "/" + pgDB + "?sslmode=" + pgSSLMode
	}
	
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "postgres"),
			Source: dbSource,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "gopher_redis:6379"),
			Password: getEnv("REDIS_PASSWORD", "123456"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "gopherpaste_secret_key"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
