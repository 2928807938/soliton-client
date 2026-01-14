package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"soliton-client/share/repository"
)

// 事务上下文键
type txKey struct{}

// GormRepository 基于 GORM 的通用仓储实现
type GormRepository[T any, ID comparable] struct {
	db *gorm.DB
}

// NewGormRepository 创建 GORM 仓储实例
func NewGormRepository[T any, ID comparable](db *gorm.DB) *GormRepository[T, ID] {
	return &GormRepository[T, ID]{
		db: db,
	}
}

// DB 获取底层 GORM DB 实例
func (r *GormRepository[T, ID]) DB() *gorm.DB {
	return r.db
}

// getDB 获取数据库连接（支持事务）
func (r *GormRepository[T, ID]) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

// Create 创建单个实体
func (r *GormRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	return r.getDB(ctx).Create(entity).Error
}

// CreateBatch 批量创建实体
func (r *GormRepository[T, ID]) CreateBatch(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.getDB(ctx).Create(entities).Error
}

// GetByID 根据主键查询
func (r *GormRepository[T, ID]) GetByID(ctx context.Context, id ID) (*T, error) {
	var entity T
	err := r.getDB(ctx).First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// Update 更新实体
func (r *GormRepository[T, ID]) Update(ctx context.Context, entity *T) error {
	return r.getDB(ctx).Save(entity).Error
}

// Delete 删除实体（逻辑删除）
func (r *GormRepository[T, ID]) Delete(ctx context.Context, id ID) error {
	var entity T
	return r.getDB(ctx).Delete(&entity, id).Error
}

// List 查询全部列表
func (r *GormRepository[T, ID]) List(ctx context.Context) ([]*T, error) {
	var entities []*T
	err := r.getDB(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Page 分页查询
func (r *GormRepository[T, ID]) Page(ctx context.Context, request *repository.PageRequest) (*repository.PageResult[*T], error) {
	db := r.getDB(ctx)

	// 应用查询条件
	if len(request.Conditions) > 0 {
		db = ApplyConditions(db, request.Conditions...)
	}

	// 统计总数
	var total int64
	var entity T
	if err := db.Model(&entity).Count(&total).Error; err != nil {
		return nil, err
	}

	// 应用排序
	for _, order := range request.OrderBy {
		if order.Desc {
			db = db.Order(order.Field + " DESC")
		} else {
			db = db.Order(order.Field + " ASC")
		}
	}

	// 应用分页
	db = db.Offset(request.Offset()).Limit(request.Size)

	// 查询数据
	var entities []*T
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}

	return repository.NewPageResult(entities, total, request.Page, request.Size), nil
}

// BeginTx 开启事务
func (r *GormRepository[T, ID]) BeginTx(ctx context.Context) (context.Context, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}
	return context.WithValue(ctx, txKey{}, tx), nil
}

// Commit 提交事务
func (r *GormRepository[T, ID]) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx.Commit().Error
	}
	return errors.New("no transaction in context")
}

// Rollback 回滚事务
func (r *GormRepository[T, ID]) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx.Rollback().Error
	}
	return errors.New("no transaction in context")
}

// WithTx 在事务中执行操作
func (r *GormRepository[T, ID]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	if err := fn(txCtx); err != nil {
		_ = r.Rollback(txCtx)
		return err
	}

	return r.Commit(txCtx)
}

// 确保实现了接口
var _ repository.BaseRepository[any, int] = (*GormRepository[any, int])(nil)
var _ repository.TransactionalRepository = (*GormRepository[any, int])(nil)
