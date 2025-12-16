package mocks

import (
	"context"
	"net/http"

	"confluence-go-sdk/pkg/confluence"
)

// MockContentService 是 IContentService 的 Mock 实现
type MockContentService struct {
	GetFunc    func(ctx context.Context, contentID string, opts *confluence.GetContentOptions) (*confluence.Content, *http.Response, error)
	CreateFunc func(ctx context.Context, content *confluence.Content) (*confluence.Content, *http.Response, error)
	UpdateFunc func(ctx context.Context, contentID string, content *confluence.Content) (*confluence.Content, *http.Response, error)
	DeleteFunc func(ctx context.Context, contentID string) (*http.Response, error)
}

func (m *MockContentService) Get(ctx context.Context, contentID string, opts *confluence.GetContentOptions) (*confluence.Content, *http.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, contentID, opts)
	}
	return nil, nil, nil
}

func (m *MockContentService) Create(ctx context.Context, content *confluence.Content) (*confluence.Content, *http.Response, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, content)
	}
	return nil, nil, nil
}

func (m *MockContentService) Update(ctx context.Context, contentID string, content *confluence.Content) (*confluence.Content, *http.Response, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, contentID, content)
	}
	return nil, nil, nil
}

func (m *MockContentService) Delete(ctx context.Context, contentID string) (*http.Response, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, contentID)
	}
	return nil, nil
}

// MockSearchService 是 ISearchService 的 Mock 实现
type MockSearchService struct {
	SearchFunc func(ctx context.Context, opts *confluence.SearchOptions) (*confluence.SearchResult, *http.Response, error)
}

func (m *MockSearchService) Search(ctx context.Context, opts *confluence.SearchOptions) (*confluence.SearchResult, *http.Response, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, opts)
	}
	return nil, nil, nil
}
