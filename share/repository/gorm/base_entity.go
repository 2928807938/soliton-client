package gorm

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity 基础实体，包含通用的审计字段
// 业务实体通过组合方式继承这些字段
type BaseEntity struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Version   int            `gorm:"default:1" json:"version"`
}

// GetID 获取实体主键
func (e *BaseEntity) GetID() int {
	return e.ID
}

// IsDeleted 判断是否已删除
func (e *BaseEntity) IsDeleted() bool {
	return e.DeletedAt.Valid
}

// SetCreatedAt 设置创建时间
func (e *BaseEntity) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// SetUpdatedAt 设置更新时间
func (e *BaseEntity) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = t
}

// IncrementVersion 版本号递增
func (e *BaseEntity) IncrementVersion() {
	e.Version++
}

// GetVersion 获取版本号
func (e *BaseEntity) GetVersion() int {
	return e.Version
}

// AuditFields 审计字段，不包含 ID，可供自定义主键类型的实体组合使用
type AuditFields struct {
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Version   int            `gorm:"default:1" json:"version"`
}

// IsDeleted 判断是否已删除
func (e *AuditFields) IsDeleted() bool {
	return e.DeletedAt.Valid
}

// SetCreatedAt 设置创建时间
func (e *AuditFields) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// SetUpdatedAt 设置更新时间
func (e *AuditFields) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = t
}

// IncrementVersion 版本号递增
func (e *AuditFields) IncrementVersion() {
	e.Version++
}

// GetVersion 获取版本号
func (e *AuditFields) GetVersion() int {
	return e.Version
}

// Touch 更新修改时间
func (e *AuditFields) Touch() {
	e.UpdatedAt = time.Now()
}

// Auditable 可审计接口，实现此接口的实体将自动填充审计字段
type Auditable interface {
	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
	IncrementVersion()
	GetVersion() int
}
