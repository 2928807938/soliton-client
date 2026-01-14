package repository

import (
	"context"
)

// QueryBuilder 查询构建器接口，提供链式调用的查询构建能力
type QueryBuilder[T any] interface {
	// Where 添加查询条件
	Where(condition *Condition) QueryBuilder[T]

	// And 添加 AND 条件
	And(conditions ...*Condition) QueryBuilder[T]

	// OrderBy 添加排序（升序）
	OrderBy(field string) QueryBuilder[T]

	// OrderByDesc 添加排序（降序）
	OrderByDesc(field string) QueryBuilder[T]

	// Limit 限制返回数量
	Limit(limit int) QueryBuilder[T]

	// Offset 设置偏移量
	Offset(offset int) QueryBuilder[T]

	// Select 指定查询字段
	Select(fields ...string) QueryBuilder[T]

	// Find 执行查询，返回结果列表
	Find(ctx context.Context) ([]*T, error)

	// First 执行查询，返回第一条结果
	First(ctx context.Context) (*T, error)

	// Count 执行统计查询
	Count(ctx context.Context) (int64, error)

	// Exists 执行存在性检查
	Exists(ctx context.Context) (bool, error)
}

// QueryOptions 查询选项，用于存储构建器的状态
type QueryOptions struct {
	Conditions []*Condition // 查询条件列表
	OrderBys   []OrderBy    // 排序规则列表
	LimitVal   int          // 限制数量
	OffsetVal  int          // 偏移量
	Fields     []string     // 查询字段
}

// NewQueryOptions 创建查询选项
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Conditions: make([]*Condition, 0),
		OrderBys:   make([]OrderBy, 0),
		LimitVal:   0,
		OffsetVal:  0,
		Fields:     make([]string, 0),
	}
}

// AddCondition 添加条件
func (o *QueryOptions) AddCondition(condition *Condition) *QueryOptions {
	o.Conditions = append(o.Conditions, condition)
	return o
}

// AddConditions 批量添加条件
func (o *QueryOptions) AddConditions(conditions ...*Condition) *QueryOptions {
	o.Conditions = append(o.Conditions, conditions...)
	return o
}

// AddOrderBy 添加排序
func (o *QueryOptions) AddOrderBy(field string, desc bool) *QueryOptions {
	o.OrderBys = append(o.OrderBys, OrderBy{Field: field, Desc: desc})
	return o
}

// SetLimit 设置限制
func (o *QueryOptions) SetLimit(limit int) *QueryOptions {
	o.LimitVal = limit
	return o
}

// SetOffset 设置偏移
func (o *QueryOptions) SetOffset(offset int) *QueryOptions {
	o.OffsetVal = offset
	return o
}

// SetFields 设置查询字段
func (o *QueryOptions) SetFields(fields ...string) *QueryOptions {
	o.Fields = fields
	return o
}
