package handlers

import (
	"net/http"
	"servon/core/internal/libs"
	"strconv"

	"github.com/gin-gonic/gin"
)

var DefaultCronManager = libs.DefaultCronManager
type CronTask = libs.CronTask

// HandleListCronTasks 处理获取所有定时任务的请求
func HandleListCronTasks(c *gin.Context) {
	tasks, err := DefaultCronManager.GetCronTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// HandleCreateCronTask 处理创建定时任务的请求
func HandleCreateCronTask(c *gin.Context) {
	var task CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	newTask, err := DefaultCronManager.CreateCronTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newTask)
}

// HandleUpdateCronTask 处理更新定时任务的请求
func HandleUpdateCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var task CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	task.ID = id

	updatedTask, err := DefaultCronManager.UpdateCronTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// HandleDeleteCronTask 处理删除定时任务的请求
func HandleDeleteCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := DefaultCronManager.DeleteCronTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleToggleCronTask 处理启用/禁用定时任务的请求
func HandleToggleCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	task, err := DefaultCronManager.ToggleCronTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
