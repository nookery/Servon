package managers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servon/components/events"
	"servon/core/internal/managers/deployers"
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
	eventBus events.IEventBus
	// logger 用于记录部署过程的日志
	logger *utils.LogUtil
	// gitUtil 用于处理Git操作
	gitUtil     *utils.GitUtil
	fileUtil    *utils.FileUtil
	stringUtil  *utils.StringUtil
	github      *github.GitHubIntegration
	logsDir     string
	tempDir     string
	projectsDir string
	deployers   []deployers.Deployer
}

func NewDeployManager(eventBus events.IEventBus, github *github.GitHubIntegration, logsDir string, tempDir string, projectsDir string) (*DeployManager, error) {
	dm := &DeployManager{
		eventBus:    eventBus,
		logger:      utils.NewTopicLogUtil(logsDir, "deploy"),
		gitUtil:     utils.NewGitUtil(utils.NewLogUtil(logsDir)),
		fileUtil:    utils.DefaultFileUtil,
		github:      github,
		logsDir:     logsDir,
		tempDir:     tempDir,
		projectsDir: projectsDir,
		deployers:   []deployers.Deployer{},
	}

	// 订阅Git Push事件
	eventBus.Subscribe(events.GitPush, dm.handleGitPushEvent)

	return dm, nil
}

// handleGitPushEvent 处理Git Push事件
func (m *DeployManager) handleGitPushEvent(event events.Event) {
	deployData, ok := event.Data.(map[string]interface{})
	if !ok {
		m.logger.ErrorMessage("无效的部署数据格式")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		m.logger.ErrorMessage("缺少仓库信息")
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

	// 获取项目名称
	projectName := m.stringUtil.GetProjectNameFromString(repoURL)

	// 部署的目标目录
	targetDir := filepath.Join(m.projectsDir, projectName)

	// 创建临时工作目录
	workDir := filepath.Join(m.tempDir, "deploy", fmt.Sprintf("%s_%s", projectName, deployID))
	m.logger.Infof("创建临时工作目录: %s", workDir)

	if err := os.MkdirAll(workDir, 0755); err != nil {
		m.logger.LogAndReturnErrorf("创建工作目录失败: %v", err)
	}
	defer func() {
		m.logger.Infof("清理临时工作目录: %s", workDir)
		os.RemoveAll(workDir)
	}()

	// 拉取代码
	m.logger.Infof("开始从仓库拉取代码: %s", repoURL)
	if err := m.gitClone(repoURL, workDir); err != nil {
		return m.logger.LogAndReturnErrorf("拉取代码失败: %v", err)
	}

	// 检测项目类型
	projectType := utils.DefaultProjectUtil.DetectProjectType(workDir)
	m.logger.Infof("检测到项目类型: %s", projectType)

	if projectType == "unknown" {
		return m.logger.LogAndReturnErrorf("未检测到项目类型，部署失败")
	}

	// 根据项目类型选择合适的部署器
	deployer := m.getDeployer(projectType)
	if deployer == nil {
		return m.logger.LogAndReturnErrorf("未找到合适的部署器")
	}

	m.logger.Infof("使用部署器: %s", deployer.GetName())

	// 执行部署
	if err := deployer.Deploy(projectName, workDir, targetDir, m.logger); err != nil {
		return m.logger.LogAndReturnErrorf("部署失败: %v", err)
	}

	m.logger.Infof("部署成功")

	return nil
}

// gitClone 从仓库拉取代码（带重试机制）
func (m *DeployManager) gitClone(repo, workDir string) error {
	const maxRetries = 3
	var lastErr error

	// 规范化仓库地址
	originalRepo := repo
	if !strings.HasPrefix(repo, "https://") && !strings.HasPrefix(repo, "git@") {
		repo = "https://github.com/" + repo
	}
	m.logger.Infof("规范化仓库地址: %s -> %s", originalRepo, repo)

	// 检查工作目录
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		m.logger.Infof("工作目录不存在，创建: %s", workDir)
		if err := os.MkdirAll(workDir, 0755); err != nil {
			return m.logger.LogAndReturnErrorf("创建工作目录失败: %v", err)
		}
	}

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			m.logger.Infof("第 %d 次重试克隆仓库...", i+1)
			time.Sleep(time.Second * time.Duration(i+1))
		}

		m.logger.Infof("开始获取 GitHub 认证信息...")
		auth, err := m.getGitHubAuth(repo)
		if err != nil {
			lastErr = fmt.Errorf("获取GitHub认证信息失败: %v", err)
			m.logger.Errorf("认证失败详情: %v", lastErr)
			continue
		}

		if auth == nil {
			m.logger.Warnf("获取到的认证信息为空，将尝试无认证克隆")
		} else {
			m.logger.Infof("成功获取认证信息 - 用户名: %s, Token长度: %d",
				auth.Username, len(auth.Password))
		}

		m.logger.Infof("开始克隆仓库 %s 到 %s", repo, workDir)
		err = m.gitUtil.CloneRepo(repo, "main", workDir, auth)
		if err == nil {
			m.logger.Infof("仓库克隆成功: %s", repo)
			// 验证克隆结果
			if files, err := os.ReadDir(workDir); err == nil {
				m.logger.Infof("克隆目录内容: %d 个文件/目录", len(files))
			}
			return nil
		}

		lastErr = err
		m.logger.Errorf("克隆失败 (尝试 %d/%d): %v", i+1, maxRetries, err)
	}

	return m.logger.LogAndReturnErrorf("克隆仓库失败（已重试%d次）- 最后错误: %v", maxRetries, lastErr)
}

// getGitHubAuth 获取GitHub认证信息
func (m *DeployManager) getGitHubAuth(repo string) (*githttp.BasicAuth, error) {
	if m.github == nil {
		return nil, m.logger.LogAndReturnErrorf("GitHub集成未初始化")
	}

	m.logger.Infof("准备获取仓库认证令牌: %s", repo)

	// 检查仓库格式
	repoName := repo
	if strings.HasPrefix(repo, "https://github.com/") {
		repoName = strings.TrimPrefix(repo, "https://github.com/")
	}
	m.logger.Infof("处理后的仓库名称: %s", repoName)

	// 验证仓库名称格式
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil, m.logger.LogAndReturnErrorf("无效的仓库名称格式: %s，应为 'owner/repo' 格式", repoName)
	}
	m.logger.Infof("仓库所有者: %s, 仓库名称: %s", parts[0], parts[1])

	token, err := m.github.GetInstallationToken(repoName)
	if err != nil {
		return nil, m.logger.LogAndReturnErrorf("获取安装令牌失败: %v", err)
	}

	if token == "" {
		return nil, m.logger.LogAndReturnErrorf("获取到的token为空")
	}
	m.logger.Infof("成功获取安装令牌 (长度: %d)", len(token))

	auth := &githttp.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}

	// 验证认证信息完整性
	if auth.Username == "" || auth.Password == "" {
		return nil, m.logger.LogAndReturnErrorf("认证信息不完整: username=%v, token_length=%d",
			auth.Username != "", len(auth.Password))
	}

	m.logger.Infof("认证信息构建成功")
	return auth, nil
}

// AddDeployer 添加新的部署器
func (m *DeployManager) AddDeployer(deployer deployers.Deployer) {
	m.logger.Infof("添加新的部署器: %T", deployer)
	m.deployers = append(m.deployers, deployer)
}

// RemoveDeployer 移除指定类型的部署器
func (m *DeployManager) RemoveDeployer(deployerType string) {
	m.logger.Infof("移除部署器: %s", deployerType)
	newDeployers := make([]deployers.Deployer, 0)
	for _, d := range m.deployers {
		if fmt.Sprintf("%T", d) != deployerType {
			newDeployers = append(newDeployers, d)
		}
	}
	m.deployers = newDeployers
}

// GetDeployers 获取所有部署器
func (m *DeployManager) GetDeployers() []deployers.Deployer {
	return m.deployers
}

// ClearDeployers 清空所有部署器
func (m *DeployManager) ClearDeployers() {
	m.logger.Info("清空所有部署器")
	m.deployers = make([]deployers.Deployer, 0)
}

// getDeployer 根据项目类型选择合适的部署器
func (m *DeployManager) getDeployer(projectType string) deployers.Deployer {
	for _, deployer := range m.deployers {
		if deployer.GetName() == projectType {
			return deployer
		}
	}
	return nil
}
