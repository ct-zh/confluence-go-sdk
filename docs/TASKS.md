# 架构规划 - 阶段 4: 内容创建与附件上传 (Content Creation & Attachment Upload)

## 1. 概述 (Overview)
Phase 3 完成了附件和层级结构的**读取**能力。为了实现完整的文档管理闭环，Phase 4 将聚焦于**写入**能力，重点攻克最复杂的 **附件上传 (Attachment Upload)** 和 **更新**。

## 2. 架构设计 (Architecture)

### 2.1 客户端重构 (Client Refactoring)
目前的 `Client.NewRequest` 仅支持 JSON Body。为了支持文件上传，必须扩展其能力以处理 `io.Reader` 和 `multipart/form-data`。

**变更点**:
- 重构 `NewRequest` 或新增 `NewUploadRequest` 方法。
- 支持设置 `Content-Type` 为 `multipart/form-data`。

### 2.2 附件上传 (Attachment Upload)
附件上传 API (`POST /rest/api/content/{id}/child/attachment`) 需要特殊的请求体构造。

```go
type IContentService interface {
    // ... existing methods
    
    // UploadAttachment 上传附件
    // comment: 附件的注释（可选）
    UploadAttachment(ctx context.Context, contentID string, filename string, data io.Reader, comment string) (*SearchResult, *http.Response, error)
}
```

### 2.3 标签管理 (Label Management) - 可选
为了更好地组织内容，支持标签 (Labels) 的添加和删除也是必要的。

## 3. 开发任务 (Tasks)

请研发工程师按以下顺序执行：

- [x] **Core Refactor**: 修改 `confluence/client.go`，使 `NewRequest` 支持 `io.Reader` 作为 Body，或添加 `NewUploadRequest` 辅助方法。
- [x] **Interface**: 更新 `confluence/interfaces.go`，在 `IContentService` 中添加 `UploadAttachment`。
- [x] **Implementation**: 在 `confluence/content.go` 中实现 `UploadAttachment`。
    - 使用 `mime/multipart` 构造请求体。
    - 确保正确设置 Boundary。
- [x] **Test**: 更新 `confluence/content_test.go`，增加附件上传的 Mock 测试。

---
**历史记录 (History)**:
- [x] Phase 1: 核心架构与基础 CRUD (Content, Search)
- [x] Phase 2: 空间管理 (Space Management)
- [x] Phase 3: 附件与层级管理 (Attachments & Hierarchy) - Read Only
