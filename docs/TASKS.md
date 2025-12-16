# 架构规划 - 阶段 3: 附件与层级管理 (Attachments & Hierarchy)

## 1. 概述 (Overview)
SDK 目前已具备 Content, Search, Space 三大核心模块。为了满足更深度的文档管理需求，下一阶段将聚焦于 **非结构化数据 (Attachments)** 和 **文档结构 (Child Pages)** 的处理。

## 2. 架构设计 (Architecture)

### 2.1 附件服务 (AttachmentService)
附件虽然属于 `Content` 的一种特殊类型，但其操作（上传、下载）涉及二进制流处理，建议独立封装或在 `ContentService` 中通过专用方法处理。为了保持接口简洁，我们将在 `ContentService` 中扩展。

**变更点**:
- 扩展 `IContentService` 接口，增加附件相关方法。
- 引入 `multipart/form-data` 处理逻辑。

```go
type IContentService interface {
    // ... existing methods
    
    // GetChildPages 获取子页面
    GetChildPages(ctx context.Context, contentID string, opts *GetChildPagesOptions) (*SearchResult, *http.Response, error)
    
    // GetAttachments 获取页面的附件列表
    GetAttachments(ctx context.Context, contentID string, opts *GetAttachmentsOptions) (*SearchResult, *http.Response, error)
}
```

*注意*: 附件的**上传** (Upload) 逻辑较为复杂（需要 `multipart`），暂定为 P1 优先级，先实现获取 (Get) 逻辑。

### 2.2 层级遍历 (Hierarchy Traversal)
用户经常需要获取某个页面下的所有子页面。虽然可以通过 CQL 实现，但提供专用的 `GetChildPages` 方法能显著提升易用性。

## 3. 开发任务 (Tasks)

请研发工程师按以下顺序执行：

- [x] **Interface**: 更新 `confluence/interfaces.go`，在 `IContentService` 中添加 `GetChildPages` 和 `GetAttachments`。
- [x] **Implementation**: 在 `confluence/content.go` 中实现上述方法。
    - `GetChildPages`: 实际上是调用 `/rest/api/content/{id}/child/page`。
    - `GetAttachments`: 调用 `/rest/api/content/{id}/child/attachment`。
- [x] **Test**: 更新 `confluence/content_test.go`，增加对应的 Mock 测试。

---
**技术债 (Tech Debt)**:
- 目前 `Client` 的 `NewRequest` 仅支持 JSON Body，后续支持附件上传时需要重构以支持 `io.Reader` 和 `multipart`。
