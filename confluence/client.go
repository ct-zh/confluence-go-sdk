package confluence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultUserAgent = "confluence-go-sdk/1.0"
)

// Client Confluence API 客户端
// 包含通用配置和各个服务模块的访问入口
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string
	auth       authenticator

	// Services
	// 将在这里挂载各个业务模块
	Content IContentService
	Search  ISearchService
	Space   ISpaceService
}

// NewClient 创建一个新的 Confluence 客户端
// baseURL: Confluence 实例地址，例如 "https://your-domain.atlassian.net/wiki"
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	parsedURL, err := url.Parse(normalizeBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	c := &Client{
		httpClient: http.DefaultClient,
		baseURL:    parsedURL,
		userAgent:  defaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	// 初始化 Services
	c.Content = &ContentService{client: c}
	c.Search = &SearchService{client: c}
	c.Space = &SpaceService{client: c}

	return c, nil
}

// NewRequest 创建一个 API 请求
// method: HTTP 方法 (GET, POST, PUT, DELETE)
// urlStr: 相对路径，例如 "/rest/api/content"
// body: 请求体。如果是 io.Reader，直接使用；如果是 struct，会自动序列化为 JSON
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// 拼接完整 URL
	u := c.baseURL.ResolveReference(rel)

	var buf io.Reader
	if body != nil {
		if r, ok := body.(io.Reader); ok {
			buf = r
		} else {
			b := new(bytes.Buffer)
			err := json.NewEncoder(b).Encode(body)
			if err != nil {
				return nil, err
			}
			buf = b
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		// 只有当 body 不是 io.Reader 时（即我们执行了 JSON 序列化），才强制设置 Content-Type 为 application/json
		// 如果是 io.Reader，调用者需要自己设置 Content-Type（例如 multipart/form-data）
		if _, ok := body.(io.Reader); !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.auth != nil {
		c.auth.SetAuth(req)
	}

	return req, nil
}

// Do 发送 API 请求并解码响应
// v: 响应结果的接收对象指针。如果为 nil，则只检查错误。
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查标准错误响应
	if err := CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		// 如果是 io.Writer，直接拷贝
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			// 否则解析 JSON
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}
