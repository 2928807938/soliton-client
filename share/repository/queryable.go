package repository

import "context"

// Operator 操作符类型
type Operator string

const (
	// 比较操作符
	OpEqual          Operator = "="
	OpNotEqual       Operator = "!="
	OpGreaterThan    Operator = ">"
	OpGreaterOrEqual Operator = ">="
	OpLessThan       Operator = "<"
	OpLessOrEqual    Operator = "<="

	// 模糊匹配
	OpLike Operator = "LIKE"

	// 集合操作
	OpIn    Operator = "IN"
	OpNotIn Operator = "NOT IN"

	// 区间操作
	OpBetween Operator = "BETWEEN"

	// 空值检查
	OpIsNull    Operator = "IS NULL"
	OpIsNotNull Operator = "IS NOT NULL"
)

// Condition 查询条件
type Condition struct {
	Field    string      // 要查询的字段名
	Operator Operator    // 操作符
	Value    interface{} // 比较的值
}

// NewCondition 创建查询条件
func NewCondition(field string, operator Operator, value interface{}) *Condition {
	return &Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// Eq 等于条件
func Eq(field string, value interface{}) *Condition {
	return NewCondition(field, OpEqual, value)
}

// NotEq 不等于条件
func NotEq(field string, value interface{}) *Condition {
	return NewCondition(field, OpNotEqual, value)
}

// Gt 大于条件
func Gt(field string, value interface{}) *Condition {
	return NewCondition(field, OpGreaterThan, value)
}

// Gte 大于等于条件
func Gte(field string, value interface{}) *Condition {
	return NewCondition(field, OpGreaterOrEqual, value)
}

// Lt 小于条件
func Lt(field string, value interface{}) *Condition {
	return NewCondition(field, OpLessThan, value)
}

// Lte 小于等于条件
func Lte(field string, value interface{}) *Condition {
	return NewCondition(field, OpLessOrEqual, value)
}

// Like 模糊匹配条件
func Like(field string, value string) *Condition {
	return NewCondition(field, OpLike, value)
}

// In 包含条件
func In(field string, values interface{}) *Condition {
	return NewCondition(field, OpIn, values)
}

// NotIn 不包含条件
func NotIn(field string, values interface{}) *Condition {
	return NewCondition(field, OpNotIn, values)
}

// Between 区间条件
func Between(field string, start, end interface{}) *Condition {
	return NewCondition(field, OpBetween, []interface{}{start, end})
}

// IsNull 为空条件
func IsNull(field string) *Condition {
	return NewCondition(field, OpIsNull, nil)
}

// IsNotNull 不为空条件
func IsNotNull(field string) *Condition {
	return NewCondition(field, OpIsNotNull, nil)
}

// QueryableRepository 可查询仓储接口，提供条件查询能力
type QueryableRepository[T any, ID comparable] interface {
	BaseRepository[T, ID]
	// Where 条件查询
	Where(ctx context.Context, conditions ...*Condition) ([]*T, error)

	// Count 统计数量
	Count(ctx context.Context, conditions ...*Condition) (int64, error)

	// Exists 存在性检查
	Exists(ctx context.Context, conditions ...*Condition) (bool, error)

	// Query 获取查询构建器
	Query() QueryBuilder[T]
}
