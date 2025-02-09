package handler

import (
	"sort"

	"servon/core"

	"github.com/gin-gonic/gin"
)

// SoftwareHandler 处理软件相关的 HTTP 请求
type SoftwareHandler struct {
	core *core.Core
}

// NewSoftwareHandler 创建软件处理器实例
func NewSoftwareHandler(core *core.Core) *SoftwareHandler {
	if core == nil {
		panic("core is nil")
	}

	return &SoftwareHandler{
		core: core,
	}
}

// HandleGetSoftwareList 处理获取软件列表的请求
func (h *SoftwareHandler) HandleGetSoftwareList(c *gin.Context) {
	softwareList := h.core.GetAllSoftware()
	names := make([]string, len(softwareList))
	sort.Strings(names)
	c.JSON(200, names)
}

// HandleInstallSoftware 处理安装软件的请求
func (h *SoftwareHandler) HandleInstallSoftware(c *gin.Context) {
	name := c.Param("name")
	msgChan := make(chan string, 100)
	err := h.core.InstallSoftware(name, msgChan)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for msg := range msgChan {
		c.SSEvent("message", msg)
		c.Writer.Flush()
	}
}

// HandleUninstallSoftware 处理卸载软件的请求
func (h *SoftwareHandler) HandleUninstallSoftware(c *gin.Context) {
	name := c.Param("name")
	msgChan := make(chan string, 100)
	err := h.core.UninstallSoftware(name, msgChan)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for msg := range msgChan {
		c.SSEvent("message", msg)
		c.Writer.Flush()
	}
}

// HandleStopSoftware 处理停止软件的请求
func (h *SoftwareHandler) HandleStopSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := h.core.StopSoftware(name); err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

// HandleGetSoftwareStatus 处理获取软件状态的请求
func (h *SoftwareHandler) HandleGetSoftwareStatus(c *gin.Context) {
	name := c.Param("name")
	status, err := h.core.GetSoftwareStatus(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, status)
}
