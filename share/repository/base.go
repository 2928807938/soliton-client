package repository

import (
	"context"
)

// BaseRepository 基础仓储接口，定义通用的 CRUD 操作
// T 为实体类型，ID 为主键类型
type BaseRepository[T any, ID comparable] interface {
	// Create 创建单个实体
	Create(ctx context.Context, entity *T) error

	// CreateBatch 批量创建实体
	CreateBatch(ctx context.Context, entities []*T) error

	// GetByID 根据主键查询
	GetByID(ctx context.Context, id ID) (*T, error)

	// Update 更新实体
	Update(ctx context.Context, entity *T) error

	// Delete 删除实体（逻辑删除）
	Delete(ctx context.Context, id ID) error

	// List 查询全部列表
	List(ctx context.Context) ([]*T, error)

	// Page 分页查询
	Page(ctx context.Context, request *PageRequest) (*PageResult[*T], error)
}

// TransactionalRepository 支持事务的仓储接口
type TransactionalRepository interface {
	// BeginTx 开启事务
	BeginTx(ctx context.Context) (context.Context, error)

	// Commit 提交事务
	Commit(ctx context.Context) error

	// Rollback 回滚事务
	Rollback(ctx context.Context) error
}

// Entity 实体接口，所有实体必须实现此接口
type Entity[ID comparable] interface {
	// GetID 获取实体主键
	GetID() ID
}
