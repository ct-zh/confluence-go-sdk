# LLM 开发指南：核心架构与规范

本文档旨在指导 AI 模型及开发者理解本项目的设计哲学与架构规范。在协助开发新模块时，请务必遵循以下准则。

## 1. 核心原则 (Core Principles)

- **少即是多 (Less is More)**: 代码应精简，拒绝过度设计。只实现当前需要的功能。
- **单一职责 (Single Responsibility)**: 每个包、接口、函数只做一件事。
- **依赖注入 (Dependency Injection)**: 
    - 严禁使用全局变量存储配置（如 `Config`）或状态。
    - 所有外部依赖（如 `http.Client`, `Logger`）必须通过构造函数或 Option 模式注入。
- **Go 惯用法 (Idiomatic Go)**: 遵循官方 Code Review Comments 规范。

## 2. 目录结构 (Directory Structure)

```
confluence-go-sdk/
├── confluence/         # 核心库代码 (对外暴露，包名为 confluence)
│   ├── client.go       # HTTP 客户端封装、鉴权逻辑
│   ├── content.go      # Content (Page/Blog) 相关业务逻辑
│   ├── space.go        # Space 相关业务逻辑
│   ├── user.go         # User 相关业务逻辑
│   ├── errors.go       # 统一错误定义
│   └── options.go      # 配置选项 (Option 模式)
├── internal/           # 内部实现 (对外隐藏)
│   └── util/           # 仅限内部使用的工具函数
├── docs/               # 文档
│   └── llm_ai_architecture.md
├── examples/           # 示例代码
├── README.md
└── go.mod
```

## 3. 代码规范 (Coding Standards)

### 3.1 命名
- **变量**: 短小精悍。上下文清晰时使用单字母（如 `ctx`, `req`）。
- **结构体**: 避免重复包名。例如在 `confluence` 包中，使用 `Client` 而不是 `ConfluenceClient`。
- **Getter**: Go 语言中 Getter 方法不加 `Get` 前缀。

### 3.2 错误处理
- 必须检查每一个 `err`。
- 错误需包含上下文信息，推荐使用 `%w` 进行 wrap。
  ```go
  if err != nil {
      return nil, fmt.Errorf("failed to get content: %w", err)
  }
  ```

### 3.3 并发与 Context
- 所有涉及 I/O 的公共方法**必须**接受 `context.Context` 作为第一个参数。
- 不要在库内部启动 goroutine，除非这是该库的核心功能。

### 3.4 依赖注入与 Option 模式
使用 Functional Options 模式处理可选参数：

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
    auth       Authenticator
}

type Option func(*Client)

func NewClient(baseURL string, opts ...Option) *Client {
    c := &Client{
        baseURL:    baseURL,
        httpClient: http.DefaultClient,
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}
```

## 4. 模块开发流程

当需要开发新的 Confluence 模块（如 `LabelService`）时：

1.  **定义数据结构**: 在对应文件中（如 `label.go`）定义请求/响应结构体。结构体字段应包含 `json` tag。
2.  **定义 Service**: 定义 `LabelService` 结构体，通常包含 `*Client`。
3.  **挂载 Service**: 将 `LabelService` 添加到 `Client` 结构体中，并在 `NewClient` 中初始化。
4.  **实现方法**: 实现具体的 CRUD 方法，确保逻辑清晰。

## 5. 提示词 (Prompting) 检查清单

在生成代码前，请自检：
1.  **是否引入了不必要的第三方库？**（优先使用标准库 `net/http`, `encoding/json`）。
2.  **是否处理了 Context？**
3.  **变量命名是否简洁准确？**
4.  **是否符合依赖注入思想？**
