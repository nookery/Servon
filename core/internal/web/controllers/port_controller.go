package controllers

import (
	"net/http"
	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

type PortController struct {
	*managers.FullManager
}

func NewPortController(manager *managers.FullManager) *PortController {
	return &PortController{FullManager: manager}
}

// HandlePortList 处理获取端口列表的请求
func (h *PortController) HandlePortList(c *gin.Context) {
	ports, err := h.GetPortList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取端口列表失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ports)
}
