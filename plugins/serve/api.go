package serve

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleSystemResources 处理系统资源监控的请求
func (p *ServePlugin) HandleSystemResources(c *gin.Context) {
	resources, err := p.GetSystemResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// HandleBasicInfo 处理基本系统信息的请求
func (p *ServePlugin) HandleBasicInfo(c *gin.Context) {
	info, err := p.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// HandleCurrentUser 处理获取当前用户的请求
func (p *ServePlugin) HandleCurrentUser(c *gin.Context) {
	user, err := p.GetCurrentUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user})
}

// HandleProcessList 处理获取进程列表的请求
func (p *ServePlugin) HandleProcessList(c *gin.Context) {
	processes, err := p.GetProcessList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

// HandleFileList 处理获取文件列表的请求
func (p *ServePlugin) HandleFileList(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := p.GetFileList(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文件列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

// HandlePortList 处理获取端口列表的请求
func (p *ServePlugin) HandlePortList(c *gin.Context) {
	ports, err := p.GetPortList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取端口列表失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ports)
}

// HandleOSInfo 处理获取操作系统信息的请求
func (p *ServePlugin) HandleOSInfo(c *gin.Context) {
	osInfo, err := p.GetOSInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"os_info": osInfo})
}

// HandleNetworkResources 处理网络资源监控的请求
func (p *ServePlugin) HandleNetworkResources(c *gin.Context) {
	networkStats, err := p.GetNetworkResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, networkStats)
}

// HandleGetSoftwareList 处理获取软件列表的请求
func (p *ServePlugin) HandleGetSoftwareList(c *gin.Context) {
	softwareList := p.GetAllSoftware()
	sort.Strings(softwareList)
	c.JSON(200, softwareList)
}

// HandleInstallSoftware 处理安装软件的请求
func (p *ServePlugin) HandleInstallSoftware(c *gin.Context) {
	name := c.Param("name")
	msgChan := make(chan string, 100)
	doneChan := make(chan error, 1)

	// 设置 SSE 头信息
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 发送初始消息
	c.SSEvent("message", map[string]string{
		"type":    "log",
		"message": "正在准备安装...",
	})
	c.Writer.Flush()

	// 在新的 goroutine 中执行安装
	go func() {
		err := p.Install(name, msgChan)
		doneChan <- err
		close(msgChan)
	}()

	// 实时发送消息到客户端
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				// msgChan 已关闭，等待 doneChan
				continue
			}
			c.SSEvent("message", map[string]string{
				"type":    "log",
				"message": msg,
			})
			// 立即刷新缓冲区，确保消息实时发送
			c.Writer.Flush()
		case err := <-doneChan:
			if err != nil {
				c.SSEvent("message", map[string]string{
					"type":    "error",
					"message": err.Error(),
				})
			} else {
				c.SSEvent("message", map[string]string{
					"type": "complete",
				})
			}
			c.Writer.Flush()
			return
		}
	}
}

// HandleUninstallSoftware 处理卸载软件的请求
func (p *ServePlugin) HandleUninstallSoftware(c *gin.Context) {
	name := c.Param("name")
	msgChan := make(chan string, 100)
	err := p.UninstallSoftware(name, msgChan)
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
func (p *ServePlugin) HandleStopSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := p.StopSoftware(name); err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

// HandleGetSoftwareStatus 处理获取软件状态的请求
func (p *ServePlugin) HandleGetSoftwareStatus(c *gin.Context) {
	name := c.Param("name")
	status, err := p.GetSoftwareStatus(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, status)
}

// CronTask 定时任务结构
type CronTask struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Command     string    `json:"command"`
	Schedule    string    `json:"schedule"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	LastRun     time.Time `json:"last_run,omitempty"`
	NextRun     time.Time `json:"next_run,omitempty"`
}
