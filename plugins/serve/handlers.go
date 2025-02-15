package serve

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servon/core"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CronTask = core.CronTask

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
	p.PrintInfo("获取端口列表")
	ports, err := p.GetPortList()
	if err != nil {
		p.PrintErrorMessage("获取端口列表失败: " + err.Error())
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

	// 发送初始消息
	c.JSON(http.StatusOK, gin.H{"message": "正在准备安装..."})

	p.AddTaskAndExecute(core.Task{
		ID: name,
		Execute: func() error {
			return p.Install(name)
		},
	}, "ServePlugin.HandleInstallSoftware")

	c.JSON(http.StatusAccepted, gin.H{"message": "安装请求已接受，正在后台处理..."})
}

// HandleUninstallSoftware 处理卸载软件的请求
func (p *ServePlugin) HandleUninstallSoftware(c *gin.Context) {
	name := c.Param("name")
	err := p.UninstallSoftware(name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
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

// HandleListCronTasks 处理获取所有定时任务的请求
func (p *ServePlugin) HandleListCronTasks(c *gin.Context) {
	tasks, err := p.GetCronTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// HandleCreateCronTask 处理创建定时任务的请求
func (p *ServePlugin) HandleCreateCronTask(c *gin.Context) {
	var task CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	newTask, err := p.CreateCronTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newTask)
}

// HandleUpdateCronTask 处理更新定时任务的请求
func (p *ServePlugin) HandleUpdateCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var task CronTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	task.ID = id

	updatedTask, err := p.UpdateCronTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// HandleDeleteCronTask 处理删除定时任务的请求
func (p *ServePlugin) HandleDeleteCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := p.DeleteCronTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleToggleCronTask 处理启用/禁用定时任务的请求
func (p *ServePlugin) HandleToggleCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	task, err := p.ToggleCronTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// HandleFileDownload handles file download requests
func (p *ServePlugin) HandleFileDownload(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}

	// Verify the file exists and is not a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot download directories"})
		return
	}

	// Serve the file
	c.File(path)
}

// HandleFileContent 处理获取文件内容的请求
func (p *ServePlugin) HandleFileContent(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供文件路径"})
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

// HandleSaveFile 处理保存文件内容的请求
func (p *ServePlugin) HandleSaveFile(c *gin.Context) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := os.WriteFile(req.Path, []byte(req.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleDeleteFile 处理删除文件的请求
func (p *ServePlugin) HandleDeleteFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供文件路径"})
		return
	}

	err := os.Remove(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleCreateFile 处理创建新文件的请求
func (p *ServePlugin) HandleCreateFile(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
		Type string `json:"type"` // "file" 或 "directory"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	// 检查路径是否已存在
	if _, err := os.Stat(req.Path); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件或目录已存在: " + req.Path})
		return
	}

	// 检查父目录是否存在且可写
	parentDir := filepath.Dir(req.Path)
	if _, err := os.Stat(parentDir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "父目录不存在: " + parentDir})
		return
	}

	if req.Type == "directory" {
		err := os.MkdirAll(req.Path, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("创建目录失败: %v (路径: %s)", err, req.Path),
			})
			return
		}
	} else {
		// 创建空文件
		f, err := os.Create(req.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("创建文件失败: %v (路径: %s)", err, req.Path),
			})
			return
		}
		f.Close()
	}

	c.Status(http.StatusOK)
}

// HandleListUsers 处理获取用户列表的请求
func (p *ServePlugin) HandleListUsers(c *gin.Context) {
	users, err := p.GetUserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// HandleCreateUser 处理创建用户的请求
func (p *ServePlugin) HandleCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := p.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleDeleteUser 处理删除用户的请求
func (p *ServePlugin) HandleDeleteUser(c *gin.Context) {
	username := c.Param("username")
	err := p.DeleteUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleStartSoftware 处理启动软件的请求
func (p *ServePlugin) HandleStartSoftware(c *gin.Context) {
	name := c.Param("name")
	if err := p.StartSoftware(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleListTasks 处理获取任务列表的请求
func (p *ServePlugin) HandleListTasks(c *gin.Context) {
	p.PrintInfo("获取任务列表")
	taskList := []string{}
	tasks := p.GetTasks()
	for _, task := range tasks {
		taskList = append(taskList, task.ID)
	}

	c.JSON(http.StatusOK, taskList)
}

// HandleRemoveTask 处理删除任务的请求
func (p *ServePlugin) HandleRemoveTask(c *gin.Context) {
	id := c.Param("id")
	p.RemoveTask(id)
	c.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}

// HandleExecuteTask 处理执行任务的请求
func (p *ServePlugin) HandleExecuteTask(c *gin.Context) {
	id := c.Param("id")
	err := p.ExecuteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "任务已执行"})
}

// GitHubManifest represents the GitHub App manifest configuration
type GitHubManifest struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	HookAttributes struct {
		URL    string `json:"url"`
		Active bool   `json:"active"`
	} `json:"hook_attributes"`
	RedirectURL        string            `json:"redirect_url"`
	CallbackURLs       []string          `json:"callback_urls"`
	Description        string            `json:"description"`
	Public             bool              `json:"public"`
	DefaultEvents      []string          `json:"default_events"`
	DefaultPermissions map[string]string `json:"default_permissions"`
}

// HandleGitHubSetup 处理GitHub App Manifest flow的设置请求
func (p *ServePlugin) HandleGitHubSetup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供必要的GitHub App信息"})
		return
	}

	baseURL := fmt.Sprintf("http://%s:%d", p.Config.Host, p.Config.Port)

	// 构建manifest
	manifest := GitHubManifest{
		Name:        req.Name,
		URL:         baseURL,
		Description: req.Description,
		Public:      true,
		HookAttributes: struct {
			URL    string `json:"url"`
			Active bool   `json:"active"`
		}{
			URL:    fmt.Sprintf("%s/web_api/github/webhook", baseURL),
			Active: true,
		},
		RedirectURL:  fmt.Sprintf("%s/web_api/github/callback", baseURL),
		CallbackURLs: []string{fmt.Sprintf("%s/web_api/github/callback", baseURL)},
		DefaultPermissions: map[string]string{
			"issues": "write",
			"checks": "write",
		},
		DefaultEvents: []string{
			"issues",
			"issue_comment",
			"check_suite",
			"check_run",
		},
	}

	// 将manifest转换为JSON
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成manifest失败"})
		return
	}

	// 生成state参数
	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成state失败"})
		return
	}
	stateStr := hex.EncodeToString(state)

	// 生成HTML表单
	html := fmt.Sprintf(`
		<form id="github-form" action="https://github.com/settings/apps/new?state=%s" method="post">
			<input type="hidden" name="manifest" value='%s'>
		</form>
		<script>document.getElementById("github-form").submit();</script>
	`, stateStr, string(manifestJSON))

	// 返回HTML
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// HandleGitHubCallback 处理GitHub App创建后的回调
func (p *ServePlugin) HandleGitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	// 调用GitHub API完成App创建
	resp, err := http.Post(
		fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code),
		"application/json",
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create GitHub App: %v", err)})
		return
	}
	defer resp.Body.Close()

	var result struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		WebhookURL string `json:"webhook_url"`
		PEM        string `json:"pem"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse GitHub response: %v", err)})
		return
	}

	// 保存App配置
	p.Config.GitHubAppID = result.ID
	p.Config.GitHubAppPrivateKey = result.PEM
	// TODO: 保存配置到文件

	c.JSON(http.StatusOK, gin.H{
		"message": "GitHub App created successfully",
		"app_id":  result.ID,
		"name":    result.Name,
	})
}

// HandleGitHubWebhook 处理来自GitHub的webhook请求
func (p *ServePlugin) HandleGitHubWebhook(c *gin.Context) {
	// 验证webhook签名
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	// 读取请求体
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}

	// 验证webhook secret
	// TODO: 实现签名验证

	// 解析事件类型
	event := c.GetHeader("X-GitHub-Event")
	switch event {
	case "push":
		// 处理push事件
		var pushEvent struct {
			Ref        string `json:"ref"`
			Before     string `json:"before"`
			After      string `json:"after"`
			Created    bool   `json:"created"`
			Deleted    bool   `json:"deleted"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		}
		if err := json.Unmarshal(payload, &pushEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid push event payload"})
			return
		}

		p.PrintInfo(fmt.Sprintf("Received push event for %s", pushEvent.Repository.FullName))
		// TODO: 处理push事件的具体逻辑

	case "pull_request":
		// 处理pull request事件
		var prEvent struct {
			Action      string `json:"action"`
			PullRequest struct {
				Number int    `json:"number"`
				Title  string `json:"title"`
				State  string `json:"state"`
			} `json:"pull_request"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		}
		if err := json.Unmarshal(payload, &prEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pull request event payload"})
			return
		}

		p.PrintInfo(fmt.Sprintf("Received PR event for %s: %s", prEvent.Repository.FullName, prEvent.Action))
		// TODO: 处理PR事件的具体逻辑

	default:
		p.PrintInfo(fmt.Sprintf("Received unhandled event type: %s", event))
	}

	c.Status(http.StatusOK)
}
