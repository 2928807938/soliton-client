package main

import (
	"context"
	"fmt"

	usersdk "github.com/2928807938/universal-service-user/sdk"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

// registerUserRoutes 注册用户相关路由
func registerUserRoutes(router *route.RouterGroup, client *usersdk.Client) {
	users := router.Group("/users")
	{
		users.POST("/register", handleRegister(client))            // 用户注册
		users.POST("/login", handleLogin(client))                  // 用户登录
		users.POST("/logout", handleLogout(client))                // 用户登出
		users.POST("/refresh", handleRefreshToken(client))         // 刷新 Token
		users.GET("/:id", handleGetUser(client))                   // 获取用户信息
		users.PUT("/:id", handleUpdateUser(client))                // 更新用户信息
		users.POST("/password/reset", handleResetPassword(client)) // 重置密码
	}

	// 验证码相关路由
	verification := router.Group("/verification")
	{
		verification.POST("/send", handleSendVerificationCode(client)) // 发送验证码
		verification.POST("/verify", handleVerifyCode(client))         // 验证验证码
	}
}

// handleRegister 用户注册
func handleRegister(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req usersdk.RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		user, err := client.Register(ctx, &req)
		if err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "注册成功",
			"user":    user,
		})
	}
}

// handleLogin 用户登录
func handleLogin(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req usersdk.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		resp, err := client.Login(ctx, &req)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, resp)
	}
}

// handleLogout 用户登出
func handleLogout(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		if err := client.Logout(ctx, req.AccessToken, req.RefreshToken); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "登出成功",
		})
	}
}

// handleRefreshToken 刷新 Token
func handleRefreshToken(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		accessToken, refreshToken, expiresIn, err := client.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    expiresIn,
		})
	}
}

// handleGetUser 获取用户信息
func handleGetUser(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := c.Param("id")
		if id == "" {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "用户 ID 不能为空",
			})
			return
		}

		// 将 string 转换为 int
		var userID int
		if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "无效的用户 ID",
			})
			return
		}

		user, err := client.GetUser(ctx, userID)
		if err != nil {
			c.JSON(consts.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, user)
	}
}

// handleUpdateUser 更新用户信息
func handleUpdateUser(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := c.Param("id")
		if id == "" {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "用户 ID 不能为空",
			})
			return
		}

		var userID int
		if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "无效的用户 ID",
			})
			return
		}

		var req usersdk.UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		user, err := client.UpdateUser(ctx, userID, &req)
		if err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "更新成功",
			"user":    user,
		})
	}
}

// handleResetPassword 重置密码
func handleResetPassword(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req usersdk.ResetPasswordRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		if err := client.ResetPassword(ctx, &req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "密码重置成功",
		})
	}
}

// handleSendVerificationCode 发送验证码
func handleSendVerificationCode(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req usersdk.SendVerificationCodeRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		if err := client.SendVerificationCode(ctx, &req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, map[string]interface{}{
			"message": "验证码已发送",
		})
	}
}

// handleVerifyCode 验证验证码
func handleVerifyCode(client *usersdk.Client) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req usersdk.VerifyCodeRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "参数错误",
			})
			return
		}

		valid, err := client.VerifyCode(ctx, &req)
		if err != nil {
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

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
