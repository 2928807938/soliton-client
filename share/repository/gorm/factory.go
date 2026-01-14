package gorm

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseType 数据库类型
type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgres"
	SQLite     DatabaseType = "sqlite"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            DatabaseType    // 数据库类型
	Host            string          // 主机地址
	Port            int             // 端口
	Database        string          // 数据库名
	Username        string          // 用户名
	Password        string          // 密码
	Charset         string          // 字符集（MySQL）
	SSLMode         string          // SSL 模式（PostgreSQL）
	MaxIdleConns    int             // 最大空闲连接数
	MaxOpenConns    int             // 最大打开连接数
	ConnMaxLifetime time.Duration   // 连接最大生命周期
	ConnMaxIdleTime time.Duration   // 连接最大空闲时间
	LogLevel        logger.LogLevel // 日志级别
	SlowThreshold   time.Duration   // 慢查询阈值
}

// DefaultConfig 默认配置
func DefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:            PostgreSQL,
		Host:            "localhost",
		Port:            5432,
		Charset:         "utf8mb4",
		SSLMode:         "disable",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
		LogLevel:        logger.Info,
		SlowThreshold:   200 * time.Millisecond,
	}
}

// DatabaseFactory 数据库工厂
type DatabaseFactory struct {
	config *DatabaseConfig
}

// NewDatabaseFactory 创建数据库工厂
func NewDatabaseFactory(config *DatabaseConfig) *DatabaseFactory {
	if config == nil {
		config = DefaultConfig()
	}
	return &DatabaseFactory{config: config}
}

// Create 创建数据库连接
func (f *DatabaseFactory) Create() (*gorm.DB, error) {
	dialector, err := f.getDialector()
	if err != nil {
		return nil, err
	}

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(f.config.LogLevel),
	}

	// 打开数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 SQL 连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(f.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(f.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(f.config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(f.config.ConnMaxIdleTime)

	// 注册审计回调
	RegisterAuditCallbacks(db)

	return db, nil
}

// getDialector 根据数据库类型获取 Dialector
func (f *DatabaseFactory) getDialector() (gorm.Dialector, error) {
	switch f.config.Type {
	case MySQL:
		return f.getMySQLDialector(), nil
	case PostgreSQL:
		return f.getPostgresDialector(), nil
	case SQLite:
		return f.getSQLiteDialector(), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", f.config.Type)
	}
}

// getMySQLDialector 获取 MySQL Dialector
func (f *DatabaseFactory) getMySQLDialector() gorm.Dialector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		f.config.Username,
		f.config.Password,
		f.config.Host,
		f.config.Port,
		f.config.Database,
		f.config.Charset,
	)
	return mysql.Open(dsn)
}

// getPostgresDialector 获取 PostgreSQL Dialector
func (f *DatabaseFactory) getPostgresDialector() gorm.Dialector {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		f.config.Host,
		f.config.Port,
		f.config.Username,
		f.config.Password,
		f.config.Database,
		f.config.SSLMode,
	)
	return postgres.Open(dsn)
}

// getSQLiteDialector 获取 SQLite Dialector
func (f *DatabaseFactory) getSQLiteDialector() gorm.Dialector {
	return sqlite.Open(f.config.Database)
}

// CreateWithDSN 使用 DSN 创建数据库连接
func CreateWithDSN(dbType DatabaseType, dsn string, config *DatabaseConfig) (*gorm.DB, error) {
	if config == nil {
		config = DefaultConfig()
	}

	var dialector gorm.Dialector
	switch dbType {
	case MySQL:
		dialector = mysql.Open(dsn)
	case PostgreSQL:
		dialector = postgres.Open(dsn)
	case SQLite:
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}

	// 打开数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 SQL 连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// 注册审计回调
	RegisterAuditCallbacks(db)

	return db, nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
