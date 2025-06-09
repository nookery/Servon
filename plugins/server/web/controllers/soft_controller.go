package controllers

import (
	"net/http"
	"servon/core/managers"
	"sort"

	"github.com/gin-gonic/gin"
)

type SoftController struct {
	*managers.FullManager
}

func NewSoftController(manager *managers.FullManager) *SoftController {
	return &SoftController{FullManager: manager}
}

// HandleGetSoftwareList 处理获取软件列表的请求
func (h *SoftController) HandleGetSoftwareList(c *gin.Context) {
	softwareList := h.GetAllSoftware()
	sort.Strings(softwareList)
	c.JSON(200, softwareList)
}

// HandleInstallSoftware 处理安装软件的请求
func (h *SoftController) HandleInstallSoftware(c *gin.Context) {
	name := c.Param("name")

	// 发送初始消息
	c.JSON(http.StatusOK, gin.H{"message": "正在准备安装..."})

	h.AddTaskAndExecute(managers.Task{
		ID: name,
		Execute: func() error {
			return h.Install(name)
		},
	}, "SoftManager.HandleInstallSoftware")

	c.JSON(http.StatusAccepted, gin.H{"message": "安装请求已接受，正在后台处理..."})
}

// HandleUninstallSoftware 处理卸载软件的请求
func (h *SoftController) HandleUninstallSoftware(c *gin.Context) {
	name := c.Param("name")
	err := h.UninstallSoftware(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

// HandleStopSoftware 处理停止软件的请求
func (h *SoftController) HandleStopSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := h.StopSoftware(name); err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

// HandleGetSoftwareStatus 处理获取软件状态的请求
func (h *SoftController) HandleGetSoftwareStatus(c *gin.Context) {
	name := c.Param("name")
	status, err := h.GetSoftwareStatus(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, status)
}

// HandleStartSoftware 处理启动软件的请求
func (h *SoftController) HandleStartSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := h.StartSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
