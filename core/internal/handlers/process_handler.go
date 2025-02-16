package handlers

import (
	"net/http"
	"servon/core/internal/libs"

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
