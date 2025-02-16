package serve

import (
	"net/http"
	"servon/core"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CronTask = core.CronTask

// HandleKillProcess 处理结束进程的请求
func (p *ServePlugin) HandleKillProcess(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的PID"})
		return
	}

	err = p.KillProcess(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
