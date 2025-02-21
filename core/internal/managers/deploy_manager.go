package managers

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"servon/core/internal/events"
	"servon/core/internal/models"
	"servon/core/internal/utils"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// DeployManager 是部署管理器，负责处理项目的自动化部署流程
// 主要功能包括：
// - 监听代码仓库的推送事件
// - 执行项目的自动化部署
// - 记录和管理部署日志
// - 发布部署相关的事件通知
type DeployManager struct {
	// eventBus 用于处理事件的发布与订阅
	eventBus *events.EventBus
	// logUtil 用于记录部署过程的日志
	logUtil *utils.LogUtil
	// gitUtil 用于处理Git操作
	gitUtil *utils.GitUtil
	// logDir 指定部署日志文件的存储目录
	logDir string
}

func NewDeployManager(eventBus *events.EventBus) (*DeployManager, error) {
	logDir := "/data/deploy"
	logUtil, err := utils.NewLogUtil(logDir)
	if err != nil {
		return nil, fmt.Errorf("初始化日志工具失败: %v", err)
	}

	dm := &DeployManager{
		eventBus: eventBus,
		logUtil:  logUtil,
		gitUtil:  utils.NewGitUtil(),
		logDir:   logDir,
	}

	// 订阅Git Push事件
	eventBus.Subscribe(events.GitPush, dm.handleGitPushEvent)

	return dm, nil
}

// handleGitPushEvent 处理Git Push事件
func (m *DeployManager) handleGitPushEvent(event events.Event) {
	deployData, ok := event.Data.(map[string]interface{})
	if !ok {
		m.logUtil.Log("错误", "无效的部署数据格式")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		m.logUtil.Log("错误", "缺少仓库信息")
		return
	}

	// 执行部署操作
	if err := m.DeployProject(repo); err != nil {
		m.logUtil.LogToFile(utils.LogLevelError, nil, "错误: 仓库 %s 部署失败: %v", repo, err)

		// 发布部署失败事件
		m.eventBus.Publish(events.Event{
			Type: events.DeployFailed,
			Data: map[string]interface{}{
				"repository": repo,
				"error":      err.Error(),
			},
		})
		return
	}

	m.logUtil.Log("信息", "仓库 %s 部署成功完成", repo)
	// 发布部署成功事件
	m.eventBus.Publish(events.Event{
		Type: events.DeployComplete,
		Data: map[string]interface{}{
			"repository": repo,
			"status":     "success",
		},
	})
}

// GetDeployLog 获取单个部署日志
func (m *DeployManager) GetDeployLog(logID string) (*models.DeployLog, error) {
	logPath := filepath.Join(m.logDir, logID+".log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		m.logUtil.Log("错误", "读取日志文件 %s 失败: %v", logID, err)
		return nil, fmt.Errorf("读取日志文件失败: %v", err)
	}

	timestamp, err := time.Parse("2006-01-02", logID)
	if err != nil {
		m.logUtil.Log("错误", "无效的日志ID格式 %s: %v", logID, err)
		return nil, fmt.Errorf("无效的日志ID格式: %v", err)
	}

	return &models.DeployLog{
		ID:        logID,
		Timestamp: timestamp,
		Message:   string(content), // 将内容放在 Message 字段
		Status:    "success",
	}, nil
}

// ListDeployLogs 获取部署日志列表
func (m *DeployManager) ListDeployLogs() ([]models.DeployLog, error) {
	files, err := os.ReadDir(m.logDir)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	logs := make([]models.DeployLog, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			// 解析日期格式的文件名 (2025-02-19.log)
			dateStr := strings.TrimSuffix(file.Name(), ".log")
			timestamp, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue // 跳过无法解析的文件名
			}

			// 读取文件内容
			content, err := os.ReadFile(filepath.Join(m.logDir, file.Name()))
			if err != nil {
				continue
			}

			logs = append(logs, models.DeployLog{
				ID:        dateStr,
				Timestamp: timestamp,
				Message:   string(content), // 将内容放在 Message 字段
				Status:    "success",
			})
		}
	}

	// 按时间戳降序排序
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Timestamp.After(logs[j].Timestamp)
	})

	return logs, nil
}

// DeleteDeployLog 删除部署日志
func (m *DeployManager) DeleteDeployLog(logID string) error {
	logPath := filepath.Join(m.logDir, logID+".log")
	if err := os.Remove(logPath); err != nil {
		return fmt.Errorf("删除日志文件失败: %v", err)
	}
	return nil
}

// DeployProject 执行实际的部署操作
func (m *DeployManager) DeployProject(repo string) error {
	// 创建部署日志文件
	header := fmt.Sprintf("部署开始时间: %s\n仓库: %s\n",
		time.Now().Format(time.RFC3339), repo)
	logID, logFile, err := m.logUtil.CreateLogFile(repo, header)
	if err != nil {
		return fmt.Errorf("创建部署日志失败: %v", err)
	}
	defer logFile.Close()

	// 创建临时工作目录
	workDir := filepath.Join(os.TempDir(), fmt.Sprintf("deploy_%s_%s", repo, logID))
	if err := os.MkdirAll(workDir, 0755); err != nil {
		m.logUtil.LogToFile(utils.LogLevelError, logFile, "创建工作目录失败: %v", err)
		return fmt.Errorf("创建工作目录失败: %v", err)
	}
	defer os.RemoveAll(workDir)

	// 1. 拉取代码
	m.logUtil.LogToFile(utils.LogLevelInfo, logFile, "正在从仓库拉取代码: %s", repo)
	if err := m.gitClone(repo, workDir); err != nil {
		m.logUtil.LogToFile(utils.LogLevelError, logFile, "拉取代码失败: %v", err)
		return fmt.Errorf("拉取代码失败: %v", err)
	}

	// 检测项目类型
	projectType := utils.DefaultProjectUtil.DetectProjectType(workDir)
	m.logUtil.LogToFile(utils.LogLevelInfo, logFile, "检测到项目类型: %s", projectType)

	// 2. 构建项目
	m.logUtil.LogToFile(utils.LogLevelInfo, logFile, "正在构建 %s 项目: %s", projectType, repo)
	if err := m.buildProject(workDir); err != nil {
		m.logUtil.LogToFile(utils.LogLevelError, logFile, "构建项目失败: %v", err)
		return fmt.Errorf("构建项目失败: %v", err)
	}

	// 3. 部署服务
	m.logUtil.LogToFile(utils.LogLevelInfo, logFile, "正在部署 %s 服务: %s", projectType, repo)
	if err := m.deployService(workDir); err != nil {
		m.logUtil.LogToFile(utils.LogLevelError, logFile, "部署服务失败: %v", err)
		return fmt.Errorf("部署服务失败: %v", err)
	}

	m.logUtil.LogToFile(utils.LogLevelInfo, logFile, "仓库 %s (%s) 部署成功", repo, projectType)
	return nil
}

// gitClone 从仓库拉取代码
func (m *DeployManager) gitClone(repo, workDir string) error {
	m.logUtil.Log(utils.LogLevelInfo, "开始克隆仓库: %s 到 %s", repo, workDir)

	// 这里可以根据需要添加认证信息
	auth := &githttp.BasicAuth{
		Username: "username", // 可以从配置中读取
		Password: "password",
	}

	err := m.gitUtil.CloneRepo(repo, "main", workDir, auth)
	if err != nil {
		m.logUtil.Log(utils.LogLevelError, "克隆仓库失败: %v", err)
		return err
	}

	m.logUtil.Log(utils.LogLevelInfo, "仓库克隆成功: %s", repo)
	return nil
}

// buildProject 构建项目
func (m *DeployManager) buildProject(workDir string) error {
	// TODO: 实现具体的项目构建逻辑
	// 可能需要根据项目类型选择不同的构建方式
	return nil
}

// deployService 部署服务
func (m *DeployManager) deployService(workDir string) error {
	// TODO: 实现具体的服务部署逻辑
	// 可能需要调用 ServiceManager 的方法
	return nil
}
