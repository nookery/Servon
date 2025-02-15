package libs

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

const (
	maxGitRetries = 3 // Git操作最大重试次数
)

type GitManager struct {
}

func NewGitManager() *GitManager {
	return &GitManager{}
}

// GitClone 克隆一个git仓库到指定目录
// url: 仓库地址
// branch: 分支
// targetDir: 目标目录
func (g *GitManager) GitClone(url string, branch string, targetDir string) error {
	// 如果目标目录存在且不为空，则返回错误
	if DefaultFilesManager.IsDirExists(targetDir) && DefaultFilesManager.GetDirSize(targetDir) > 0 {
		return fmt.Errorf("目标目录 %s 已存在且不为空", targetDir)
	}

	var lastErr error

	// 首先尝试不使用代理克隆
	for i := 0; i < maxGitRetries; i++ {
		if i > 0 {
			PrintCommandOutput(fmt.Sprintf("克隆失败，第 %d 次重试...", i))
			time.Sleep(time.Second * 2) // 重试前等待一段时间
		}

		err := g.cloneRepo(url, branch, targetDir)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	// 如果常规克隆失败，尝试使用代理
	if !DefaultProxyManager.IsProxyOn() {
		PrintCommandOutput("常规克隆失败，尝试开启代理重新克隆...")
		software, err := DefaultProxyManager.OpenProxy()
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

			err := g.cloneRepo(url, branch, targetDir)
			if err == nil {
				return nil
			}
			lastErr = err
		}
	}

	return fmt.Errorf("克隆失败（已尝试使用代理）: %v", lastErr)
}

// cloneRepo 执行实际的克隆操作
func (g *GitManager) cloneRepo(url string, branch string, targetDir string) error {
	PrintInfo("克隆仓库: " + url + " 的分支 " + branch + " 到 " + targetDir)

	options := &git.CloneOptions{
		URL:           url,
		Progress:      &progressWriter{}, // 使用自定义的进度写入器
		ReferenceName: plumbing.ReferenceName(branch),
	}

	_, err := git.PlainClone(targetDir, false, options)
	return err
}

// progressWriter 实现了io.Writer接口，用于处理git操作的进度输出
type progressWriter struct{}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	DefaultPrinter.PrintCommandOutput(string(p))
	return len(p), nil
}

// --- 以下是 cobra 命令 ---

// GetGitRootCommand 获取git命令根命令，返回一个 cobra.Command
func (g *GitManager) GetGitRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "git",
		Short: "git 命令，比默认的 git 命令更智能",
	}

	rootCmd.AddCommand(g.GetCloneCommand())

	return rootCmd
}

// GetCloneCommand 获取克隆命令，返回一个 cobra.Command
func (g *GitManager) GetCloneCommand() *cobra.Command {
	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "克隆一个git仓库",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			branch, _ := cmd.Flags().GetString("branch")
			targetDir, _ := cmd.Flags().GetString("target-dir")

			err := g.GitClone(url, branch, targetDir)
			if err != nil {
				PrintError(err)
			}
		},
	}

	cloneCmd.Flags().StringP("url", "u", "", "仓库地址")
	cloneCmd.Flags().StringP("branch", "b", "master", "分支")
	cloneCmd.Flags().StringP("target-dir", "t", "", "目标目录")

	return cloneCmd
}
