package confluence

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SearchService 处理搜索相关操作
type SearchService struct {
	client *Client
}

// SearchResult 搜索结果分页结构
type SearchResult struct {
	Results []Content `json:"results"`
	Start   int       `json:"start"`
	Limit   int       `json:"limit"`
	Size    int       `json:"size"`
	Links   *Links    `json:"_links,omitempty"`
}

// SearchOptions 搜索选项
type SearchOptions struct {
	CQL    string // Confluence Query Language
	Limit  int    // 每页数量
	Start  int    // 起始位置
	Expand []string // 需要展开的字段
}

// Search 执行 CQL 搜索
// 文档: https://developer.atlassian.com/cloud/confluence/rest/v1/api-group-search/#api-wiki-rest-api-search-get
func (s *SearchService) Search(ctx context.Context, opts *SearchOptions) (*SearchResult, *http.Response, error) {
	u := "rest/api/search"
	
	if opts != nil {
		q := url.Values{}
		if opts.CQL != "" {
			q.Set("cql", opts.CQL)
		}
		if opts.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Start > 0 {
			q.Set("start", fmt.Sprintf("%d", opts.Start))
		}
		if len(opts.Expand) > 0 {
			for _, exp := range opts.Expand {
				q.Add("expand", exp)
			}
		}
		u += "?" + q.Encode()
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
