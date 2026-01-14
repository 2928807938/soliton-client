// Package bom 是 Bill of Materials 模块，用于统一管理所有依赖版本
// 其他模块通过 replace 指令引用此模块，自动继承依赖版本
package bom

import (
	// Hertz HTTP 框架
	_ "github.com/cloudwego/hertz/pkg/app"
	_ "github.com/cloudwego/hertz/pkg/app/server"
	_ "github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/cloudwego/hertz/pkg/protocol/consts"

	// Kitex RPC 框架
	_ "github.com/cloudwego/kitex/client"
	_ "github.com/cloudwego/kitex/pkg/klog"
	_ "github.com/cloudwego/kitex/server"

	// 通用工具
	_ "github.com/bytedance/sonic"
	_ "github.com/google/uuid"

	// 数据库
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"

	// 配置管理
	_ "github.com/spf13/viper"

	// 验证器
	_ "github.com/go-playground/validator/v10"

	// 缓存
	_ "github.com/redis/go-redis/v9"
)
