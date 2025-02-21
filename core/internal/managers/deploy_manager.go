package managers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"servon/core/internal/events"
	"servon/core/internal/models"
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
	// logger 用于记录部署过程的日志
	logger *log.Logger
	// logDir 指定部署日志文件的存储目录
	logDir string
}

func NewDeployManager(eventBus *events.EventBus) (*DeployManager, error) {
	// 确保日志目录存在
	logDir := "/data/deploy"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建部署日志目录失败: %v", err)
	}

	// 创建或打开日志文件（按日期分文件）
	logFile, err := os.OpenFile(
		filepath.Join(logDir, time.Now().Format("2006-01-02")+".log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("打开部署日志文件失败: %v", err)
	}

	dm := &DeployManager{
		eventBus: eventBus,
		logger:   log.New(logFile, "", log.LstdFlags|log.Lmicroseconds),
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
		m.logger.Printf("错误: 无效的部署数据格式")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		m.logger.Printf("错误: 缺少仓库信息")
		return
	}

	// 执行部署操作
	if err := m.DeployProject(repo); err != nil {
		m.logger.Printf("错误: 仓库 %s 部署失败: %v", repo, err)

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

	m.logger.Printf("信息: 仓库 %s 部署成功完成", repo)
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
		m.logger.Printf("错误: 读取日志文件失败: %v", err)
		return nil, fmt.Errorf("读取日志文件失败: %v", err)
	}

	timestamp, err := time.Parse("2006-01-02", logID)
	if err != nil {
		m.logger.Printf("错误: 无效的日志ID格式: %v", err)
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

// createNewLogFile 创建新的日志文件
func (m *DeployManager) createNewLogFile(repo string) (string, error) {
	timestamp := time.Now()
	logID := timestamp.Format("2006-01-02-150405")
	logPath := filepath.Join(m.logDir, logID+".log")

	file, err := os.Create(logPath)
	if err != nil {
		return "", fmt.Errorf("创建日志文件失败: %v", err)
	}
	defer file.Close()

	// 写入日志头部信息
	_, err = fmt.Fprintf(file, "部署开始时间: %s\n仓库: %s\n\n",
		timestamp.Format(time.RFC3339), repo)
	if err != nil {
		return "", fmt.Errorf("写入日志头部信息失败: %v", err)
	}

	return logID, nil
}

// DeployProject 执行实际的部署操作
func (m *DeployManager) DeployProject(repo string) error {
	// 创建新的日志文件
	logID, err := m.createNewLogFile(repo)
	if err != nil {
		return fmt.Errorf("创建部署日志失败: %v", err)
	}

	m.logger.Printf("信息: 开始部署仓库: %s (日志ID: %s)", repo, logID)

	// 发布部署开始事件
	m.eventBus.Publish(events.Event{
		Type: events.DeployStart,
		Data: map[string]interface{}{
			"repository": repo,
		},
	})

	m.logger.Printf("信息: 正在从仓库拉取代码: %s", repo)
	// 1. 拉取代码

	m.logger.Printf("信息: 正在构建项目: %s", repo)
	// 2. 构建项目

	m.logger.Printf("信息: 正在部署服务: %s", repo)
	// 3. 部署服务

	return nil
}
