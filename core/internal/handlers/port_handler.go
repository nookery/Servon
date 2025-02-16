package handlers

import (
	"net/http"
	"servon/core/internal/libs"

	"github.com/gin-gonic/gin"
)

// HandlePortList 处理获取端口列表的请求
func HandlePortList(c *gin.Context) {
	ports, err := libs.DefaultPortManager.GetPortList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取端口列表失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ports)
}
