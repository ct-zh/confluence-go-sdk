package confluence

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

// IntegrationConfig 集成测试配置
type IntegrationConfig struct {
	BaseURL       string `json:"base_url"`
	Username      string `json:"username"`
	Token         string `json:"token"`
	TestSpaceKey  string `json:"test_space_key"`
	TestContentID string `json:"test_content_id"`
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
