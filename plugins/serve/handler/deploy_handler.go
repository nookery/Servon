package handler

import (
	"io"
	"net/http"
	"strconv"

	"servon/plugins/deploy"

	"github.com/gin-gonic/gin"
)

type DeployHandler struct{}

func NewDeployHandler() *DeployHandler {
	return &DeployHandler{}
}

// HandleListProjects 获取所有项目
func (h *DeployHandler) HandleListProjects(c *gin.Context) {
	projects, err := deploy.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// HandleCreateProject 创建项目
func (h *DeployHandler) HandleCreateProject(c *gin.Context) {
	var project deploy.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	newProject, err := deploy.CreateProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newProject)
}

// HandleUpdateProject 更新项目
func (h *DeployHandler) HandleUpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	var project deploy.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	project.ID = id

	updatedProject, err := deploy.UpdateProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedProject)
}

// HandleDeleteProject 删除项目
func (h *DeployHandler) HandleDeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	if err := deploy.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "项目已删除"})
}

// HandleBuildProject 构建项目
func (h *DeployHandler) HandleBuildProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	// 使用 Server-Sent Events
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 开始构建并发送日志
	logChan := make(chan string)
	errChan := make(chan error)

	go func() {
		err := deploy.BuildProject(id, logChan)
		if err != nil {
			errChan <- err
		}
		close(logChan)
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case log, ok := <-logChan:
			if !ok {
				return false
			}
			c.SSEvent("message", log)
			return true
		case err := <-errChan:
			c.SSEvent("error", err.Error())
			return false
		}
	})
}

// HandleProjectLogs 获取项目日志
func (h *DeployHandler) HandleProjectLogs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	logs, err := deploy.GetProjectLogs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
