package confluence

import (
	"context"
	"net/http"
)

// IContentService 定义 ContentService 的行为接口
// 使用接口可以让用户在测试时轻松注入 Mock 对象
type IContentService interface {
	Get(ctx context.Context, contentID string, opts *GetContentOptions) (*Content, *http.Response, error)
	Create(ctx context.Context, content *Content) (*Content, *http.Response, error)
	Update(ctx context.Context, contentID string, content *Content) (*Content, *http.Response, error)
	Delete(ctx context.Context, contentID string) (*http.Response, error)
	// GetChildPages 获取子页面
	GetChildPages(ctx context.Context, contentID string, opts *GetChildPagesOptions) (*SearchResult, *http.Response, error)
	// GetAttachments 获取页面的附件列表
	GetAttachments(ctx context.Context, contentID string, opts *GetAttachmentsOptions) (*SearchResult, *http.Response, error)
}

// ISearchService 定义 SearchService 的行为接口
type ISearchService interface {
	Search(ctx context.Context, opts *SearchOptions) (*SearchResult, *http.Response, error)
}

// ISpaceService 定义 SpaceService 的行为接口
type ISpaceService interface {
	Get(ctx context.Context, spaceKey string, opts *GetSpaceOptions) (*Space, *http.Response, error)
	Create(ctx context.Context, space *Space, opts *CreateSpaceOptions) (*Space, *http.Response, error)
	Update(ctx context.Context, spaceKey string, space *Space) (*Space, *http.Response, error)
	Delete(ctx context.Context, spaceKey string) (*http.Response, error)
}
