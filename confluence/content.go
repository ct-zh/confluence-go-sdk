package confluence

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ContentService 处理页面(Page)和博客(Blogpost)相关操作
type ContentService struct {
	client *Client
}

// Content 表示 Confluence 中的内容实体
// 仅包含核心字段，遵循 Less is More 原则
type Content struct {
	ID        string    `json:"id,omitempty"` // Create 时不需要 ID
	Type      string    `json:"type"`         // page, blogpost
	Status    string    `json:"status,omitempty"`
	Title     string    `json:"title"`
	Space     *Space    `json:"space,omitempty"`
	Body      *Body     `json:"body,omitempty"`
	Version   *Version  `json:"version,omitempty"`
	Ancestors []Content `json:"ancestors,omitempty"`
	Links     *Links    `json:"_links,omitempty"`
}

// Version 版本信息
type Version struct {
	Number    int    `json:"number"`
	Message   string `json:"message,omitempty"`
	MinorEdit bool   `json:"minorEdit,omitempty"`
}


// Body 内容正文
type Body struct {
	Storage *BodyContent `json:"storage,omitempty"` // 存储格式 (XHTML)
	View    *BodyContent `json:"view,omitempty"`    // 显示格式 (HTML)
}

// BodyContent 具体格式的内容
type BodyContent struct {
	Value          string `json:"value"`
	Representation string `json:"representation"` // storage, view
}

// Links 通用链接结构
type Links struct {
	Self string `json:"self"`
	Base string `json:"base"`
	Web  string `json:"webui"`
}

// GetContentOptions 获取内容的选项
type GetContentOptions struct {
	Expand []string // 需要展开的字段，如 "body.storage", "space", "version"
}

// Get 根据 ID 获取单篇内容
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content/#api-wiki-rest-api-content-id-get
func (s *ContentService) Get(ctx context.Context, contentID string, opts *GetContentOptions) (*Content, *http.Response, error) {
	path := fmt.Sprintf("rest/api/content/%s", contentID)

	if opts != nil && len(opts.Expand) > 0 {
		path += "?expand=" + strings.Join(opts.Expand, ",")
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var content Content
	resp, err := s.client.Do(req, &content)
	if err != nil {
		return nil, resp, err
	}

	return &content, resp, nil
}

// Create 创建新内容 (Page 或 Blogpost)
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content/#api-wiki-rest-api-content-post
func (s *ContentService) Create(ctx context.Context, content *Content) (*Content, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, "POST", "rest/api/content", content)
	if err != nil {
		return nil, nil, err
	}

	var createdContent Content
	resp, err := s.client.Do(req, &createdContent)
	if err != nil {
		return nil, resp, err
	}

	return &createdContent, resp, nil
}

// Update 更新现有内容
// 注意: 必须提供 content.Version.Number (通常是当前版本号 + 1)
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content/#api-wiki-rest-api-content-id-put
func (s *ContentService) Update(ctx context.Context, contentID string, content *Content) (*Content, *http.Response, error) {
	path := fmt.Sprintf("rest/api/content/%s", contentID)
	req, err := s.client.NewRequest(ctx, "PUT", path, content)
	if err != nil {
		return nil, nil, err
	}

	var updatedContent Content
	resp, err := s.client.Do(req, &updatedContent)
	if err != nil {
		return nil, resp, err
	}

	return &updatedContent, resp, nil
}

// Delete 删除内容
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content/#api-wiki-rest-api-content-id-delete
func (s *ContentService) Delete(ctx context.Context, contentID string) (*http.Response, error) {
	path := fmt.Sprintf("rest/api/content/%s", contentID)
	req, err := s.client.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
