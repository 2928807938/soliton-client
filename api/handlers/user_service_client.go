package handlers

import (
	"net/http"
	"time"
)

// UserServiceClient HTTP 客户端
type UserServiceClient struct {
	BaseURL    string
	TenantID   string
	HTTPClient *http.Client
}

// NewUserServiceClient 创建用户服务客户端
func NewUserServiceClient(baseURL, tenantID string) *UserServiceClient {
	return &UserServiceClient{
		BaseURL:  baseURL,
		TenantID: tenantID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
