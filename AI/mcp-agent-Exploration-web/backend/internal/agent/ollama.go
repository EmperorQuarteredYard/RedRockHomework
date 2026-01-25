package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// OllamaMessage Ollama API 消息结构
type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaRequest Ollama API 请求结构
type OllamaRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  OllamaOptions   `json:"options"`
}

// OllamaOptions Ollama 选项
type OllamaOptions struct {
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	NumCtx      int     `json:"num_ctx"`
}

// OllamaResponse Ollama API 响应结构
type OllamaResponse struct {
	Model     string        `json:"model"`
	CreatedAt time.Time     `json:"created_at"`
	Message   OllamaMessage `json:"message"`
	Done      bool          `json:"done"`
}

// OllamaClient Ollama 客户端
type OllamaClient struct {
	baseURL   string
	modelName string
	client    *http.Client
	logger    *zap.Logger
}

// NewOllamaClient 创建新的 Ollama 客户端
func NewOllamaClient(baseURL, modelName string, logger *zap.Logger) *OllamaClient {
	return &OllamaClient{
		baseURL:   baseURL,
		modelName: modelName,
		client: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxConnsPerHost:     10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		logger: logger,
	}
}

// Chat 与 Ollama 对话
func (c *OllamaClient) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := []OllamaMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	request := OllamaRequest{
		Model:    c.modelName,
		Messages: messages,
		Stream:   false,
		Options: OllamaOptions{
			Temperature: 0.3,
			TopP:        0.9,
			NumCtx:      4096,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API 错误: %s, 响应: %s", resp.Status, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	processingTime := time.Since(startTime)
	c.logger.Debug("Ollama 调用完成",
		zap.Duration("processing_time", processingTime),
		zap.String("model", ollamaResp.Model),
	)

	return ollamaResp.Message.Content, nil
}

// CheckConnection 检查 Ollama 连接
func (c *OllamaClient) CheckConnection(ctx context.Context) (bool, string) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return false, fmt.Sprintf("创建请求失败: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Sprintf("解析响应失败: %v", err)
	}

	// 检查模型是否可用
	for _, model := range result.Models {
		if model.Name == c.modelName {
			return true, "model available"
		}
	}

	return true, "model not found"
}
