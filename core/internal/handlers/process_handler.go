package handlers

import (
	"net/http"
	"servon/core/internal/libs"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ProcessUtil = libs.DefaultProcessManager

// HandleProcessList 处理获取进程列表的请求
func HandleProcessList(c *gin.Context) {
	processes, err := ProcessUtil.GetProcessList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

// HandleKillProcess 处理结束进程的请求
func HandleKillProcess(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的PID"})
		return
	}

	err = ProcessUtil.KillProcess(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
