package managers

import (
	"fmt"
	"servon/core/internal/utils"
	"time"
)

const (
	maxGitRetries = 3 // Git操作最大重试次数
)

type GitManager struct {
	gitUtil     *utils.GitUtil
	SoftManager *SoftManager
}

func NewGitManager(softManager *SoftManager) *GitManager {
	return &GitManager{
		gitUtil:     utils.NewGitUtil(nil),
		SoftManager: softManager,
	}
}

// GitClone 克隆一个git仓库到指定目录
// url: 仓库地址
// branch: 分支
// targetDir: 目标目录
func (g *GitManager) GitClone(url string, branch string, targetDir string) error {
	// 如果目标目录存在且不为空，则返回错误
	if utils.IsDirExists(targetDir) && utils.GetDirSize(targetDir) > 0 {
		return fmt.Errorf("目标目录 %s 已存在且不为空", targetDir)
	}

	var lastErr error

	// 首先尝试不使用代理克隆
	for i := 0; i < maxGitRetries; i++ {
		if i > 0 {
			PrintCommandOutput(fmt.Sprintf("克隆失败，第 %d 次重试...", i))
			time.Sleep(time.Second * 2) // 重试前等待一段时间
		}

		err := g.gitUtil.CloneRepo(url, branch, targetDir, nil)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	// 如果常规克隆失败，尝试使用代理
	if !g.SoftManager.IsProxyOn() {
		PrintCommandOutput("常规克隆失败，尝试开启代理重新克隆...")
		software, err := g.SoftManager.OpenProxy()
		if err != nil {
			return fmt.Errorf("开启代理失败: %v，上一次克隆错误: %v", err, lastErr)
		}

		PrintAlert("使用代理软件: " + software + " 克隆仓库...")

		// 使用代理重试克隆
		for i := 0; i < maxGitRetries; i++ {
			if i > 0 {
				PrintCommandOutput(fmt.Sprintf("代理克隆失败，第 %d 次重试...", i))
				time.Sleep(time.Second * 2)
			}

			err := g.gitUtil.CloneRepo(url, branch, targetDir, nil)
			if err == nil {
				return nil
			}
			lastErr = err
		}
	}

	return fmt.Errorf("克隆失败（已尝试使用代理）: %v", lastErr)
}

// progressWriter 实现了io.Writer接口，用于处理git操作的进度输出
type progressWriter struct{}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	PrintCommandOutput(string(p))
	return len(p), nil
}
