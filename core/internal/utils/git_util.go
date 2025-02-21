package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitUtil 提供Git操作相关的功能
type GitUtil struct {
}

// NewGitUtil 创建新的Git工具实例
func NewGitUtil() *GitUtil {
	return &GitUtil{}
}

// CloneRepo 克隆代码仓库
func (g *GitUtil) CloneRepo(url, branch, targetDir string, auth *http.BasicAuth) error {
	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 克隆选项
	cloneOpts := &git.CloneOptions{
		URL:          url,
		Progress:     os.Stdout,
		SingleBranch: true,
		Auth:         auth,
	}

	// 如果指定了分支，设置分支
	if branch != "" {
		cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}

	// 尝试克隆
	_, err := git.PlainClone(targetDir, false, cloneOpts)
	if err != nil {
		return fmt.Errorf("克隆仓库失败: %v", err)
	}

	return nil
}

// PullRepo 拉取最新代码
func (g *GitUtil) PullRepo(repoPath string, auth *http.BasicAuth) error {
	// 打开仓库
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("打开仓库失败: %v", err)
	}

	// 获取工作目录
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 拉取更新
	err = w.Pull(&git.PullOptions{
		Auth:         auth,
		Progress:     os.Stdout,
		SingleBranch: true,
	})

	if err == git.NoErrAlreadyUpToDate {
		return nil
	}

	if err != nil {
		return fmt.Errorf("拉取更新失败: %v", err)
	}

	return nil
}

// GetLastCommit 获取最后一次提交信息
func (g *GitUtil) GetLastCommit(repoPath string) (*object.Commit, error) {
	// 打开仓库
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("打开仓库失败: %v", err)
	}

	// 获取HEAD引用
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("获取HEAD引用失败: %v", err)
	}

	// 获取提交对象
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("获取提交对象失败: %v", err)
	}

	return commit, nil
}
