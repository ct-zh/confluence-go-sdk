# Tasks: Standard Template Library

作为架构师，我已定义了标准模板库的架构。请按照以下步骤进行实现。

## 1. 基础设施建设
- [ ] **创建包结构**: 在 `confluence/template` 下创建 `loader.go` 和 `constants.go`。
- [ ] **实现 Embed FS**: 使用 `//go:embed *.html` 将所有模板文件嵌入。
- [ ] **定义加载接口**:
    ```go
    // Loader 定义模板加载行为
    type Loader interface {
        Load(name string) (string, error)
    }
    
    // BuiltinLoader 加载内置模板
    // CustomLoader 加载用户文件
    ```

## 2. 迁移与新增模板
- [ ] **SwaggerDoc**: 将 `integration_test.go` 中的 Swagger 模板迁移至 `confluence/template/swagger.html`。
- [ ] **ReleaseNote**: 创建一个标准的发布说明模板 `release_note.html`。
- [ ] **MeetingNotes**: 创建一个会议纪要模板 `meeting_notes.html`。

## 3. 集成与测试
- [ ] **更新 TemplateRenderer**: 使其支持通过 `Loader` 获取模板内容。
- [ ] **单元测试**: 验证内置模板是否能被正确加载和解析。
- [ ] **集成测试**: 更新 `TestIntegration_CreatePageFromSwagger` 使用内置模板。
