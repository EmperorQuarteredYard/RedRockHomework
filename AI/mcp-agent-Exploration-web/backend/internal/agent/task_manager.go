package agent

import (
	"backend/internal/models"
	"sync"
	"time"
)

// TaskManager 任务管理器
type TaskManager struct {
	tasks map[string]*models.TaskDecompositionResponse
	mu    sync.RWMutex
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*models.TaskDecompositionResponse),
	}
}

// SaveTask 保存任务
func (tm *TaskManager) SaveTask(task *models.TaskDecompositionResponse) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 限制任务数量，最多保存1000个任务
	if len(tm.tasks) >= 1000 {
		// 删除最旧的任务
		var oldestKey string
		var oldestTime time.Time = time.Now()

		for key, task := range tm.tasks {
			if task.CreatedAt.Before(oldestTime) {
				oldestTime = task.CreatedAt
				oldestKey = key
			}
		}

		if oldestKey != "" {
			delete(tm.tasks, oldestKey)
		}
	}

	tm.tasks[task.TaskID] = task
}

// GetTask 获取任务
func (tm *TaskManager) GetTask(taskID string) (*models.TaskDecompositionResponse, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	task, exists := tm.tasks[taskID]
	return task, exists
}

// ListTasks 列出所有任务
func (tm *TaskManager) ListTasks(limit, offset int) []*models.TaskDecompositionResponse {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var tasks []*models.TaskDecompositionResponse
	count := 0

	// 按创建时间排序（这里简化处理，实际可能需要更好的排序）
	for _, task := range tm.tasks {
		if count >= offset && len(tasks) < limit {
			tasks = append(tasks, task)
		}
		count++
	}

	return tasks
}

// DeleteTask 删除任务
func (tm *TaskManager) DeleteTask(taskID string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tasks[taskID]; exists {
		delete(tm.tasks, taskID)
		return true
	}

	return false
}

// CleanupOldTasks 清理旧任务
func (tm *TaskManager) CleanupOldTasks(maxAge time.Duration) int {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	deletedCount := 0

	for taskID, task := range tm.tasks {
		if task.CreatedAt.Before(cutoff) {
			delete(tm.tasks, taskID)
			deletedCount++
		}
	}

	return deletedCount
}
