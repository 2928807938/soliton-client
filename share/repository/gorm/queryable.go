package gorm

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"soliton-client/share/repository"
)

// QueryableGormRepository 可查询的 GORM 仓储实现
type QueryableGormRepository[T any, ID comparable] struct {
	*GormRepository[T, ID]
}

// NewQueryableGormRepository 创建可查询的 GORM 仓储实例
func NewQueryableGormRepository[T any, ID comparable](db *gorm.DB) *QueryableGormRepository[T, ID] {
	return &QueryableGormRepository[T, ID]{
		GormRepository: NewGormRepository[T, ID](db),
	}
}

// Where 条件查询
func (r *QueryableGormRepository[T, ID]) Where(ctx context.Context, conditions ...*repository.Condition) ([]*T, error) {
	db := ApplyConditions(r.getDB(ctx), conditions...)
	var entities []*T
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count 统计数量
func (r *QueryableGormRepository[T, ID]) Count(ctx context.Context, conditions ...*repository.Condition) (int64, error) {
	db := ApplyConditions(r.getDB(ctx), conditions...)
	var count int64
	var entity T
	if err := db.Model(&entity).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists 存在性检查
func (r *QueryableGormRepository[T, ID]) Exists(ctx context.Context, conditions ...*repository.Condition) (bool, error) {
	count, err := r.Count(ctx, conditions...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Query 获取查询构建器
func (r *QueryableGormRepository[T, ID]) Query() repository.QueryBuilder[T] {
	return NewGormQueryBuilder[T](r.GormRepository.DB())
}

// ApplyConditions 将条件列表应用到 GORM 查询（包级函数）
func ApplyConditions(db *gorm.DB, conditions ...*repository.Condition) *gorm.DB {
	for _, cond := range conditions {
		db = ApplyCondition(db, cond)
	}
	return db
}

// ApplyCondition 应用单个条件到 GORM 查询（包级函数）
func ApplyCondition(db *gorm.DB, cond *repository.Condition) *gorm.DB {
	switch cond.Operator {
	case repository.OpEqual:
		return db.Where(fmt.Sprintf("%s = ?", cond.Field), cond.Value)
	case repository.OpNotEqual:
		return db.Where(fmt.Sprintf("%s != ?", cond.Field), cond.Value)
	case repository.OpGreaterThan:
		return db.Where(fmt.Sprintf("%s > ?", cond.Field), cond.Value)
	case repository.OpGreaterOrEqual:
		return db.Where(fmt.Sprintf("%s >= ?", cond.Field), cond.Value)
	case repository.OpLessThan:
		return db.Where(fmt.Sprintf("%s < ?", cond.Field), cond.Value)
	case repository.OpLessOrEqual:
		return db.Where(fmt.Sprintf("%s <= ?", cond.Field), cond.Value)
	case repository.OpLike:
		return db.Where(fmt.Sprintf("%s LIKE ?", cond.Field), cond.Value)
	case repository.OpIn:
		return db.Where(fmt.Sprintf("%s IN ?", cond.Field), cond.Value)
	case repository.OpNotIn:
		return db.Where(fmt.Sprintf("%s NOT IN ?", cond.Field), cond.Value)
	case repository.OpBetween:
		if values, ok := cond.Value.([]interface{}); ok && len(values) == 2 {
			return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", cond.Field), values[0], values[1])
		}
		return db
	case repository.OpIsNull:
		return db.Where(fmt.Sprintf("%s IS NULL", cond.Field))
	case repository.OpIsNotNull:
		return db.Where(fmt.Sprintf("%s IS NOT NULL", cond.Field))
	default:
		return db
	}
}

// GormQueryBuilder GORM 查询构建器实现
type GormQueryBuilder[T any] struct {
	db      *gorm.DB
	options *repository.QueryOptions
}

// NewGormQueryBuilder 创建 GORM 查询构建器
func NewGormQueryBuilder[T any](db *gorm.DB) *GormQueryBuilder[T] {
	return &GormQueryBuilder[T]{
		db:      db,
		options: repository.NewQueryOptions(),
	}
}

// Where 添加查询条件
func (b *GormQueryBuilder[T]) Where(condition *repository.Condition) repository.QueryBuilder[T] {
	b.options.AddCondition(condition)
	return b
}

// And 添加 AND 条件
func (b *GormQueryBuilder[T]) And(conditions ...*repository.Condition) repository.QueryBuilder[T] {
	b.options.AddConditions(conditions...)
	return b
}

// OrderBy 添加排序（升序）
func (b *GormQueryBuilder[T]) OrderBy(field string) repository.QueryBuilder[T] {
	b.options.AddOrderBy(field, false)
	return b
}

// OrderByDesc 添加排序（降序）
func (b *GormQueryBuilder[T]) OrderByDesc(field string) repository.QueryBuilder[T] {
	b.options.AddOrderBy(field, true)
	return b
}

// Limit 限制返回数量
func (b *GormQueryBuilder[T]) Limit(limit int) repository.QueryBuilder[T] {
	b.options.SetLimit(limit)
	return b
}

// Offset 设置偏移量
func (b *GormQueryBuilder[T]) Offset(offset int) repository.QueryBuilder[T] {
	b.options.SetOffset(offset)
	return b
}

// Select 指定查询字段
func (b *GormQueryBuilder[T]) Select(fields ...string) repository.QueryBuilder[T] {
	b.options.SetFields(fields...)
	return b
}

// build 构建 GORM 查询
func (b *GormQueryBuilder[T]) build(ctx context.Context) *gorm.DB {
	db := b.db.WithContext(ctx)

	// 应用查询条件
	for _, cond := range b.options.Conditions {
		db = ApplyCondition(db, cond)
	}

	// 应用字段选择
	if len(b.options.Fields) > 0 {
		db = db.Select(b.options.Fields)
	}

	// 应用排序
	for _, order := range b.options.OrderBys {
		if order.Desc {
			db = db.Order(order.Field + " DESC")
		} else {
			db = db.Order(order.Field + " ASC")
		}
	}

	// 应用分页
	if b.options.LimitVal > 0 {
		db = db.Limit(b.options.LimitVal)
	}
	if b.options.OffsetVal > 0 {
		db = db.Offset(b.options.OffsetVal)
	}

	return db
}

// Find 执行查询，返回结果列表
func (b *GormQueryBuilder[T]) Find(ctx context.Context) ([]*T, error) {
	var entities []*T
	if err := b.build(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// First 执行查询，返回第一条结果
func (b *GormQueryBuilder[T]) First(ctx context.Context) (*T, error) {
	var entity T
	if err := b.build(ctx).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// Count 执行统计查询
func (b *GormQueryBuilder[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	db := b.db.WithContext(ctx)

	// 只应用查询条件
	for _, cond := range b.options.Conditions {
		db = ApplyCondition(db, cond)
	}

	if err := db.Model(&entity).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists 执行存在性检查
func (b *GormQueryBuilder[T]) Exists(ctx context.Context) (bool, error) {
	count, err := b.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Page 执行分页查询
func (b *GormQueryBuilder[T]) Page(ctx context.Context, page, size int) (*repository.PageResult[*T], error) {
	// 统计总数
	total, err := b.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 设置分页参数
	b.options.SetOffset((page - 1) * size)
	b.options.SetLimit(size)

	// 查询数据
	entities, err := b.Find(ctx)
	if err != nil {
		return nil, err
	}

	return repository.NewPageResult(entities, total, page, size), nil
}

// 确保实现了接口
var _ repository.QueryableRepository[any, int] = (*QueryableGormRepository[any, int])(nil)
var _ repository.QueryBuilder[any] = (*GormQueryBuilder[any])(nil)
