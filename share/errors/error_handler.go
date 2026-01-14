package errors

import (
	"context"
	"errors"
	"net/http"
	"soliton-client/share/types"

	"github.com/cloudwego/hertz/pkg/app"
)

// HandleError 统一错误处理
// 支持处理 AppError 及其继承类型（如 UserError）
func HandleError(ctx context.Context, c *app.RequestContext, err error) {
	// 使用 errors.As 支持嵌入类型的解包
	var appErr *AppError
	if errors.As(err, &appErr) {
		status := getHTTPStatus(appErr.Code)
		c.JSON(status, types.Error(appErr.Code, appErr.Message))
		return
	}

	c.JSON(http.StatusInternalServerError, types.Error(InternalError, "内部服务错误"))
}

// getHTTPStatus 根据业务错误码获取对应的 HTTP 状态码
// 错误码分段规则:
//
//	10000-10999: 通用错误
//	11000-11999: User 模块
//	12000-12999: Order 模块
//	...以此类推
func getHTTPStatus(code int) int {
	// 根据错误码末尾判断类型
	switch code % 100 {
	case 1: // xxx01: bad_request
		return http.StatusBadRequest
	case 2: // xxx02: unauthorized
		return http.StatusUnauthorized
	case 3: // xxx03: forbidden
		return http.StatusForbidden
	case 4: // xxx04: not_found
		return http.StatusNotFound
	case 5: // xxx05: conflict
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// IsAppError 判断是否为 AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError 将 error 转换为 AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
