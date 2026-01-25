package handlers

import (
	"backend/internal/models"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler 创建新的健康检查处理器
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthCheckResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// 获取系统信息
	var cpuUsage float64
	var memUsage float64

	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	if memStat, err := mem.VirtualMemory(); err == nil {
		memUsage = memStat.UsedPercent
	}

	uptime := time.Since(h.startTime)

	response := models.HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    formatUptime(uptime),
	}

	// 添加系统信息
	systemInfo := map[string]interface{}{
		"go_version":    runtime.Version(),
		"num_goroutine": runtime.NumGoroutine(),
		"cpu_usage":     cpuUsage,
		"memory_usage":  memUsage,
		"num_cpu":       runtime.NumCPU(),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"health": response,
			"system": systemInfo,
		},
	})
}

// formatUptime 格式化运行时间
func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分%d秒", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	}
	return fmt.Sprintf("%d秒", seconds)
}
