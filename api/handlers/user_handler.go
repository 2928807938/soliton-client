package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"soliton-client/share/types"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(router *route.RouterGroup, userClient *UserServiceClient) {
	users := router.Group("/users")
	{
		users.POST("/register", handleRegister(userClient))            // 用户注册
		users.POST("/login", handleLogin(userClient))                  // 用户登录
		users.POST("/logout", handleLogout(userClient))                // 用户登出
		users.POST("/refresh", handleRefreshToken(userClient))         // 刷新 Token
		users.GET("/:id", handleGetUser(userClient))                   // 获取用户信息
		users.PUT("/:id", handleUpdateUser(userClient))                // 更新用户信息
		users.POST("/password/reset", handleResetPassword(userClient)) // 重置密码
	}

	// 验证码相关路由
	verification := router.Group("/verification")
	{
		verification.POST("/send", handleSendVerificationCode(userClient)) // 发送验证码
		verification.POST("/verify", handleVerifyCode(userClient))         // 验证验证码
	}
}

// doRequest 发送HTTP请求到用户服务
func doRequest(c *UserServiceClient, method, path string, body interface{}) (*types.Response, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("序列化请求失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if c.TenantID != "" {
		req.Header.Set("X-Tenant-Id", c.TenantID)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("请求用户服务失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var apiResp types.Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("解析响应失败: %v, body: %s", err, string(respBody))
	}

	return &apiResp, resp.StatusCode, nil
}

// handleRegister 用户注册
func handleRegister(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "POST", "/api/v1/users/register", req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "注册成功",
			"user":    resp.Data,
		})
	}
}

// handleLogin 用户登录
func handleLogin(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "POST", "/api/v1/auth/login", req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, resp.Data)
	}
}

// handleLogout 用户登出
func handleLogout(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头获取 Authorization token
		token := string(c.GetHeader("Authorization"))

		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			req = make(map[string]interface{})
		}

		// 创建带 token 的请求
		reqBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", userClient.BaseURL+"/api/v1/auth/logout", bytes.NewBuffer(reqBody))
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"error": "创建请求失败",
			})
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		if userClient.TenantID != "" {
			httpReq.Header.Set("X-Tenant-Id", userClient.TenantID)
		}
		if token != "" {
			httpReq.Header.Set("Authorization", token)
		}

		httpResp, err := userClient.HTTPClient.Do(httpReq)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": "请求用户服务失败",
			})
			return
		}
		defer httpResp.Body.Close()

		respBody, _ := io.ReadAll(httpResp.Body)
		var apiResp types.Response
		json.Unmarshal(respBody, &apiResp)

		if apiResp.Code != 0 {
			c.JSON(httpResp.StatusCode, map[string]interface{}{
				"error": apiResp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "登出成功",
		})
	}
}

// handleRefreshToken 刷新 Token
func handleRefreshToken(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "POST", "/api/v1/auth/refresh", req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, resp.Data)
	}
}

// handleGetUser 获取用户信息
func handleGetUser(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := c.Param("id")
		if id == "" {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "用户 ID 不能为空",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "GET", "/api/v1/users/"+id, nil)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, resp.Data)
	}
}

// handleUpdateUser 更新用户信息
func handleUpdateUser(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := c.Param("id")
		if id == "" {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "用户 ID 不能为空",
			})
			return
		}

		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "PUT", "/api/v1/users/"+id, req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "更新成功",
			"user":    resp.Data,
		})
	}
}

// handleResetPassword 重置密码
func handleResetPassword(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "POST", "/api/v1/users/password/reset", req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "密码重置成功",
		})
	}
}

// handleSendVerificationCode 发送验证码
func handleSendVerificationCode(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "POST", "/api/v1/verification/code/send", req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "验证码已发送",
		})
	}
}

// handleVerifyCode 验证验证码
func handleVerifyCode(userClient *UserServiceClient) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, statusCode, err := doRequest(userClient, "POST", "/api/v1/verification/code/verify", req)
		if err != nil {
			c.JSON(consts.StatusServiceUnavailable, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if resp.Code != 0 {
			c.JSON(statusCode, map[string]interface{}{
				"error": resp.Message,
			})
			return
		}

		// 检查验证结果
		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			c.JSON(consts.StatusOK, map[string]interface{}{
				"message": "验证成功",
			})
			return
		}

		valid, _ := data["valid"].(bool)
		if !valid {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "验证码无效或已过期",
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "验证成功",
		})
	}
}
