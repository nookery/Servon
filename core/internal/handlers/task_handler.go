package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleListTasks 处理获取任务列表的请求
func HandleListTasks(c *gin.Context) {
	taskList := []string{}
	tasks := DefaultTaskManager.GetTasks()
	for _, task := range tasks {
		taskList = append(taskList, task.ID)
	}

	c.JSON(http.StatusOK, taskList)
}

// HandleRemoveTask 处理删除任务的请求
func HandleRemoveTask(c *gin.Context) {
	id := c.Param("id")
	DefaultTaskManager.RemoveTask(id)
	c.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}

// HandleExecuteTask 处理执行任务的请求
func HandleExecuteTask(c *gin.Context) {
	id := c.Param("id")
	err := DefaultTaskManager.ExecuteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "任务已执行"})
}
