package managers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"servon/core/internal/events"
	"servon/core/internal/managers/github"
	"servon/core/internal/utils"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// DeployManager 是部署管理器，负责处理项目的自动化部署流程
// 主要功能包括：
// - 监听代码仓库的推送事件
// - 执行项目的自动化部署
// - 记录和管理部署日志
// - 发布部署相关的事件通知
type DeployManager struct {
	// eventBus 用于处理事件的发布与订阅
	eventBus *events.EventBus
	// logger 用于记录部署过程的日志
	logger *utils.LogUtil
	// gitUtil 用于处理Git操作
	gitUtil *utils.GitUtil
	github  *github.GitHubIntegration
}

func NewDeployManager(eventBus *events.EventBus, github *github.GitHubIntegration) (*DeployManager, error) {
	dm := &DeployManager{
		eventBus: eventBus,
		logger:   utils.NewLogUtil(deployLogDir),
		gitUtil:  utils.NewGitUtil(),
		github:   github,
	}

	// 订阅Git Push事件
	eventBus.Subscribe(events.GitPush, dm.handleGitPushEvent)

	return dm, nil
}

// handleGitPushEvent 处理Git Push事件
func (m *DeployManager) handleGitPushEvent(event events.Event) {
	deployData, ok := event.Data.(map[string]interface{})
	if !ok {
		m.logger.Error("无效的部署数据格式")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		m.logger.Error("缺少仓库信息")
		return
	}

	// 执行部署操作
	if err := m.DeployProject(repo); err != nil {
		m.logger.Errorf("错误: 仓库 %s 部署失败: %v", repo, err)

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

	m.logger.Infof("仓库 %s 部署成功完成", repo)
	// 发布部署成功事件
	m.eventBus.Publish(events.Event{
		Type: events.DeployComplete,
		Data: map[string]interface{}{
			"repository": repo,
			"status":     "success",
		},
	})
}

// DeployProject 执行实际的部署操作
func (m *DeployManager) DeployProject(repoURL string) error {
	// 生成唯一的部署ID，根据当前日期和时间
	deployID := time.Now().Format("20060102150405")

	// 创建临时工作目录
	workDir := filepath.Join(os.TempDir(), fmt.Sprintf("deploy_%s_%s", repoURL, deployID))
	m.logger.Infof("创建临时工作目录: %s", workDir)

	if err := os.MkdirAll(workDir, 0755); err != nil {
		m.logger.LogAndReturnErrorf("创建工作目录失败: %v", err)
	}
	defer func() {
		m.logger.Infof("清理临时工作目录: %s", workDir)
		os.RemoveAll(workDir)
	}()

	// 1. 拉取代码
	m.logger.Infof("开始从仓库拉取代码: %s", repoURL)
	if err := m.gitClone(repoURL, workDir); err != nil {
		return m.logger.LogAndReturnErrorf("拉取代码失败: %v", err)
	}

	// 检测项目类型
	m.logger.Infof("开始检测项目类型...")
	projectType := utils.DefaultProjectUtil.DetectProjectType(workDir)
	m.logger.Infof("检测到项目类型: %s", projectType)

	// 2. 构建项目
	m.logger.Infof("开始构建 %s 项目: %s", projectType, repoURL)
	if err := m.buildProject(workDir); err != nil {
		return m.logger.LogAndReturnErrorf("构建项目失败: %v", err)
	}

	// 3. 部署服务
	m.logger.Infof("开始部署 %s 服务: %s", projectType, repoURL)
	if err := m.deployService(workDir); err != nil {
		return m.logger.LogAndReturnErrorf("部署服务失败: %v", err)
	}

	m.logger.Infof("仓库 %s (%s) 部署成功完成", repoURL, projectType)
	return nil
}

// gitClone 从仓库拉取代码（带重试机制）
func (m *DeployManager) gitClone(repo, workDir string) error {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			m.logger.Infof("第 %d 次重试克隆仓库...", i+1)
			time.Sleep(time.Second * time.Duration(i+1)) // 递增重试延迟
		}

		// 获取认证信息
		m.logger.Infof("正在获取 GitHub 认证信息...")
		auth, err := m.getGitHubAuth(repo)
		if err != nil {
			lastErr = fmt.Errorf("获取GitHub认证信息失败: %v", err)
			m.logger.Errorf("获取GitHub认证信息失败: %v", err)
			continue
		}
		m.logger.Infof("成功获取 GitHub 认证信息")

		// 尝试克隆
		m.logger.Infof("开始克隆仓库 %s 到 %s", repo, workDir)
		err = m.gitUtil.CloneRepo(repo, "main", workDir, auth)
		if err == nil {
			m.logger.Infof("仓库克隆成功: %s", repo)
			return nil
		}

		lastErr = err
		m.logger.Errorf("克隆失败: %v", err)
	}

	return m.logger.LogAndReturnErrorf("克隆仓库失败（已重试%d次）: %v", maxRetries, lastErr)
}

// getGitHubAuth 获取GitHub认证信息
func (m *DeployManager) getGitHubAuth(repo string) (*githttp.BasicAuth, error) {
	if m.github == nil {
		return nil, m.logger.LogAndReturnErrorf("GitHub集成未初始化")
	}

	m.logger.Infof("正在获取仓库 %s 的安装令牌", repo)
	token, err := m.github.GetInstallationToken(repo)
	if err != nil {
		return nil, m.logger.LogAndReturnErrorf("获取GitHub认证令牌失败: %v", err)
	}
	m.logger.Infof("成功获取安装令牌")

	return &githttp.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}, nil
}

// buildProject 构建项目
func (m *DeployManager) buildProject(workDir string) error {
	// TODO: 实现具体的项目构建逻辑
	// 可能需要根据项目类型选择不同的构建方式
	return nil
}

// deployService 部署服务
func (m *DeployManager) deployService(workDir string) error {
	// TODO: 实现具体的服务部署逻辑
	// 可能需要调用 ServiceManager 的方法
	return nil
}
