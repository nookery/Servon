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
)

type DeployLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Repo      string    `json:"repo"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
}

type DeployManager struct {
	eventBus *events.EventBus
	logger   *log.Logger
	logDir   string
}

func NewDeployManager(eventBus *events.EventBus) (*DeployManager, error) {
	// 确保日志目录存在
	logDir := "/data/deploy"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create deploy log directory: %v", err)
	}

	// 创建或打开日志文件（按日期分文件）
	logFile, err := os.OpenFile(
		filepath.Join(logDir, time.Now().Format("2006-01-02")+".log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open deploy log file: %v", err)
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
		m.logger.Printf("ERROR: Invalid deploy data format")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		m.logger.Printf("ERROR: Repository information missing")
		return
	}

	// 执行部署操作
	if err := m.deployProject(repo); err != nil {
		m.logger.Printf("ERROR: Deploy failed for repository %s: %v", repo, err)

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

	m.logger.Printf("INFO: Deploy completed successfully for repository %s", repo)
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
func (m *DeployManager) GetDeployLog(logID string) (*DeployLog, error) {
	logPath := filepath.Join(m.logDir, logID+".log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		m.logger.Printf("ERROR: failed to read log file: %v", err)
		return nil, fmt.Errorf("failed to read log file: %v", err)
	}

	timestamp, err := time.Parse("2006-01-02", logID)
	if err != nil {
		m.logger.Printf("ERROR: invalid log ID format: %v", err)
		return nil, fmt.Errorf("invalid log ID format: %v", err)
	}

	return &DeployLog{
		ID:        logID,
		Timestamp: timestamp,
		Message:   string(content), // 将内容放在 Message 字段
		Status:    "success",
	}, nil
}

// ListDeployLogs 获取部署日志列表
func (m *DeployManager) ListDeployLogs() ([]DeployLog, error) {
	files, err := os.ReadDir(m.logDir)
	if err != nil {
		return nil, printer.PrintAndReturnErrorf("failed to read log directory: %v", err)
	}

	logs := make([]DeployLog, 0, len(files))
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

			logs = append(logs, DeployLog{
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
		return fmt.Errorf("failed to delete log file: %v", err)
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
		return "", fmt.Errorf("failed to create log file: %v", err)
	}
	defer file.Close()

	// 写入日志头部信息
	_, err = fmt.Fprintf(file, "Deploy started at: %s\nRepository: %s\n\n",
		timestamp.Format(time.RFC3339), repo)
	if err != nil {
		return "", fmt.Errorf("failed to write log header: %v", err)
	}

	return logID, nil
}

// deployProject 执行实际的部署操作
func (m *DeployManager) deployProject(repo string) error {
	// 创建新的日志文件
	logID, err := m.createNewLogFile(repo)
	if err != nil {
		return fmt.Errorf("failed to create deploy log: %v", err)
	}

	m.logger.Printf("INFO: Starting deployment for repository: %s (LogID: %s)", repo, logID)

	// 发布部署开始事件
	m.eventBus.Publish(events.Event{
		Type: events.DeployStart,
		Data: map[string]interface{}{
			"repository": repo,
		},
	})

	m.logger.Printf("INFO: Pulling code from repository: %s", repo)
	// 1. 拉取代码

	m.logger.Printf("INFO: Building project for repository: %s", repo)
	// 2. 构建项目

	m.logger.Printf("INFO: Deploying service for repository: %s", repo)
	// 3. 部署服务

	return nil
}
