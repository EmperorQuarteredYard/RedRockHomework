package agent

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// LLMAgent LLM Agent 核心
type LLMAgent struct {
	modelName    string
	baseURL      string
	tools        map[string]models.ToolInfo
	taskManager  *TaskManager
	requestCount int64
	startTime    time.Time
	logger       *zap.Logger
	ollamaClient *OllamaClient
	mu           sync.RWMutex
}

// NewLLMAgent 创建新的 LLM Agent
func NewLLMAgent(modelName, baseURL string, logger *zap.Logger) *LLMAgent {
	agent := &LLMAgent{
		modelName:    modelName,
		baseURL:      baseURL,
		startTime:    time.Now(),
		logger:       logger,
		ollamaClient: NewOllamaClient(baseURL, modelName, logger),
		tools:        make(map[string]models.ToolInfo),
	}

	agent.initTools()
	agent.taskManager = NewTaskManager()

	return agent
}

// initTools 初始化工具
func (a *LLMAgent) initTools() {
	a.tools = map[string]models.ToolInfo{
		string(models.ToolSearchWeb): {
			Name:        string(models.ToolSearchWeb),
			Description: "搜索网络信息，获取最新数据和资料",
			Parameters: map[string]string{
				"query":       "搜索关键词",
				"max_results": "最大结果数",
			},
			EstimatedTimeRange: map[string]int{
				"simple":  5,
				"medium":  15,
				"complex": 30,
			},
		},
		string(models.ToolCalculate): {
			Name:        string(models.ToolCalculate),
			Description: "执行数学计算和数据处理",
			Parameters: map[string]string{
				"expression": "数学表达式",
				"variables":  "变量字典",
			},
			EstimatedTimeRange: map[string]int{
				"simple":  1,
				"medium":  5,
				"complex": 15,
			},
		},
		string(models.ToolAnalyzeText): {
			Name:        string(models.ToolAnalyzeText),
			Description: "分析文本内容，提取关键信息",
			Parameters: map[string]string{
				"text":          "需要分析的文本",
				"analysis_type": "分析类型（情感、主题、摘要等）",
			},
			EstimatedTimeRange: map[string]int{
				"simple":  3,
				"medium":  10,
				"complex": 20,
			},
		},
		string(models.ToolPlanProject): {
			Name:        string(models.ToolPlanProject),
			Description: "规划项目结构和任务分配",
			Parameters: map[string]string{
				"project_description": "项目描述",
				"team_size":           "团队规模",
			},
			EstimatedTimeRange: map[string]int{
				"simple":  10,
				"medium":  30,
				"complex": 60,
			},
		},
		string(models.ToolEstimateTime): {
			Name:        string(models.ToolEstimateTime),
			Description: "估计任务所需时间",
			Parameters: map[string]string{
				"task":       "任务描述",
				"complexity": "复杂度级别",
			},
			EstimatedTimeRange: map[string]int{
				"simple":  2,
				"medium":  5,
				"complex": 10,
			},
		},
		string(models.ToolDecompose): {
			Name:        string(models.ToolDecompose),
			Description: "拆解复杂任务为子任务",
			Parameters: map[string]string{
				"main_task": "主任务描述",
				"max_depth": "最大拆解深度",
			},
			EstimatedTimeRange: map[string]int{
				"simple":  5,
				"medium":  15,
				"complex": 30,
			},
		},
	}
}

// DecomposeTask 拆解任务
func (a *LLMAgent) DecomposeTask(ctx context.Context, request models.TaskDecompositionRequest) (*models.TaskDecompositionResponse, error) {
	atomic.AddInt64(&a.requestCount, 1)

	startTime := time.Now()
	a.logger.Info("开始拆解任务",
		zap.String("task", request.TaskDescription),
		zap.Int("max_steps", request.MaxSteps),
	)

	// 调用 Ollama API
	systemPrompt := fmt.Sprintf(`你是一个专业的任务规划专家。请将用户的任务拆解成具体的步骤，并为每个步骤估计所需时间。

请以JSON数组格式回复，每个元素包含：
{
  "step_number": 序号,
  "description": "步骤描述",
  "estimated_minutes": 时间(分钟),
  "tools_needed": ["工具列表"],
  "complexity": "simple/medium/complex"
}

可用的工具：
- search_web: 搜索网络信息
- calculate: 执行数学计算
- analyze_text: 分析文本内容
- plan_project: 规划项目
- estimate_time: 估计时间
- decompose_task: 拆解任务

请确保：
1. 步骤数量不超过%d个
2. 每个步骤有明确的描述和合理的时间估计
3. 工具选择要合理
4. 总时间要合理`, request.MaxSteps)

	userPrompt := fmt.Sprintf(`请拆解以下任务：
任务：%s

要求：
1. 将任务分解为可执行的子任务
2. 为每个子任务估计合理的时间
3. 注明每个步骤需要的工具
4. 评估每个步骤的复杂度

请返回JSON数组：`, request.TaskDescription)

	response, err := a.ollamaClient.Chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		a.logger.Error("调用 Ollama 失败", zap.Error(err))
		return nil, fmt.Errorf("调用模型失败: %v", err)
	}

	// 解析响应
	steps, totalMinutes := a.parseTaskDecomposition(response, request.MaxSteps)

	// 创建响应
	taskResponse := &models.TaskDecompositionResponse{
		TaskID:              generateTaskID(),
		OriginalTask:        request.TaskDescription,
		Steps:               steps,
		TotalMinutes:        totalMinutes,
		TotalHours:          totalMinutes / 60,
		EstimatedCompletion: time.Now().Add(time.Duration(totalMinutes) * time.Minute),
		CreatedAt:           time.Now(),
	}

	// 保存任务
	a.taskManager.SaveTask(taskResponse)

	processingTime := time.Since(startTime).Milliseconds()
	a.logger.Info("任务拆解完成",
		zap.String("task_id", taskResponse.TaskID),
		zap.Int("steps", len(steps)),
		zap.Float64("total_minutes", totalMinutes),
		zap.Int64("processing_time_ms", processingTime),
	)

	return taskResponse, nil
}

// parseTaskDecomposition 解析任务拆解响应
func (a *LLMAgent) parseTaskDecomposition(response string, maxSteps int) ([]models.TaskStep, float64) {
	var steps []models.TaskStep
	var totalMinutes float64

	// 尝试解析 JSON - 修复正则表达式语法
	jsonRegex := regexp.MustCompile(`(?s)json\s*(.*?)\s*`)
	jsonMatch := jsonRegex.FindStringSubmatch(response)

	if jsonMatch != nil {
		var jsonSteps []struct {
			StepNumber       int      `json:"step_number"`
			Description      string   `json:"description"`
			EstimatedMinutes float64  `json:"estimated_minutes"`
			ToolsNeeded      []string `json:"tools_needed"`
			Complexity       string   `json:"complexity"`
		}

		if err := json.Unmarshal([]byte(jsonMatch[1]), &jsonSteps); err == nil {
			for i, stepData := range jsonSteps {
				if i >= maxSteps {
					break
				}

				step := models.TaskStep{
					StepNumber:       stepData.StepNumber,
					Description:      stepData.Description,
					EstimatedMinutes: stepData.EstimatedMinutes,
					ToolsNeeded:      stepData.ToolsNeeded,
					Complexity:       stepData.Complexity,
					Status:           "pending",
				}

				if step.StepNumber == 0 {
					step.StepNumber = i + 1
				}

				steps = append(steps, step)
				totalMinutes += step.EstimatedMinutes
			}
		}
	}

	// 如果解析失败，创建默认步骤
	if len(steps) == 0 {
		steps = []models.TaskStep{
			{
				StepNumber:       1,
				Description:      "分析任务需求",
				EstimatedMinutes: 15,
				ToolsNeeded:      []string{string(models.ToolAnalyzeText)},
				Complexity:       "medium",
				Status:           "pending",
			},
			{
				StepNumber:       2,
				Description:      "制定执行计划",
				EstimatedMinutes: 30,
				ToolsNeeded:      []string{string(models.ToolPlanProject)},
				Complexity:       "medium",
				Status:           "pending",
			},
			{
				StepNumber:       3,
				Description:      "执行核心任务",
				EstimatedMinutes: 60,
				ToolsNeeded:      []string{},
				Complexity:       "medium",
				Status:           "pending",
			},
			{
				StepNumber:       4,
				Description:      "测试和优化",
				EstimatedMinutes: 30,
				ToolsNeeded:      []string{},
				Complexity:       "simple",
				Status:           "pending",
			},
		}

		for _, step := range steps {
			totalMinutes += step.EstimatedMinutes
		}
	}

	// 限制总时间不超过8小时
	if totalMinutes > 480 {
		scaleFactor := 480 / totalMinutes
		for i := range steps {
			steps[i].EstimatedMinutes *= scaleFactor
		}
		totalMinutes = 480
	}

	return steps, totalMinutes
}

// GetStatus 获取 Agent 状态
func (a *LLMAgent) GetStatus(ctx context.Context) *models.AgentStatus {
	ollamaOK, ollamaDetail := a.ollamaClient.CheckConnection(ctx)

	status := "ready"
	if !ollamaOK {
		status = "degraded"
	}

	tools := make([]string, 0, len(a.tools))
	for toolName := range a.tools {
		tools = append(tools, toolName)
	}

	return &models.AgentStatus{
		ModelName:      a.modelName,
		OllamaStatus:   map[bool]string{true: "connected", false: "disconnected"}[ollamaOK],
		OllamaDetail:   ollamaDetail,
		ToolsAvailable: tools,
		UptimeSeconds:  time.Since(a.startTime).Seconds(),
		RequestsCount:  atomic.LoadInt64(&a.requestCount),
		Status:         status,
		StartedAt:      a.startTime,
	}
}

// GetTools 获取工具列表
func (a *LLMAgent) GetTools() []models.ToolInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()

	tools := make([]models.ToolInfo, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, tool)
	}

	return tools
}

// Chat 聊天接口
func (a *LLMAgent) Chat(ctx context.Context, message string, conversationID string) (*models.ChatResponse, error) {
	atomic.AddInt64(&a.requestCount, 1)

	startTime := time.Now()

	response, err := a.ollamaClient.Chat(ctx, "你是一个有帮助的助手。", message)
	if err != nil {
		return nil, err
	}

	processingTime := int(time.Since(startTime).Milliseconds())

	return &models.ChatResponse{
		Response:       response,
		ConversationID: conversationID,
		MessageID:      generateMessageID(),
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}, nil
}

// GetTask 获取任务详情
func (a *LLMAgent) GetTask(taskID string) (*models.TaskDecompositionResponse, bool) {
	return a.taskManager.GetTask(taskID)
}

// QuickEstimate 快速时间估计
func (a *LLMAgent) QuickEstimate(ctx context.Context, task string, complexity string) map[string]interface{} {
	request := models.TaskDecompositionRequest{
		TaskDescription: task,
		MaxSteps:        3,
		IncludeTools:    false,
		ComplexityLevel: complexity,
	}

	response, err := a.DecomposeTask(ctx, request)
	if err != nil {
		return map[string]interface{}{
			"task":                 task,
			"estimated_minutes":    60,
			"formatted_duration":   "1小时",
			"estimated_completion": time.Now().Add(time.Hour).Format(time.RFC3339),
			"confidence":           "low",
			"error":                err.Error(),
		}
	}

	confidence := "medium"
	if len(response.Steps) > 2 {
		confidence = "high"
	} else if len(response.Steps) <= 1 {
		confidence = "low"
	}

	return map[string]interface{}{
		"task":                 task,
		"estimated_minutes":    response.TotalMinutes,
		"formatted_duration":   formatDuration(response.TotalMinutes),
		"estimated_completion": response.EstimatedCompletion.Format(time.RFC3339),
		"confidence":           confidence,
	}
}

// 辅助函数
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

func formatDuration(minutes float64) string {
	if minutes < 60 {
		return fmt.Sprintf("%.1f分钟", minutes)
	} else if minutes < 1440 {
		hours := minutes / 60
		return fmt.Sprintf("%.1f小时", hours)
	} else {
		days := minutes / 1440
		return fmt.Sprintf("%.1f天", days)
	}
}
