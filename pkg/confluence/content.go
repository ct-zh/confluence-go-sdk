package confluence

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// ContentService 处理页面(Page)和博客(Blogpost)相关操作
type ContentService struct {
	client *Client
}

// Content 表示 Confluence 中的内容实体
// 仅包含核心字段，遵循 Less is More 原则
type Content struct {
	ID        string    `json:"id,omitempty"` // 内容 ID，创建时无需提供
	Type      string    `json:"type"`         // 内容类型：page (页面) 或 blogpost (博客)
	Status    string    `json:"status,omitempty"` // 内容状态：current (当前), trashed (回收站), historical (历史)
	Title     string    `json:"title"`        // 内容标题
	Space     *Space    `json:"space,omitempty"` // 所属空间
	Body      *Body     `json:"body,omitempty"`  // 内容正文
	Version   *Version  `json:"version,omitempty"` // 版本信息
	Ancestors []Content `json:"ancestors,omitempty"` // 祖先节点（用于构建层级）
	Links     *Links    `json:"_links,omitempty"`    // HATEOAS 链接
}

// Version 版本信息
type Version struct {
	Number    int    `json:"number"`            // 版本号
	Message   string `json:"message,omitempty"` // 版本注释
	MinorEdit bool   `json:"minorEdit,omitempty"` // 是否为微小修改（不触发通知）
}

// Body 内容正文
type Body struct {
	Storage *BodyContent `json:"storage,omitempty"` // 存储格式 (XHTML)，用于编辑和更新
	View    *BodyContent `json:"view,omitempty"`    // 显示格式 (HTML)，用于前端展示
}

// BodyContent 具体格式的内容
type BodyContent struct {
	Value          string `json:"value"`          // 内容字符串
	Representation string `json:"representation"` // 表现形式：storage, view, wiki, plain 等
}

// Links 通用链接结构
type Links struct {
	Self string `json:"self"`  // 自身资源的 API 链接
	Base string `json:"base"`  // Confluence 实例的基础 URL
	Web  string `json:"webui"` // 浏览器可访问的 Web UI 链接
}

// GetContentOptions 获取内容的选项
type GetContentOptions struct {
	Expand []string // 需要展开的字段，如 "body.storage", "space", "version"
}

// GetChildPagesOptions 获取子页面的选项
type GetChildPagesOptions struct {
	Expand []string // 需要展开的字段
	Start  int      // 起始位置
	Limit  int      // 每页数量
}

// GetAttachmentsOptions 获取附件的选项
type GetAttachmentsOptions struct {
	Expand []string // 需要展开的字段
	Start  int      // 起始位置
	Limit  int      // 每页数量
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

// Update 更新内容
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

// GetChildPages 获取子页面
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content-children-and-descendants/#api-wiki-rest-api-content-id-child-page-get
func (s *ContentService) GetChildPages(ctx context.Context, contentID string, opts *GetChildPagesOptions) (*SearchResult, *http.Response, error) {
	u := fmt.Sprintf("rest/api/content/%s/child/page", url.PathEscape(contentID))

	if opts != nil {
		q := url.Values{}
		addPaginationParams(q, opts.Start, opts.Limit, opts.Expand)
		if len(q) > 0 {
			u += "?" + q.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result SearchResult
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetAttachments 获取页面的附件列表
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content-children-and-descendants/#api-wiki-rest-api-content-id-child-attachment-get
func (s *ContentService) GetAttachments(ctx context.Context, contentID string, opts *GetAttachmentsOptions) (*SearchResult, *http.Response, error) {
	u := fmt.Sprintf("rest/api/content/%s/child/attachment", url.PathEscape(contentID))

	if opts != nil {
		q := url.Values{}
		addPaginationParams(q, opts.Start, opts.Limit, opts.Expand)
		if len(q) > 0 {
			u += "?" + q.Encode()
		}
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result SearchResult
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UploadAttachment 上传附件
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-content-children-and-descendants/#api-wiki-rest-api-content-id-child-attachment-post
func (s *ContentService) UploadAttachment(ctx context.Context, contentID string, filename string, data io.Reader, comment string) (*SearchResult, *http.Response, error) {
	u := fmt.Sprintf("rest/api/content/%s/child/attachment", url.PathEscape(contentID))

	// 使用 io.Pipe 实现流式上传，避免大文件占用过多内存
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// 使用 channel 捕获 goroutine 中的错误
	errChan := make(chan error, 1)

	go func() {
		defer pw.Close()
		defer close(errChan)

		// 添加文件
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			errChan <- err
			return
		}
		if _, err := io.Copy(part, data); err != nil {
			errChan <- err
			return
		}

		// 添加注释（可选）
		if comment != "" {
			if err := writer.WriteField("comment", comment); err != nil {
				errChan <- err
				return
			}
		}

		// 关闭 writer 以写入结尾 boundary
		if err := writer.Close(); err != nil {
			errChan <- err
			return
		}
	}()

	// 创建请求
	req, err := s.client.NewRequest(ctx, "POST", u, pr)
	if err != nil {
		pr.Close() // 确保 pipe 被关闭
		return nil, nil, err
	}

	// 手动设置 Content-Type，包含 boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 这是一个特殊 header，用于绕过 XSRF 检查（某些 Confluence 版本需要）
	req.Header.Set("X-Atlassian-Token", "nocheck")

	var result SearchResult
	resp, err := s.client.Do(req, &result)

	// 检查 pipe goroutine 是否发生错误
	if pipeErr := <-errChan; pipeErr != nil {
		// 如果 http 请求也失败了，优先返回 pipe 错误（通常是根本原因）
		// 或者合并错误信息
		if err != nil {
			return nil, resp, fmt.Errorf("upload error: %v, request error: %v", pipeErr, err)
		}
		return nil, resp, pipeErr
	}

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// addPaginationParams 添加分页和展开参数
func addPaginationParams(q url.Values, start, limit int, expand []string) {
	q.Set("start", fmt.Sprintf("%d", start))
	if limit != 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
}
