package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlePortList 处理获取端口列表的请求
func (h *WebHandler) HandlePortList(c *gin.Context) {
	ports, err := h.App.PortManager.GetPortList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取端口列表失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ports)
}
