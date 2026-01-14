package repository

// PageRequest 分页请求
type PageRequest struct {
	Page       int          `json:"page"`       // 页码（从 1 开始）
	Size       int          `json:"size"`       // 每页数量
	Conditions []*Condition `json:"conditions"` // 查询条件列表
	OrderBy    []OrderBy    `json:"order_by"`   // 排序规则
}

// OrderBy 排序规则
type OrderBy struct {
	Field string `json:"field"` // 排序字段
	Desc  bool   `json:"desc"`  // 是否降序
}

// NewPageRequest 创建分页请求
func NewPageRequest(page, size int) *PageRequest {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	return &PageRequest{
		Page:       page,
		Size:       size,
		Conditions: make([]*Condition, 0),
		OrderBy:    make([]OrderBy, 0),
	}
}

// WithCondition 添加查询条件
func (p *PageRequest) WithCondition(condition *Condition) *PageRequest {
	p.Conditions = append(p.Conditions, condition)
	return p
}

// WithOrderBy 添加排序规则
func (p *PageRequest) WithOrderBy(field string, desc bool) *PageRequest {
	p.OrderBy = append(p.OrderBy, OrderBy{Field: field, Desc: desc})
	return p
}

// Offset 计算偏移量
func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Size
}

// PageResult 分页结果
type PageResult[T any] struct {
	Items      []T   `json:"items"`       // 当前页数据列表
	Total      int64 `json:"total"`       // 总记录数
	Page       int   `json:"page"`        // 当前页码
	Size       int   `json:"size"`        // 每页数量
	TotalPages int   `json:"total_pages"` // 总页数
}

// NewPageResult 创建分页结果
func NewPageResult[T any](items []T, total int64, page, size int) *PageResult[T] {
	totalPages := int(total) / size
	if int(total)%size != 0 {
		totalPages++
	}
	return &PageResult[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}
}

// HasNext 是否有下一页
func (p *PageResult[T]) HasNext() bool {
	return p.Page < p.TotalPages
}

// HasPrev 是否有上一页
func (p *PageResult[T]) HasPrev() bool {
	return p.Page > 1
}

// IsEmpty 是否为空
func (p *PageResult[T]) IsEmpty() bool {
	return len(p.Items) == 0
}
