// Package git_util 提供Git操作功能
package git_util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitUtil 提供Git操作相关的功能
type GitUtil struct{}

// NewGitUtil 创建新的Git工具实例
func NewGitUtil() *GitUtil {
	return &GitUtil{}
}

// CloneRepo 克隆代码仓库
func (g *GitUtil) CloneRepo(url, branch, targetDir string, auth *http.BasicAuth) error {
	if auth != nil {
		// 确保URL使用HTTPS格式
		if !strings.HasPrefix(url, "https://") {
			url = "https://" + strings.TrimPrefix(url, "http://")
		}

		// 在URL中嵌入认证信息
		urlParts := strings.Split(url, "//")
		url = fmt.Sprintf("%s//%s:%s@%s",
			urlParts[0],
			auth.Username,
			auth.Password,
			urlParts[1])
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 克隆选项
	cloneOptions := &git.CloneOptions{
		URL: url,
	}

	// 如果指定了分支，则只克隆该分支
	if branch != "" {
		cloneOptions.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
		cloneOptions.SingleBranch = true
	}

	// 执行克隆
	_, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("克隆仓库失败: %v", err)
	}

	return nil
}

// GetCurrentBranch 获取当前分支名
func (g *GitUtil) GetCurrentBranch(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("打开仓库失败: %v", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("获取HEAD失败: %v", err)
	}

	branchName := head.Name().Short()
	return branchName, nil
}

// GetCommitInfo 获取最新提交信息
func (g *GitUtil) GetCommitInfo(repoPath string) (*object.Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("打开仓库失败: %v", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("获取HEAD失败: %v", err)
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("获取提交对象失败: %v", err)
	}

	return commit, nil
}

// IsGitRepo 检查目录是否为Git仓库
func (g *GitUtil) IsGitRepo(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

// GetRemoteURL 获取远程仓库URL
func (g *GitUtil) GetRemoteURL(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("打开仓库失败: %v", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("获取远程仓库失败: %v", err)
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("没有找到远程仓库")
	}

	// 通常使用origin作为默认远程仓库
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			if len(remote.Config().URLs) > 0 {
				return remote.Config().URLs[0], nil
			}
		}
	}

	// 如果没有origin，返回第一个远程仓库的URL
	if len(remotes[0].Config().URLs) > 0 {
		return remotes[0].Config().URLs[0], nil
	}

	return "", fmt.Errorf("没有找到远程仓库URL")
}

// PullLatest 拉取最新代码
func (g *GitUtil) PullLatest(repoPath string, auth *http.BasicAuth) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("打开仓库失败: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("获取工作树失败: %v", err)
	}

	pullOptions := &git.PullOptions{}
	if auth != nil {
		pullOptions.Auth = auth
	}

	err = w.Pull(pullOptions)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("拉取代码失败: %v", err)
	}

	return nil
}
