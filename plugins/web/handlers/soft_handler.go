package handlers

import (
	"net/http"
	"servon/core"
	"sort"

	"github.com/gin-gonic/gin"
)

// HandleGetSoftwareList 处理获取软件列表的请求
func (h *WebHandler) HandleGetSoftwareList(c *gin.Context) {
	softwareList := h.App.SoftManager.GetAllSoftware()
	sort.Strings(softwareList)
	c.JSON(200, softwareList)
}

// HandleInstallSoftware 处理安装软件的请求
func (h *WebHandler) HandleInstallSoftware(c *gin.Context) {
	name := c.Param("name")

	// 发送初始消息
	c.JSON(http.StatusOK, gin.H{"message": "正在准备安装..."})

	h.App.TaskManager.AddTaskAndExecute(core.Task{
		ID: name,
		Execute: func() error {
			return h.App.SoftManager.Install(name)
		},
	}, "SoftManager.HandleInstallSoftware")

	c.JSON(http.StatusAccepted, gin.H{"message": "安装请求已接受，正在后台处理..."})
}

// HandleUninstallSoftware 处理卸载软件的请求
func (h *WebHandler) HandleUninstallSoftware(c *gin.Context) {
	name := c.Param("name")
	err := h.App.SoftManager.UninstallSoftware(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

// HandleStopSoftware 处理停止软件的请求
func (h *WebHandler) HandleStopSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := h.App.SoftManager.StopSoftware(name); err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

// HandleGetSoftwareStatus 处理获取软件状态的请求
func (h *WebHandler) HandleGetSoftwareStatus(c *gin.Context) {
	name := c.Param("name")
	status, err := h.App.SoftManager.GetSoftwareStatus(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, status)
}

// HandleStartSoftware 处理启动软件的请求
func (h *WebHandler) HandleStartSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := h.App.SoftManager.StartSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
