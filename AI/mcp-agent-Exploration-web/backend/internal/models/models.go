package models

import (
	"time"
)

// ToolName 工具名称枚举
type ToolName string

const (
	ToolSearchWeb    ToolName = "search_web"
	ToolCalculate    ToolName = "calculate"
	ToolAnalyzeText  ToolName = "analyze_text"
	ToolPlanProject  ToolName = "plan_project"
	ToolEstimateTime ToolName = "estimate_time"
	ToolDecompose    ToolName = "decompose_task"
)

// TaskStep 任务步骤
type TaskStep struct {
	StepNumber       int      `json:"step_number"`
	Description      string   `json:"description"`
	EstimatedMinutes float64  `json:"estimated_minutes"`
	ToolsNeeded      []string `json:"tools_needed"`
	Dependencies     []int    `json:"dependencies"`
	Complexity       string   `json:"complexity"`
	ActualMinutes    float64  `json:"actual_minutes,omitempty"`
	Status           string   `json:"status"`
}

// TaskDecompositionRequest 任务拆解请求
type TaskDecompositionRequest struct {
	TaskDescription string `json:"task_description" binding:"required,min=5,max=1000"`
	MaxSteps        int    `json:"max_steps" binding:"min=1,max=20"`
	IncludeTools    bool   `json:"include_tools"`
	ComplexityLevel string `json:"complexity_level" binding:"oneof=auto simple medium complex"`
}

// TaskDecompositionResponse 任务拆解响应
type TaskDecompositionResponse struct {
	TaskID              string     `json:"task_id"`
	OriginalTask        string     `json:"original_task"`
	Steps               []TaskStep `json:"steps"`
	TotalMinutes        float64    `json:"total_minutes"`
	TotalHours          float64    `json:"total_hours"`
	EstimatedCompletion time.Time  `json:"estimated_completion"`
	CreatedAt           time.Time  `json:"created_at"`
}

// ToolInfo 工具信息
type ToolInfo struct {
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Parameters         map[string]string `json:"parameters"`
	EstimatedTimeRange map[string]int    `json:"estimated_time_range"`
}

// AgentStatus Agent状态
type AgentStatus struct {
	ModelName      string    `json:"model_name"`
	OllamaStatus   string    `json:"ollama_status"`
	OllamaDetail   string    `json:"ollama_detail"`
	ToolsAvailable []string  `json:"tools_available"`
	UptimeSeconds  float64   `json:"uptime_seconds"`
	RequestsCount  int64     `json:"requests_count"`
	Status         string    `json:"status"`
	StartedAt      time.Time `json:"started_at"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string    `json:"role" binding:"required,oneof=user assistant system"`
	Content   string    `json:"content" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Message        string `json:"message" binding:"required,min=1,max=2000"`
	ConversationID string `json:"conversation_id"`
	UseHistory     bool   `json:"use_history"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Response       string    `json:"response"`
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id"`
	Timestamp      time.Time `json:"timestamp"`
	ProcessingTime int       `json:"processing_time_ms"`
}

// APIResponse 标准API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}
