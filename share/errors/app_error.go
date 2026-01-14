package errors

import "fmt"

// AppError 应用错误基类
type AppError struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
	Err     error  `json:"-"`       // 原始错误
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装原始错误
func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ==================== 通用错误 ====================
// 错误码分段: 1xxxx
// 示例: 10001, 10002...

const (
	// 通用错误码 10000-10999
	Success       = 200   // 成功
	BadRequest    = 10001 // 请求参数错误
	Unauthorized  = 10002 // 未授权
	Forbidden     = 10003 // 禁止访问
	NotFound      = 10004 // 资源不存在
	Conflict      = 10005 // 资源冲突
	InternalError = 10006 // 内部错误
)

// ErrBadRequest 请求参数错误
func ErrBadRequest(message string) *AppError {
	return New(BadRequest, message)
}

// ErrNotFound 资源不存在
func ErrNotFound(message string) *AppError {
	return New(NotFound, message)
}

// ErrUnauthorized 未授权
func ErrUnauthorized(message string) *AppError {
	return New(Unauthorized, message)
}

// ErrForbidden 禁止访问
func ErrForbidden(message string) *AppError {
	return New(Forbidden, message)
}

// ErrConflict 资源冲突
func ErrConflict(message string) *AppError {
	return New(Conflict, message)
}

// ErrInternal 内部错误
func ErrInternal(message string, err error) *AppError {
	return Wrap(InternalError, message, err)
}
