package confluence

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SpaceService 处理空间(Space)相关操作
type SpaceService struct {
	client *Client
}

// Space 空间信息 (增强版)
type Space struct {
	ID          int64             `json:"id,omitempty"`
	Key         string            `json:"key"`
	Name        string            `json:"name"`
	Type        string            `json:"type,omitempty"` // global, personal
	Description *SpaceDescription `json:"description,omitempty"`
	Homepage    *Content          `json:"homepage,omitempty"`
	Links       *Links            `json:"_links,omitempty"`
}

// SpaceDescription 空间描述
type SpaceDescription struct {
	Plain *BodyContent `json:"plain,omitempty"`
	View  *BodyContent `json:"view,omitempty"`
}

// GetSpaceOptions 获取空间的选项
type GetSpaceOptions struct {
	Expand []string // 需要展开的字段，如 "description.plain", "homepage"
}

// CreateSpaceOptions 创建空间的选项
type CreateSpaceOptions struct {
	Private bool // 是否创建私有空间
}

// Get 获取指定 Space 信息
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-space/#api-wiki-rest-api-space-spacekey-get
func (s *SpaceService) Get(ctx context.Context, spaceKey string, opts *GetSpaceOptions) (*Space, *http.Response, error) {
	path := fmt.Sprintf("rest/api/space/%s", url.PathEscape(spaceKey))

	if opts != nil && len(opts.Expand) > 0 {
		path += "?expand=" + strings.Join(opts.Expand, ",")
	}

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var space Space
	resp, err := s.client.Do(req, &space)
	if err != nil {
		return nil, resp, err
	}

	return &space, resp, nil
}

// Create 创建一个新的 Space
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-space/#api-wiki-rest-api-space-post
func (s *SpaceService) Create(ctx context.Context, space *Space, opts *CreateSpaceOptions) (*Space, *http.Response, error) {
	// 目前暂不支持 CreateSpaceOptions.Private 的处理，因为涉及复杂的 Permissions 结构
	// 留待后续完善
	req, err := s.client.NewRequest(ctx, "POST", "rest/api/space", space)
	if err != nil {
		return nil, nil, err
	}

	var createdSpace Space
	resp, err := s.client.Do(req, &createdSpace)
	if err != nil {
		return nil, resp, err
	}

	return &createdSpace, resp, nil
}

// Update 更新 Space 信息
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-space/#api-wiki-rest-api-space-spacekey-put
func (s *SpaceService) Update(ctx context.Context, spaceKey string, space *Space) (*Space, *http.Response, error) {
	path := fmt.Sprintf("rest/api/space/%s", url.PathEscape(spaceKey))
	req, err := s.client.NewRequest(ctx, "PUT", path, space)
	if err != nil {
		return nil, nil, err
	}

	var updatedSpace Space
	resp, err := s.client.Do(req, &updatedSpace)
	if err != nil {
		return nil, resp, err
	}

	return &updatedSpace, resp, nil
}

// Delete 删除 Space
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-space/#api-wiki-rest-api-space-spacekey-delete
func (s *SpaceService) Delete(ctx context.Context, spaceKey string) (*http.Response, error) {
	path := fmt.Sprintf("rest/api/space/%s", url.PathEscape(spaceKey))
	req, err := s.client.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
