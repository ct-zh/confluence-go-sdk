package confluence

import (
	"context"
	"io"
	"net/http"
)

// IContentService 定义 ContentService 的行为接口
// 使用接口可以让用户在测试时轻松注入 Mock 对象
type IContentService interface {
	// Get 获取指定 ID 的内容
	Get(ctx context.Context, contentID string, opts *GetContentOptions) (*Content, *http.Response, error)
	// Create 创建新内容 (页面或博客)
	Create(ctx context.Context, content *Content) (*Content, *http.Response, error)
	// Update 更新现有内容
	Update(ctx context.Context, contentID string, content *Content) (*Content, *http.Response, error)
	// Delete 删除内容（移动到回收站）
	Delete(ctx context.Context, contentID string) (*http.Response, error)
	// GetChildPages 获取子页面
	GetChildPages(ctx context.Context, contentID string, opts *GetChildPagesOptions) (*SearchResult, *http.Response, error)
	// GetAttachments 获取页面的附件列表
	GetAttachments(ctx context.Context, contentID string, opts *GetAttachmentsOptions) (*SearchResult, *http.Response, error)
	// UploadAttachment 上传附件
	// comment: 附件的注释（可选）
	UploadAttachment(ctx context.Context, contentID string, filename string, data io.Reader, comment string) (*SearchResult, *http.Response, error)
}

// ISearchService 定义 SearchService 的行为接口
type ISearchService interface {
	// Search 执行 CQL 搜索
	Search(ctx context.Context, opts *SearchOptions) (*SearchResult, *http.Response, error)
}

// ISpaceService 定义 SpaceService 的行为接口
type ISpaceService interface {
	// Get 获取指定 Key 的空间信息
	Get(ctx context.Context, spaceKey string, opts *GetSpaceOptions) (*Space, *http.Response, error)
	// Create 创建新空间
	Create(ctx context.Context, space *Space, opts *CreateSpaceOptions) (*Space, *http.Response, error)
	// Update 更新空间信息
	Update(ctx context.Context, spaceKey string, space *Space) (*Space, *http.Response, error)
	// Delete 删除空间
	Delete(ctx context.Context, spaceKey string) (*http.Response, error)
}
