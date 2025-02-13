package libs

import (
	"github.com/go-git/go-git/v5"
)

type GitManager struct {
}

func NewGitManager() *GitManager {
	return &GitManager{}
}

// GitClone 克隆一个git仓库到指定目录
// url: 仓库地址
// targetDir: 目标目录
func (g *GitManager) GitClone(url string, targetDir string) error {
	options := &git.CloneOptions{
		URL:      url,
		Progress: &progressWriter{}, // 使用自定义的进度写入器
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
