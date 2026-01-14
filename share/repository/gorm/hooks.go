package gorm

import (
	"time"

	"gorm.io/gorm"
)

// BeforeCreate GORM 创建前钩子
// 自动设置创建时间、更新时间和版本号
func (e *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = now
	}
	if e.Version == 0 {
		e.Version = 1
	}
	return nil
}

// BeforeUpdate GORM 更新前钩子
// 自动更新更新时间和版本号
func (e *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	e.Version++
	return nil
}

// RegisterAuditCallbacks 注册审计回调到 GORM
// 为所有实现 Auditable 接口的实体自动填充审计字段
func RegisterAuditCallbacks(db *gorm.DB) {
	// 创建前回调
	db.Callback().Create().Before("gorm:create").Register("audit:before_create", func(tx *gorm.DB) {
		if tx.Statement.Schema == nil {
			return
		}

		now := time.Now()

		// 设置创建时间
		if field := tx.Statement.Schema.LookUpField("CreatedAt"); field != nil {
			if _, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue); isZero {
				_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, now)
			}
		}

		// 设置更新时间
		if field := tx.Statement.Schema.LookUpField("UpdatedAt"); field != nil {
			if _, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue); isZero {
				_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, now)
			}
		}

		// 设置版本号
		if field := tx.Statement.Schema.LookUpField("Version"); field != nil {
			if val, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue); isZero || val == 0 {
				_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, 1)
			}
		}
	})

	// 更新前回调
	db.Callback().Update().Before("gorm:update").Register("audit:before_update", func(tx *gorm.DB) {
		if tx.Statement.Schema == nil {
			return
		}

		// 设置更新时间
		if field := tx.Statement.Schema.LookUpField("UpdatedAt"); field != nil {
			_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, time.Now())
		}

		// 版本号递增（乐观锁）
		if field := tx.Statement.Schema.LookUpField("Version"); field != nil {
			if val, _ := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue); val != nil {
				if version, ok := val.(int); ok {
					_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, version+1)
				}
			}
		}
	})
}
