package handler

import (
	"fmt"
	"net/http"
	"servon/cmd/system"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CronTask 定时任务结构
type CronTask struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Command     string    `json:"command"`
	Schedule    string    `json:"schedule"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	LastRun     time.Time `json:"last_run,omitempty"`
	NextRun     time.Time `json:"next_run,omitempty"`
}

// 将 handler 中的 CronTask 转换为 system.CronTask
func (t *CronTask) toSystemTask() system.CronTask {
	return system.CronTask{
		ID:          t.ID,
		Name:        t.Name,
		Command:     t.Command,
		Schedule:    t.Schedule,
		Description: t.Description,
		Enabled:     t.Enabled,
		LastRun:     t.LastRun,
		NextRun:     t.NextRun,
	}
}

// 将 system.CronTask 转换为 handler 中的 CronTask
func fromSystemTask(t *system.CronTask) *CronTask {
	return &CronTask{
		ID:          t.ID,
		Name:        t.Name,
		Command:     t.Command,
		Schedule:    t.Schedule,
		Description: t.Description,
		Enabled:     t.Enabled,
		LastRun:     t.LastRun,
		NextRun:     t.NextRun,
	}
}

// 添加一个辅助函数来处理错误响应
func handleError(c *gin.Context, err error) {
	if ve, ok := err.(system.ValidationErrors); ok {
		// 处理验证错误
		c.JSON(http.StatusBadRequest, ve)
	} else {
		// 处理其他错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": []system.ValidationError{
				{
					Field:   "general",
					Message: err.Error(),
				},
			},
		})
	}
}

// HandleListCronTasks 获取所有定时任务
func (h *Handler) HandleListCronTasks(c *gin.Context) {
	tasks, err := system.GetCronTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换任务列表
	result := make([]*CronTask, len(tasks))
	for i, task := range tasks {
		result[i] = fromSystemTask(task)
	}
	c.JSON(http.StatusOK, result)
}

// HandleCreateCronTask 创建定时任务
func (h *Handler) HandleCreateCronTask(c *gin.Context) {
	var task CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		handleError(c, fmt.Errorf("无效的请求数据"))
		return
	}

	systemTask := task.toSystemTask()
	newTask, err := system.CreateCronTask(systemTask)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, fromSystemTask(newTask))
}

// HandleUpdateCronTask 更新定时任务
func (h *Handler) HandleUpdateCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, fmt.Errorf("无效的任务ID"))
		return
	}

	var task CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		handleError(c, fmt.Errorf("无效的请求数据"))
		return
	}
	task.ID = id

	systemTask := task.toSystemTask()
	updatedTask, err := system.UpdateCronTask(systemTask)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, fromSystemTask(updatedTask))
}

// HandleDeleteCronTask 删除定时任务
func (h *Handler) HandleDeleteCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, fmt.Errorf("无效的任务ID"))
		return
	}

	if err := system.DeleteCronTask(id); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}

// HandleToggleCronTask 启用/禁用定时任务
func (h *Handler) HandleToggleCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, fmt.Errorf("无效的任务ID"))
		return
	}

	task, err := system.ToggleCronTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, fromSystemTask(task))
}
