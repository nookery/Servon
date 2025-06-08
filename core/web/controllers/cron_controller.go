package controllers

import (
	"net/http"
	"servon/components/cron_util"
	"servon/core/managers"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CronController struct {
	*managers.FullManager
}

func NewCronController(manager *managers.FullManager) *CronController {
	return &CronController{
		FullManager: manager,
	}
}

// HandleListCronTasks 处理获取所有定时任务的请求
func (c *CronController) HandleListCronTasks(ctx *gin.Context) {
	tasks, err := c.GetCronTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

// HandleCreateCronTask 处理创建定时任务的请求
func (c *CronController) HandleCreateCronTask(ctx *gin.Context) {
	var task cron_util.CronTask
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	newTask, err := c.CreateCronTask(task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, newTask)
}

// HandleUpdateCronTask 处理更新定时任务的请求
func (c *CronController) HandleUpdateCronTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var task cron_util.CronTask
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	task.ID = id

	updatedTask, err := c.UpdateCronTask(task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedTask)
}

// HandleDeleteCronTask 处理删除定时任务的请求
func (c *CronController) HandleDeleteCronTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := c.DeleteCronTask(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

// HandleToggleCronTask 处理启用/禁用定时任务的请求
func (c *CronController) HandleToggleCronTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	task, err := c.ToggleCronTask(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}
