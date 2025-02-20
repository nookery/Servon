package github

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"servon/core/internal/models"
	"time"
)

const (
	githubLogDir    = "/data/github/integration"   // GitHub集成日志目录
	installationDir = "/data/github/installations" // GitHub安装数据目录
	configDir       = "/data/github/config"        // GitHub安装配置目录
	timeFormat      = "2006-01-02"                 // 日期格式
)

var DefaultGitHubLogger = NewGitHubLogger()

// GitHubLogger 处理GitHub集成的日志记录
type GitHubLogger struct {
	logDir string
}

// NewGitHubLogger 创建一个新的GitHub日志记录器
func NewGitHubLogger() *GitHubLogger {
	return &GitHubLogger{
		logDir: githubLogDir,
	}
}

// ensureLogDir 确保日志目录存在
func (l *GitHubLogger) ensureLogDir() error {
	return os.MkdirAll(l.logDir, 0755)
}

// WriteLog 写入日志到文件
func (l *GitHubLogger) WriteLog(logType LogType, message string) error {
	if err := l.ensureLogDir(); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	today := time.Now().Format(timeFormat)
	logFile := filepath.Join(l.logDir, fmt.Sprintf("github-%s.log", today))

	// 以追加模式打开文件
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer f.Close()

	// 写入时间戳和消息
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s][%s] %s %s\n",
		timestamp,
		logType.Name,
		logType.Symbol,
		message)

	if _, err := f.WriteString(logEntry); err != nil {
		return fmt.Errorf("写入日志失败: %v", err)
	}

	return nil
}

// GetLogDir 返回日志目录路径
func (l *GitHubLogger) GetLogDir() string {
	return l.logDir
}

// GetLogFiles 获取日志目录中的所有文件信息
func (l *GitHubLogger) GetLogFiles() ([]models.FileInfo, error) {
	if err := l.ensureLogDir(); err != nil {
		return nil, fmt.Errorf("确保日志目录存在失败: %v", err)
	}

	files, err := os.ReadDir(l.logDir)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	var fileInfos []models.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		fileInfo := models.FileInfo{
			Name:    file.Name(),
			Path:    filepath.Join(l.logDir, file.Name()),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime(),
			IsDir:   file.IsDir(),
			Owner:   "system",
			Group:   "system",
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

// LogInfo 记录信息级别的日志
func (l *GitHubLogger) LogInfo(message string) {
	l.WriteLog(LogType{Name: "INFO", Symbol: "ℹ"}, message)
}

func (l *GitHubLogger) LogInfof(format string, args ...interface{}) {
	l.WriteLog(LogType{Name: "INFO", Symbol: "ℹ"}, fmt.Sprintf(format, args...))
}

// LogError 记录错误级别的日志
func (l *GitHubLogger) LogError(message string) {
	l.WriteLog(LogType{Name: "ERROR", Symbol: "❌"}, message)
}

func (l *GitHubLogger) LogErrorf(format string, args ...interface{}) {
	l.WriteLog(LogType{Name: "ERROR", Symbol: "❌"}, fmt.Sprintf(format, args...))
}

// SaveInstallationData 保存安装数据到指定目录
func (l *GitHubLogger) SaveInstallationData(installationID int64, data []byte) error {
	// 确保目录存在
	if err := os.MkdirAll(installationDir, 0755); err != nil {
		return fmt.Errorf("创建安装数据目录失败: %v", err)
	}

	filename := filepath.Join(installationDir, fmt.Sprintf("%d.json", installationID))
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入安装数据失败: %v", err)
	}

	return nil
}

// SaveRawInstallationData 保存原始安装数据，使用时间戳作为文件名
func (l *GitHubLogger) SaveRawInstallationData(payload []byte) error {
	// 确保目录存在
	if err := os.MkdirAll(installationDir, 0755); err != nil {
		return fmt.Errorf("创建安装数据目录失败: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(installationDir, fmt.Sprintf("raw_%s.json", timestamp))

	if err := os.WriteFile(filename, payload, 0644); err != nil {
		return fmt.Errorf("写入原始安装数据失败: %v", err)
	}

	return nil
}

// SaveInstallationConfig 保存安装配置到指定目录
func (l *GitHubLogger) SaveInstallationConfig(installation *Installation) error {
	config := InstallationConfig{
		InstallationID: installation.ID,
		AccountID:      installation.AccountID,
		AccountLogin:   installation.AccountLogin,
		AccountType:    installation.AccountType,
		AppID:          installation.AppID,
		Permissions:    installation.Permissions,
		Events:         installation.Events,
		Repositories:   installation.Repositories,
		CreatedAt:      installation.CreatedAt,
	}

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 生成配置文件路径
	configPath := filepath.Join(configDir, fmt.Sprintf("installation_%d.json", installation.ID))

	// 序列化并保存配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
