package handlers

import (
	"net/http"
	"servon/core"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandleListCronTasks 处理获取所有定时任务的请求
func (h *WebHandler) HandleListCronTasks(c *gin.Context) {
	tasks, err := h.App.CronManager.GetCronTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// HandleCreateCronTask 处理创建定时任务的请求
func (h *WebHandler) HandleCreateCronTask(c *gin.Context) {
	var task core.CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	newTask, err := h.App.CronManager.CreateCronTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newTask)
}

// HandleUpdateCronTask 处理更新定时任务的请求
func (h *WebHandler) HandleUpdateCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var task core.CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	task.ID = id

	updatedTask, err := h.App.CronManager.UpdateCronTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// HandleDeleteCronTask 处理删除定时任务的请求
func (h *WebHandler) HandleDeleteCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := h.App.CronManager.DeleteCronTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleToggleCronTask 处理启用/禁用定时任务的请求
func (h *WebHandler) HandleToggleCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	task, err := h.App.CronManager.ToggleCronTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
