package controllers

import (
	"net/http"
	"servon/core/managers"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	*managers.FullManager
}

func NewTaskController(manager *managers.FullManager) *TaskController {
	return &TaskController{FullManager: manager}
}

// HandleListTasks 处理获取任务列表的请求
func (h *TaskController) HandleListTasks(c *gin.Context) {
	taskList := []string{}
	tasks := h.GetTasks()
	for _, task := range tasks {
		taskList = append(taskList, task.ID)
	}

	c.JSON(http.StatusOK, taskList)
}

// HandleRemoveTask 处理删除任务的请求
func (h *TaskController) HandleRemoveTask(c *gin.Context) {
	id := c.Param("id")
	h.RemoveTask(id)
	c.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}

// HandleExecuteTask 处理执行任务的请求
func (h *TaskController) HandleExecuteTask(c *gin.Context) {
	id := c.Param("id")
	err := h.ExecuteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "任务已执行"})
}
