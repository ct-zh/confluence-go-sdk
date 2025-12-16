package confluence

import (
	"net/http"
	"strings"
)

// Option 定义配置 Client 的函数签名
type Option func(*Client)

// WithHTTPClient 允许注入自定义的 HTTP Client
// 场景：设置超时、代理、或使用 mock client 进行测试
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithAuth 设置基础认证信息 (Basic Auth)
// email: 用户邮箱
// token: API Token
func WithAuth(email, token string) Option {
	return func(c *Client) {
		c.auth = &basicAuth{
			email: email,
			token: token,
		}
	}
}

// WithToken 设置 Bearer Token 认证 (OAuth / PAT)
func WithToken(token string) Option {
	return func(c *Client) {
		c.auth = &tokenAuth{
			token: token,
		}
	}
}

// authenticator 接口定义认证行为
type authenticator interface {
	SetAuth(req *http.Request)
}

// basicAuth 实现 Basic Authentication
type basicAuth struct {
	email string
	token string
}

func (b *basicAuth) SetAuth(req *http.Request) {
	req.SetBasicAuth(b.email, b.token)
}

// tokenAuth 实现 Bearer Token Authentication
type tokenAuth struct {
	token string
}

func (t *tokenAuth) SetAuth(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+t.token)
}

// normalizeBaseURL 确保 BaseURL 格式正确（移除尾部斜杠）
func normalizeBaseURL(url string) string {
	return strings.TrimRight(url, "/")
}
