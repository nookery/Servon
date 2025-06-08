package controllers

import (
	"net/http"
	"servon/core/managers"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProcessController struct {
	*managers.FullManager
}

func NewProcessController(manager *managers.FullManager) *ProcessController {
	return &ProcessController{FullManager: manager}
}

// HandleProcessList 处理获取进程列表的请求
func (h *ProcessController) HandleProcessList(c *gin.Context) {
	processes, err := h.GetProcessList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

// HandleKillProcess 处理结束进程的请求
func (h *ProcessController) HandleKillProcess(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的PID"})
		return
	}

	err = h.KillProcess(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
