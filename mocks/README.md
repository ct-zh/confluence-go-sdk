# 单元测试指南：使用 Mock 对象

本 SDK 提供了 `mocks` 包，帮助您在编写业务代码的单元测试时，轻松模拟 Confluence API 的行为，而无需发起真实的 HTTP 请求。

## 核心接口

我们在 `confluence` 包中定义了核心服务接口：

- `IContentService`: 包含 `Get`, `Create`, `Update`, `Delete` 方法。
- `ISearchService`: 包含 `Search` 方法。

## 如何使用 Mock

以下是一个示例，展示如何测试一个依赖 Confluence SDK 的业务函数。

假设您有一个函数 `GetPageTitle`：

```go
func GetPageTitle(service confluence.IContentService, pageID string) (string, error) {
    content, _, err := service.Get(context.Background(), pageID, nil)
    if err != nil {
        return "", err
    }
    return content.Title, nil
}
```

您可以这样编写测试：

```go
package main_test

import (
	"context"
	"errors"
	"testing"

	"confluence-go-sdk/confluence"
	"confluence-go-sdk/mocks"
)

func TestGetPageTitle(t *testing.T) {
	// 1. 创建 Mock 对象
	mockService := &mocks.MockContentService{
		// 2. 注入模拟行为
		GetFunc: func(ctx context.Context, id string, opts *confluence.GetContentOptions) (*confluence.Content, *http.Response, error) {
			if id == "error-id" {
				return nil, nil, errors.New("mock error")
			}
			return &confluence.Content{Title: "Mock Page Title"}, nil, nil
		},
	}

	// 3. 执行测试 (成功场景)
	title, err := GetPageTitle(mockService, "123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if title != "Mock Page Title" {
		t.Errorf("Expected 'Mock Page Title', got '%s'", title)
	}

	// 4. 执行测试 (失败场景)
	_, err = GetPageTitle(mockService, "error-id")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
```

## 优势

1.  **无需网络**：测试运行极快。
2.  **确定性**：您可以精确控制 API 返回的数据（包括错误情况）。
3.  **简单**：Mock 实现基于函数回调，无需复杂的 Mock 框架。
