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

// GitUtil 提供Git操作相关的功能，如果提供了logger，则使用logger记录日志，否则不记录日志
type GitUtil struct {
	logger *LogUtil
}

// NewGitUtil 创建新的Git工具实例
func NewGitUtil(logger *LogUtil) *GitUtil {
	return &GitUtil{
		logger: logger,
	}
}

// CloneRepo 克隆代码仓库
func (g *GitUtil) CloneRepo(url, branch, targetDir string, auth *http.BasicAuth) error {
	// 记录克隆开始
	if g.logger != nil {
		g.logger.Infof("开始克隆仓库\n")
		g.logger.Infof("URL: %s\n", url)
		g.logger.Infof("分支: %s\n", branch)
		g.logger.Infof("目标目录: %s\n", targetDir)
		g.logger.Infof("使用认证: %v\n", auth != nil)
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 克隆选项
	cloneOpts := &git.CloneOptions{
		URL:          url,
		Progress:     os.Stdout, // 显示进度
		SingleBranch: true,
		Auth:         auth,
	}

	// 如果指定了分支，设置分支
	if branch != "" {
		if g.logger != nil {
			g.logger.Infof("设置克隆分支: %s\n", branch)
		}
		cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}

	// 尝试克隆
	if g.logger != nil {
		g.logger.Infof("执行克隆操作...\n")
	}
	repo, err := git.PlainClone(targetDir, false, cloneOpts)
	if err != nil {
		if auth != nil {
			// 脱敏处理，只显示token长度
			if g.logger != nil {
				g.logger.Infof("克隆失败 (使用认证 - token长度: %d): %v\n", len(auth.Password), err)
			}
		} else {
			if g.logger != nil {
				g.logger.Infof("克隆失败 (无认证): %v\n", err)
			}
		}
		return fmt.Errorf("克隆仓库失败: %v", err)
	}

	// 验证克隆结果
	if repo != nil {
		head, err := repo.Head()
		if err == nil {
			if g.logger != nil {
				g.logger.Infof("克隆成功 - 当前HEAD: %s\n", head.Hash())
			}
		}
	}

	if g.logger != nil {
		g.logger.Infof("仓库克隆完成\n")
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
