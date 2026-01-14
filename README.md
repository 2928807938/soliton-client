# soliton-client

基于 Go 语言的领域驱动设计（DDD）项目，采用多模块工作区（Go Workspace）+ BOM 依赖管理。

## 技术栈

- **语言**: Go 1.24.11
- **HTTP 框架**: Hertz (CloudWeGo)
- **RPC 框架**: Kitex (CloudWeGo)
- **ORM**: GORM
- **数据库**: PostgreSQL 16
- **缓存**: Redis 7
- **依赖管理**: BOM (Bill of Materials)
- **容器化**: Docker + Docker Compose

## 快速开始

### 1. 同步依赖

```bash
go work sync
```

### 2. 启动数据库服务

```bash
docker-compose up -d postgres redis
```

### 3. 运行应用

```bash
go run ./cmd/api/main.go
```

访问 http://localhost:8080/health 检查服务状态。

## 项目结构

```
soliton-client/
├── go.work                   # Go 工作区配置
├── bom/                      # BOM 依赖管理模块
├── share/                    # 公共组件模块
│   ├── errors/               # 错误定义
│   ├── utils/                # 工具函数
│   ├── types/                # 通用类型
│   └── middleware/           # 中间件
├── user/                     # 用户聚合模块
│   ├── domain/               # 领域层
│   │   ├── entity/           # 领域实体
│   │   ├── repository/       # 仓储接口
│   │   ├── service/          # 领域服务
│   │   ├── valueobject/      # 值对象
│   │   └── event/            # 领域事件
│   └── infrastructure/       # 基础设施层
│       ├── entity/           # 数据库实体 (PO)
│       ├── converter/        # 转换器
│       └── repository/       # 仓储实现
├── api/                      # API 聚合模块
│   └── user-api/             # 用户 API
│       ├── dto/              # 数据传输对象
│       ├── service/          # 应用服务
│       └── http/             # HTTP 处理器
└── cmd/
    └── api/                  # 主程序入口
```

## 环境变量

- `DB_HOST`: PostgreSQL 主机（默认：localhost）
- `DB_PORT`: PostgreSQL 端口（默认：5432）
- `DB_USER`: 数据库用户（默认：postgres）
- `DB_PASSWORD`: 数据库密码（默认：postgres）
- `DB_NAME`: 数据库名称（默认：soliton-client）
- `REDIS_HOST`: Redis 主机（默认：localhost）
- `REDIS_PORT`: Redis 端口（默认：6379）

## 常用命令

```bash
# 构建
make build

# 运行
make run

# 测试
make test

# 同步依赖
make tidy

# 启动 Docker 服务
make docker-up

# 停止 Docker 服务
make docker-down
```

## License

MIT
