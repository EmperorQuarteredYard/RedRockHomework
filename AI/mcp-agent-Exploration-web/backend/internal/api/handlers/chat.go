package handlers

import (
	"backend/internal/agent"
	"backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	agent  *agent.LLMAgent
	logger *zap.Logger
}

// NewChatHandler 创建新的聊天处理器
func NewChatHandler(agent *agent.LLMAgent, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		agent:  agent,
		logger: logger,
	}
}

// Chat 聊天接口
// @Summary 与Agent对话
// @Description 与LLM Agent进行对话交互
// @Tags chat
// @Accept json
// @Produce json
// @Param request body models.ChatRequest true "聊天请求"
// @Success 200 {object} models.ChatResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/chat [post]
func (h *ChatHandler) Chat(c *gin.Context) {
	var request models.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("聊天请求参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "无效的请求参数: " + err.Error(),
		})
		return
	}

	h.logger.Info("处理聊天请求",
		zap.String("message", request.Message[:min(len(request.Message), 100)]),
		zap.String("conversation_id", request.ConversationID),
	)

	response, err := h.agent.Chat(c.Request.Context(), request.Message, request.ConversationID)
	if err != nil {
		h.logger.Error("聊天请求失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "聊天失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
