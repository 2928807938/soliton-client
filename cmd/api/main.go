package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"soliton-client/api/handlers"
)

func main() {
	// 初始化数据库
	db, err := initDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化 Redis
	_ = initRedis()

	// 初始化用户服务HTTP客户端
	userClient := initUserServiceClient()

	// 初始化 Hertz 服务器
	port := getEnv("PORT", "8080")
	h := server.Default(
		server.WithHostPorts(":"+port),
		server.WithMaxRequestBodySize(10*1024*1024), // 10MB
	)

	// 注册路由
	registerRoutes(h, db, userClient)

	// 启动服务
	log.Printf("服务启动在 :%s", port)
	h.Spin()
}

// registerRoutes 注册所有路由
func registerRoutes(h *server.Hertz, db *gorm.DB, userClient *handlers.UserServiceClient) {
	// 健康检查
	h.GET("/health", healthCheck(db))

	// API v1 路由组
	v1 := h.Group("/api/v1")
	{
		// Ping 测试
		v1.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
			c.JSON(consts.StatusOK, map[string]interface{}{
				"message": "pong",
			})
		})

		// 用户相关路由
		handlers.RegisterUserRoutes(v1, userClient)
	}
}

// healthCheck 健康检查处理器
func healthCheck(db *gorm.DB) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		health := map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		}

		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			health["status"] = "error"
			health["database"] = "disconnected"
			c.JSON(consts.StatusServiceUnavailable, health)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			health["status"] = "error"
			health["database"] = "ping failed"
			c.JSON(consts.StatusServiceUnavailable, health)
			return
		}

		health["database"] = "connected"
		c.JSON(consts.StatusOK, health)
	}
}

func initDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "soliton-client")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func initRedis() *redis.Client {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := 0

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Redis 连接失败: %v (使用 %s)", err, addr)
	} else {
		log.Printf("Redis 连接成功: %s", addr)
	}

	return client
}

func initUserServiceClient() *handlers.UserServiceClient {
	// 连接到 universal-service-user 服务
	baseURL := getEnv("USER_SERVICE_URL", "http://universal-service-user:8080")
	tenantID := getEnv("USER_SERVICE_TENANT_ID", "")

	// 如果没有配置 tenant_id，需要先注册应用获取 tenant_id
	if tenantID == "" {
		log.Println("警告: 未配置 USER_SERVICE_TENANT_ID，请先注册应用获取 tenant_id")
	}

	userClient := handlers.NewUserServiceClient(baseURL, tenantID)
	log.Printf("用户服务 HTTP 客户端初始化成功，连接到: %s", baseURL)
	return userClient
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
