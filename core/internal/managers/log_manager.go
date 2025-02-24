package managers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"servon/core/internal/utils"
)

var DefaultLogDir = "/data/logs"

var DefaultLogManager, _ = NewLogManager(DefaultLogDir)

// LogEntry 表示一条日志记录
type LogEntry struct {
	Timestamp time.Time       `json:"time"`
	Level     string          `json:"level"`
	Caller    string          `json:"caller"`
	Message   string          `json:"message"`
	Extra     json.RawMessage `json:"extra,omitempty"`
}

// LogManager 负责管理系统日志
type LogManager struct {
	baseLogDir string
	*utils.LogUtil
}

// NewLogManager 创建日志管理器实例
func NewLogManager(baseLogDir string) (*LogManager, error) {
	if baseLogDir == "" {
		baseLogDir = DefaultLogDir
	}

	// 确保日志目录存在
	if err := os.MkdirAll(baseLogDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	return &LogManager{
		baseLogDir: baseLogDir,
		LogUtil:    utils.NewLogUtil(baseLogDir),
	}, nil
}

// ListLogFiles 列出指定目录下的所有日志文件
func (m *LogManager) ListLogFiles(subDir string) ([]string, error) {
	dir := filepath.Join(m.baseLogDir, subDir)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	var logFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" {
			logFiles = append(logFiles, filepath.Join(subDir, file.Name()))
		}
	}

	// 按修改时间排序，最新的在前
	sort.Slice(logFiles, func(i, j int) bool {
		fi, _ := os.Stat(filepath.Join(m.baseLogDir, logFiles[i]))
		fj, _ := os.Stat(filepath.Join(m.baseLogDir, logFiles[j]))
		return fi.ModTime().After(fj.ModTime())
	})

	return logFiles, nil
}

// ReadLogEntries 读取指定日志文件的内容
func (m *LogManager) ReadLogEntries(logFile string, limit int) ([]LogEntry, error) {
	fullPath := filepath.Join(m.baseLogDir, logFile)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && (limit == 0 || len(entries) < limit) {
		var entry LogEntry
		// 先解析为临时结构，处理时间格式
		var tempEntry struct {
			Time    string          `json:"time"`
			Level   string          `json:"level"`
			Caller  string          `json:"caller"`
			Message string          `json:"message"`
			Extra   json.RawMessage `json:"extra,omitempty"`
		}

		if err := json.Unmarshal(scanner.Bytes(), &tempEntry); err != nil {
			m.WarnfConsole("解析日志行失败: %v", err)
			continue
		}

		// 解析时间字符串
		t, err := time.Parse("2006-01-02 15:04:05.000", tempEntry.Time)
		if err != nil {
			m.WarnfConsole("解析时间失败: %v", err)
			continue
		}

		// 构建最终的日志条目
		entry = LogEntry{
			Timestamp: t,
			Level:     tempEntry.Level,
			Caller:    tempEntry.Caller,
			Message:   tempEntry.Message,
			Extra:     tempEntry.Extra,
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// SearchLogs 在日志中搜索指定关键字
func (m *LogManager) SearchLogs(subDir, keyword string) ([]LogEntry, error) {
	files, err := m.ListLogFiles(subDir)
	if err != nil {
		return nil, err
	}

	var results []LogEntry
	for _, file := range files {
		entries, err := m.ReadLogEntries(file, 0)
		if err != nil {
			m.Warnf("读取日志文件 %s 失败: %v", file, err)
			continue
		}

		for _, entry := range entries {
			if entry.containsKeyword(keyword) {
				results = append(results, entry)
			}
		}
	}

	return results, nil
}

// containsKeyword 检查日志条目是否包含指定关键字
func (e *LogEntry) containsKeyword(keyword string) bool {
	return contains(e.Message, keyword) ||
		contains(e.Level, keyword) ||
		contains(e.Caller, keyword)
}

// contains 检查字符串是否包含关键字（不区分大小写）
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// CleanOldLogs 清理指定天数之前的日志
func (m *LogManager) CleanOldLogs(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)

	err := filepath.Walk(m.baseLogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".log" && info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				m.Warnf("删除旧日志文件失败 %s: %v", path, err)
				return err
			}
			m.Infof("已删除旧日志文件: %s", path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("清理旧日志失败: %v", err)
	}
	return nil
}

// GetLogStats 获取日志统计信息
func (m *LogManager) GetLogStats(subDir string) (map[string]int, error) {
	stats := map[string]int{
		"error": 0,
		"warn":  0,
		"info":  0,
		"debug": 0,
	}

	// 如果日志文件不存在，直接返回空统计
	logPath := filepath.Join(subDir, "app.log")
	if _, err := os.Stat(filepath.Join(m.baseLogDir, logPath)); os.IsNotExist(err) {
		return stats, nil
	}

	entries, err := m.ReadLogEntries(logPath, 0)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if count, exists := stats[strings.ToLower(entry.Level)]; exists {
			stats[strings.ToLower(entry.Level)] = count + 1
		}
	}

	return stats, nil
}

// DeleteLogFile 删除指定的日志文件
func (m *LogManager) DeleteLogFile(logPath string) error {
	// 确保文件路径在日志目录内，防止任意文件删除
	fullPath := filepath.Join(m.baseLogDir, logPath)
	if !strings.HasPrefix(fullPath, m.baseLogDir) {
		return fmt.Errorf("无效的日志文件路径")
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("日志文件不存在: %s", logPath)
	}

	// 删除文件
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("删除日志文件失败: %v", err)
	}

	m.Infof("已删除日志文件: %s", logPath)
	return nil
}

// ClearLogFile 清空指定的日志文件
func (m *LogManager) ClearLogFile(logFile string) error {
	fullPath := filepath.Join(m.baseLogDir, logFile)
	// 清空文件内容
	return os.Truncate(fullPath, 0)
}
