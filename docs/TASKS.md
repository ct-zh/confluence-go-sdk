# 架构规划 - 阶段 2: 空间管理 (Space Management)

## 1. 概述 (Overview)
目前 SDK 已支持基础的 Content (Page/Blogpost) 操作和 Search 功能。为了完善 SDK 的能力，下一个关键组件是 **Space Management (空间管理)**。

## 2. 架构设计 (Architecture)

我们将引入 `SpaceService`，遵循现有的 `Service` 模式。

### 2.1 接口定义 (Interface Definition)
在 `confluence/interfaces.go` 中添加 `ISpaceService`：

```go
// ISpaceService 定义 SpaceService 的行为接口
type ISpaceService interface {
    // Get 获取指定 Space 信息
    Get(ctx context.Context, spaceKey string, opts *GetSpaceOptions) (*Space, *http.Response, error)
    
    // Create 创建一个新的 Space
    Create(ctx context.Context, space *Space, opts *CreateSpaceOptions) (*Space, *http.Response, error)
    
    // Update 更新 Space 信息
    Update(ctx context.Context, spaceKey string, space *Space) (*Space, *http.Response, error)
    
    // Delete 删除 Space
    Delete(ctx context.Context, spaceKey string) (*http.Response, error)
}
```

### 2.2 数据结构 (Data Structures)
需要扩展 `Space` 结构体（可能需要从 `content.go` 中移动或增强），并在 `confluence/space.go` 中定义。

```go
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

type SpaceDescription struct {
    Plain *BodyContent `json:"plain,omitempty"`
    View  *BodyContent `json:"view,omitempty"`
}

type GetSpaceOptions struct {
    Expand []string // description.plain, homepage
}

type CreateSpaceOptions struct {
    Private bool // 是否创建私有空间
}
```

### 2.3 客户端集成 (Client Integration)
在 `confluence/client.go` 中注册新服务。

```go
type Client struct {
    // ... 现有字段
    Space ISpaceService
}

// NewClient ...
    c.Space = &SpaceService{client: c}
```

## 3. 开发任务 (Tasks)

请研发工程师按以下顺序执行：

- [ ] **Refactor**: 检查 `confluence/content.go` 中的 `Space` 定义，准备将其增强或移动到新文件。
- [ ] **Interface**: 在 `confluence/interfaces.go` 中添加 `ISpaceService` 接口定义。
- [ ] **Implementation**: 创建 `confluence/space.go`，实现 `SpaceService` 及其方法 (`Get`, `Create`, `Update`, `Delete`)。
- [ ] **Integration**: 更新 `confluence/client.go`，将 `SpaceService` 挂载到 `Client`。
- [ ] **Test**: 创建 `confluence/space_test.go`，编写单元测试（使用 Mock 或集成测试）。

---
**注意**: 请保持代码风格一致，确保所有公共方法都有注释。
