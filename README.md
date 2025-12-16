# Confluence Go SDK

简洁、优雅且符合 Go 惯用法的 Confluence REST API 客户端。

本项目遵循 "Less is More" 的设计哲学，旨在提供最轻量、核心逻辑最清晰的 API 调用体验。

## 特性

- **依赖注入**：核心组件解耦，易于测试和扩展。
- **零本地配置**：所有参数由调用方显式传入，无隐式环境依赖。
- **上下文感知**：全链路 `context.Context` 支持，从容控制并发与超时。
- **类型安全**：核心业务模型强类型定义。

## 安装

```bash
go get github.com/your-username/confluence-go-sdk
```

## 快速开始

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/your-username/confluence-go-sdk/confluence"
)

func main() {
	// 1. 基础配置（显式传入，无隐式依赖）
	baseURL := "https://your-domain.atlassian.net/wiki"
	user := "your-email@example.com"
	token := "your-api-token"

	// 2. 创建客户端 (符合依赖注入思想，底层 HTTP Client 可被替换)
	// 默认使用 http.DefaultClient，也可注入自定义 Client 以控制超时或代理
	httpClient := &http.Client{Timeout: 10 * time.Second}
	
	client := confluence.NewClient(
		baseURL, 
		confluence.WithHTTPClient(httpClient),
		confluence.WithAuth(user, token),
	)

	// 3. 调用 API - 获取内容
	ctx := context.Background()
	content, _, err := client.Content.Get(ctx, "123456", nil)
	if err != nil {
		log.Fatalf("获取内容失败: %v", err)
	}
	fmt.Printf("标题: %s\n", content.Title)

	// 4. 调用 API - 搜索内容 (CQL)
	searchOpts := &confluence.SearchOptions{
		CQL:   "type=page AND space = 'DS'",
		Limit: 10,
	}
	results, _, _ := client.Search.Search(ctx, searchOpts)
	fmt.Printf("找到 %d 个页面\n", results.Size)

	// 5. 调用 API - 空间操作
	space, _, _ := client.Space.Get(ctx, "DS", nil)
	fmt.Printf("空间名称: %s\n", space.Name)

	// 6. 调用 API - 获取子页面 (Hierarchy)
	childOpts := &confluence.GetChildPagesOptions{
		Limit: 5,
		Expand: []string{"version"},
	}
	children, _, _ := client.Content.GetChildPages(ctx, "123456", childOpts)
	fmt.Printf("找到 %d 个子页面\n", children.Size)

	// 7. 调用 API - 获取附件 (Attachments)
	attachOpts := &confluence.GetAttachmentsOptions{
		Limit: 10,
	}
	attachments, _, _ := client.Content.GetAttachments(ctx, "123456", attachOpts)
	for _, att := range attachments.Results {
		fmt.Printf("附件: %s (Type: %s)\n", att.Title, att.Type)
	}
}
```

## 设计哲学

- **简洁 (Simplicity)**: 仅实现核心逻辑，拒绝过度封装。
- **显式 (Explicitness)**: 错误显式处理，参数显式传递。
- **组合 (Composition)**: 优先使用组合而非继承，模块职责单一。

## 参考文档

- [Confluence Cloud REST API v1 (Official)](https://developer.atlassian.com/cloud/confluence/rest/v1/intro/) - 功能最全，支持 CQL。
- [Confluence Cloud REST API v2 (Official)](https://developer.atlassian.com/cloud/confluence/rest/v2/intro/) - 新一代 API，性能更优。

## 许可证

MIT
