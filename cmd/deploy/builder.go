package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"servon/cmd/utils"
	"strings"
	"sync"
)

// 添加日志相关的变量和结构
var (
	projectLogs = make(map[int][]string)
	logsMutex   sync.RWMutex
)

// GetProjectLogs 获取项目日志
func GetProjectLogs(id int) ([]string, error) {
	logsMutex.RLock()
	defer logsMutex.RUnlock()

	logs, exists := projectLogs[id]
	if !exists {
		return []string{}, nil
	}
	return logs, nil
}

// addLog 添加日志
func addLog(id int, log string) {
	logsMutex.Lock()
	defer logsMutex.Unlock()

	if _, exists := projectLogs[id]; !exists {
		projectLogs[id] = make([]string, 0)
	}

	// 限制日志数量，只保留最新的1000条
	logs := projectLogs[id]
	if len(logs) >= 1000 {
		logs = logs[1:]
	}
	projectLogs[id] = append(logs, log)
}

// BuildProject 构建项目
func BuildProject(id int, logChan chan<- string) error {
	utils.Info("开始构建项目: %d", id)

	projectsMu.RLock()
	project, exists := projects[id]
	projectsMu.RUnlock()

	if !exists {
		utils.Error("项目不存在: %d", id)
		return fmt.Errorf("项目不存在")
	}

	// 修改日志发送逻辑
	sendLog := func(log string) {
		utils.Debug("[项目 %d] %s", id, log)
		addLog(id, log)
		logChan <- log
	}

	// 更新项目状态
	project.Status = "building"
	saveProjects()
	sendLog("开始部署...")

	utils.Info("项目构建完成: [%d] %s", id, project.Name)
	sendLog("部署完成")
	return nil
}

// gitSync 同步Git仓库
func gitSync(project *Project, dir string, logFn func(string)) error {
	logFn("正在同步代码...")

	if _, err := os.Stat(filepath.Join(dir, ".git")); os.IsNotExist(err) {
		// 克隆仓库
		cmd := exec.Command("git", "clone", "-b", project.Branch, project.GitRepo, ".")
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("克隆代码失败: %v\n%s", err, output)
		}
	} else {
		// 更新仓库
		cmd := exec.Command("git", "pull", "origin", project.Branch)
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("更新代码失败: %v\n%s", err, output)
		}
	}

	logFn("代码同步完成")
	return nil
}

// build 执行构建命令
func build(project *Project, dir string, logFn func(string)) error {
	if project.BuildCmd == "" {
		return nil
	}

	logFn("开始构建...")

	// 设置环境变量
	env := os.Environ()
	for _, e := range project.Environment {
		env = append(env, fmt.Sprintf("%s=%s", e.Key, e.Value))
	}

	// 分割命令和参数
	parts := strings.Fields(project.BuildCmd)
	if len(parts) == 0 {
		return fmt.Errorf("无效的构建命令")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("构建失败: %v\n%s", err, output)
	}

	logFn(fmt.Sprintf("构建输出:\n%s", output))
	logFn("构建完成")
	return nil
}
