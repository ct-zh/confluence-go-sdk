package confluence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	tpl "confluence-go-sdk/pkg/confluence/template"
)

// IntegrationConfig 集成测试配置
type IntegrationConfig struct {
	BaseURL         string `json:"base_url"`
	Username        string `json:"username"`
	Token           string `json:"token"`
	TestSpaceKey    string `json:"test_space_key"`
	TestContentID   string `json:"test_content_id"`
	TestWritePageID string `json:"test_write_page_id"` // 用于写操作测试的父页面ID
}

// loadIntegrationConfig 加载集成测试配置
// 如果配置文件不存在，返回 nil，表示跳过集成测试
func loadIntegrationConfig(t *testing.T) *IntegrationConfig {
	data, err := os.ReadFile("test_config.json")
	if os.IsNotExist(err) {
		t.Skip("Skipping integration test: test_config.json not found")
		return nil
	}
	if err != nil {
		t.Fatalf("Failed to read test_config.json: %v", err)
	}

	var config IntegrationConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse test_config.json: %v", err)
	}

	if config.BaseURL == "" || config.Token == "" {
		t.Skip("Skipping integration test: BaseURL or Token is empty")
	}

	return &config
}

// createIntegrationClient 创建用于集成测试的 Client
func createIntegrationClient(t *testing.T, config *IntegrationConfig) *Client {
	var opts []Option
	if config.Username != "" {
		// 如果提供了用户名，使用 Basic Auth (Email + API Token)
		opts = append(opts, WithAuth(config.Username, config.Token))
	} else {
		// 如果没有用户名，使用 Bearer Token (PAT)
		opts = append(opts, WithToken(config.Token))
	}

	client, err := NewClient(config.BaseURL, opts...)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return client
}

// TestIntegration_GetSpace 集成测试：获取空间信息 (只读)
func TestIntegration_GetSpace(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	if config.TestSpaceKey == "" {
		t.Skip("Skipping TestIntegration_GetSpace: test_space_key not configured")
	}

	ctx := context.Background()
	space, _, err := client.Space.Get(ctx, config.TestSpaceKey, nil)
	if err != nil {
		t.Errorf("Failed to get space %s: %v", config.TestSpaceKey, err)
		return
	}

	t.Logf("Successfully retrieved space: %s (ID: %d)", space.Name, space.ID)
	if space.Key != config.TestSpaceKey {
		t.Errorf("Expected space key %s, got %s", config.TestSpaceKey, space.Key)
	}
}

// TestIntegration_GetContent 集成测试：获取内容信息 (只读)
func TestIntegration_GetContent(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	if config.TestContentID == "" {
		t.Skip("Skipping TestIntegration_GetContent: test_content_id not configured")
	}

	ctx := context.Background()
	content, _, err := client.Content.Get(ctx, config.TestContentID, nil)
	if err != nil {
		t.Errorf("Failed to get content %s: %v", config.TestContentID, err)
		return
	}

	t.Logf("Successfully retrieved content: %s (Type: %s)", content.Title, content.Type)
	if content.ID != config.TestContentID {
		t.Errorf("Expected content ID %s, got %s", config.TestContentID, content.ID)
	}
}

// TestIntegration_Search 集成测试：搜索 (只读)
func TestIntegration_Search(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	// 搜索最近更新的内容
	cql := "type in (page,blogpost) order by lastModified desc"
	opts := &SearchOptions{
		CQL:   cql,
		Limit: 5,
	}

	ctx := context.Background()
	result, _, err := client.Search.Search(ctx, opts)
	if err != nil {
		t.Errorf("Failed to search: %v", err)
		return
	}

	t.Logf("Successfully searched. Found %d results.", len(result.Results))
	for _, item := range result.Results {
		t.Logf(" - %s (ID: %s)", item.Title, item.ID)
	}
}

// TestIntegration_GetChildPages 集成测试：获取子页面 (只读)
func TestIntegration_GetChildPages(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	if config.TestContentID == "" {
		t.Skip("Skipping TestIntegration_GetChildPages: test_content_id not configured")
	}

	ctx := context.Background()
	// 获取前 5 个子页面
	opts := &GetChildPagesOptions{
		Limit: 5,
	}
	result, _, err := client.Content.GetChildPages(ctx, config.TestContentID, opts)
	if err != nil {
		t.Errorf("Failed to get child pages: %v", err)
		return
	}

	t.Logf("Successfully retrieved %d child pages.", len(result.Results))
	for _, page := range result.Results {
		t.Logf(" - %s (ID: %s)", page.Title, page.ID)
	}
}

// TestIntegration_GetAttachments 集成测试：获取附件 (只读)
func TestIntegration_GetAttachments(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	if config.TestContentID == "" {
		t.Skip("Skipping TestIntegration_GetAttachments: test_content_id not configured")
	}

	ctx := context.Background()
	// 获取前 5 个附件
	opts := &GetAttachmentsOptions{
		Limit: 5,
	}
	result, _, err := client.Content.GetAttachments(ctx, config.TestContentID, opts)
	if err != nil {
		t.Errorf("Failed to get attachments: %v", err)
		return
	}

	t.Logf("Successfully retrieved %d attachments.", len(result.Results))
	for _, att := range result.Results {
		t.Logf(" - %s (Type: %s, ID: %s)", att.Title, att.Type, att.ID)
	}
}

// TestIntegration_ChildPageLifecycle 集成测试：子页面生命周期 (增删改查)
func TestIntegration_ChildPageLifecycle(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	parentID := config.TestWritePageID
	if parentID == "" {
		// Fallback for current user session if not configured
		parentID = "73203845"
		t.Logf("Using hardcoded parentID: %s", parentID)
	}

	if config.TestSpaceKey == "" {
		t.Skip("Skipping TestIntegration_ChildPageLifecycle: test_space_key not configured")
	}

	ctx := context.Background()
	timestamp := time.Now().Unix()
	title := fmt.Sprintf("Integration Test Page %d", timestamp)

	// 1. Create
	newContent := &Content{
		Title:     title,
		Type:      "page",
		Space:     &Space{Key: config.TestSpaceKey},
		Ancestors: []Content{{ID: parentID}},
		Body: &Body{
			Storage: &BodyContent{
				Value:          "<p>Created by Integration Test</p>",
				Representation: "storage",
			},
		},
	}

	created, _, err := client.Content.Create(ctx, newContent)
	if err != nil {
		t.Fatalf("Failed to create child page: %v", err)
	}
	t.Logf("Successfully created page: %s (ID: %s)", created.Title, created.ID)

	// Ensure cleanup
	defer func() {
		_, err := client.Content.Delete(ctx, created.ID)
		if err == nil {
			t.Logf("Cleanup: Successfully deleted page %s", created.ID)
		}
	}()

	// 2. Get
	fetched, _, err := client.Content.Get(ctx, created.ID, nil)
	if err != nil {
		t.Errorf("Failed to get created page: %v", err)
	}
	if fetched.Title != title {
		t.Errorf("Expected title %s, got %s", title, fetched.Title)
	}

	// 3. Update
	updatedTitle := title + " - Updated"
	fetched.Title = updatedTitle
	// 乐观锁：Version 需要递增
	if fetched.Version == nil {
		fetched.Version = &Version{Number: 1}
	}
	fetched.Version.Number++

	updated, _, err := client.Content.Update(ctx, created.ID, fetched)
	if err != nil {
		t.Errorf("Failed to update page: %v", err)
	} else {
		t.Logf("Successfully updated page: %s (Version: %d)", updated.Title, updated.Version.Number)

		if updated.Title != updatedTitle {
			t.Errorf("Expected updated title %s, got %s", updatedTitle, updated.Title)
		}
	}

	// 4. Delete is handled by defer, but we can test explicit delete here
	// 为了验证 Delete 功能，我们手动删除，然后让 defer 里的删除失败（或者 defer 检查 ID 是否存在）
	// 这里简单起见，让 defer 做清理。如果想明确测试 Delete，可以在这里做。
	_, err = client.Content.Delete(ctx, created.ID)
	if err != nil {
		t.Errorf("Failed to delete page: %v", err)
	} else {
		t.Logf("Successfully deleted page: %s", created.ID)
	}

	// 5. Verify Delete
	_, resp, err := client.Content.Get(ctx, created.ID, nil)
	if err == nil {
		t.Errorf("Expected error getting deleted page, got nil")
	}
	if resp != nil && resp.StatusCode != 404 {
		t.Errorf("Expected 404 for deleted page, got %d", resp.StatusCode)
	}
}

// TestIntegration_UploadAttachment 集成测试：上传附件
func TestIntegration_UploadAttachment(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	pageID := config.TestWritePageID
	if pageID == "" {
		pageID = "73203845"
		t.Logf("Using hardcoded pageID for attachment: %s", pageID)
	}

	ctx := context.Background()
	filename := fmt.Sprintf("test_attachment_%d.txt", time.Now().Unix())
	content := []byte("This is a test attachment content created by integration test.")
	reader := bytes.NewReader(content)

	results, _, err := client.Content.UploadAttachment(ctx, pageID, filename, reader, "Integration test attachment")
	if err != nil {
		t.Fatalf("Failed to upload attachment: %v", err)
	}

	if len(results.Results) == 0 {
		t.Errorf("Expected at least one attachment result")
	} else {
		t.Logf("Successfully uploaded attachment: %s (ID: %s)", results.Results[0].Title, results.Results[0].ID)
	}
}

// TestIntegration_CreatePageFromSwagger 集成测试：使用模板生成 Swagger 文档页面
func TestIntegration_CreatePageFromSwagger(t *testing.T) {
	config := loadIntegrationConfig(t)
	client := createIntegrationClient(t, config)

	parentID := config.TestWritePageID
	if parentID == "" {
		parentID = "73203845"
		t.Logf("Using hardcoded parentID: %s", parentID)
	}

	if config.TestSpaceKey == "" {
		t.Skip("Skipping TestIntegration_CreatePageFromSwagger: test_space_key not configured")
	}

	// 1. 模拟 Swagger 数据
	type Parameter struct {
		Name        string
		Type        string
		Description string
		Required    bool
	}
	type Endpoint struct {
		Method     string
		Path       string
		Summary    string
		Parameters []Parameter
		Response   string // JSON example
	}
	type SwaggerDoc struct {
		Title       string
		Version     string
		Description string
		Endpoints   []Endpoint
	}

	docData := SwaggerDoc{
		Title:       "User Service API",
		Version:     "v1.0.0",
		Description: "API for managing users in the system.",
		Endpoints: []Endpoint{
			{
				Method:  "GET",
				Path:    "/users/{id}",
				Summary: "Get user by ID",
				Parameters: []Parameter{
					{Name: "id", Type: "string", Description: "User UUID", Required: true},
				},
				Response: `{"id": "123", "name": "John Doe", "email": "john@example.com"}`,
			},
			{
				Method:  "POST",
				Path:    "/users",
				Summary: "Create new user",
				Parameters: []Parameter{
					{Name: "name", Type: "string", Description: "Full Name", Required: true},
					{Name: "email", Type: "string", Description: "Email Address", Required: true},
				},
				Response: `{"id": "124", "status": "created"}`,
			},
		},
	}

	// 2. 渲染模板
	// 使用内置的 Swagger 模板
	renderer := NewTemplateRenderer(nil) // nil means use BuiltinLoader

	htmlContent, err := renderer.LoadAndRender(tpl.SwaggerDoc, docData)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 3. 创建页面
	ctx := context.Background()
	pageTitle := fmt.Sprintf("API Docs - %s (%d)", docData.Title, time.Now().Unix())

	newContent := &Content{
		Title:     pageTitle,
		Type:      "page",
		Space:     &Space{Key: config.TestSpaceKey},
		Ancestors: []Content{{ID: parentID}},
		Body: &Body{
			Storage: &BodyContent{
				Value:          htmlContent,
				Representation: "storage",
			},
		},
	}

	created, _, err := client.Content.Create(ctx, newContent)
	if err != nil {
		t.Fatalf("Failed to create swagger page: %v", err)
	}
	t.Logf("Successfully created Swagger doc page: %s (ID: %s)", created.Title, created.ID)

	// Cleanup
	defer func() {
		_, err := client.Content.Delete(ctx, created.ID)
		if err == nil {
			t.Logf("Cleanup: Successfully deleted swagger page %s", created.ID)
		}
	}()
}
