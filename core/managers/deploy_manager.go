package managers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servon/components/events"
	"servon/components/git_util"
	"servon/components/github"
	"servon/components/utils"
	"servon/core/contract"

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
	// gitUtil 用于处理Git操作
	gitUtil     *git_util.GitUtil
	fileUtil    *utils.FileUtil
	stringUtil  *utils.StringUtil
	github      *github.GitHubIntegration
	logsDir     string
	tempDir     string
	projectsDir string
	deployers   []contract.SuperDeployer
}

func NewDeployManager(eventBus events.IEventBus, github *github.GitHubIntegration, logsDir string, tempDir string, projectsDir string) (*DeployManager, error) {
	dm := &DeployManager{
		eventBus:    eventBus,
		gitUtil:     git_util.NewGitUtil(),
		fileUtil:    utils.DefaultFileUtil,
		github:      github,
		logsDir:     logsDir,
		tempDir:     tempDir,
		projectsDir: projectsDir,
		deployers:   []contract.SuperDeployer{},
	}

	// 订阅Git Push事件
	eventBus.Subscribe(events.GitPush, dm.handleGitPushEvent)

	return dm, nil
}

// handleGitPushEvent 处理Git Push事件
func (m *DeployManager) handleGitPushEvent(event events.Event) {
	deployData, ok := event.Data.(map[string]interface{})
	if !ok {
		fmt.Println("无效的部署数据格式")
		return
	}

	repo, ok := deployData["repository"].(string)
	if !ok {
		fmt.Println("缺少仓库信息")
		return
	}

	// 执行部署操作
	if err := m.DeployProject(repo); err != nil {
		fmt.Printf("错误: 仓库 %s 部署失败: %v\n", repo, err)

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

	fmt.Printf("仓库 %s 部署成功完成\n", repo)
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
	fmt.Printf("创建临时工作目录: %s\n", workDir)

	if err := os.MkdirAll(workDir, 0755); err != nil {
		fmt.Printf("创建工作目录失败: %v\n", err)
		return fmt.Errorf("创建工作目录失败: %v", err)
	}
	defer func() {
		fmt.Printf("清理临时工作目录: %s\n", workDir)
		os.RemoveAll(workDir)
	}()

	// 拉取代码
	fmt.Printf("开始从仓库拉取代码: %s\n", repoURL)
	if err := m.gitClone(repoURL, workDir); err != nil {
		fmt.Printf("拉取代码失败: %v\n", err)
		return fmt.Errorf("拉取代码失败: %v", err)
	}

	// 检测项目类型
	projectType := utils.DefaultProjectUtil.DetectProjectType(workDir)
	fmt.Printf("检测到项目类型: %s\n", projectType)

	if projectType == "unknown" {
		fmt.Printf("未检测到项目类型，部署失败\n")
		return fmt.Errorf("未检测到项目类型，部署失败")
	}

	// 根据项目类型选择合适的部署器
	deployer := m.getDeployer(projectType)
	if deployer == nil {
		fmt.Printf("未找到合适的部署器\n")
		return fmt.Errorf("未找到合适的部署器")
	}

	fmt.Printf("使用部署器: %s\n", deployer.GetName())

	// 执行部署
	if err := deployer.Deploy(projectName, workDir, targetDir); err != nil {
		fmt.Printf("部署失败: %v\n", err)
		return fmt.Errorf("部署失败: %v", err)
	}

	fmt.Println("部署成功")

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
	fmt.Printf("规范化仓库地址: %s -> %s\n", originalRepo, repo)

	// 检查工作目录
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		fmt.Printf("工作目录不存在，创建: %s\n", workDir)
		if err := os.MkdirAll(workDir, 0755); err != nil {
			fmt.Printf("创建工作目录失败: %v\n", err)
			return fmt.Errorf("创建工作目录失败: %v", err)
		}
	}

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			fmt.Printf("第 %d 次重试克隆仓库...\n", i+1)
			time.Sleep(time.Second * time.Duration(i+1))
		}

		fmt.Println("开始获取 GitHub 认证信息...")
		auth, err := m.getGitHubAuth(repo)
		if err != nil {
			lastErr = fmt.Errorf("获取GitHub认证信息失败: %v", err)
			fmt.Printf("认证失败详情: %v\n", lastErr)
			continue
		}

		if auth == nil {
			fmt.Println("获取到的认证信息为空，将尝试无认证克隆")
		} else {
			fmt.Printf("成功获取认证信息 - 用户名: %s, Token长度: %d\n",
				auth.Username, len(auth.Password))
		}

		fmt.Printf("开始克隆仓库 %s 到 %s\n", repo, workDir)
		err = m.gitUtil.CloneRepo(repo, "main", workDir, auth)
		if err == nil {
			fmt.Printf("仓库克隆成功: %s\n", repo)
			// 验证克隆结果
			if files, err := os.ReadDir(workDir); err == nil {
				fmt.Printf("克隆目录内容: %d 个文件/目录\n", len(files))
			}
			return nil
		}

		lastErr = err
		fmt.Printf("克隆失败 (尝试 %d/%d): %v\n", i+1, maxRetries, err)
	}

	fmt.Printf("克隆仓库失败（已重试%d次）- 最后错误: %v\n", maxRetries, lastErr)
	return fmt.Errorf("克隆仓库失败（已重试%d次）- 最后错误: %v", maxRetries, lastErr)
}

// getGitHubAuth 获取GitHub认证信息
func (m *DeployManager) getGitHubAuth(repo string) (*githttp.BasicAuth, error) {
	if m.github == nil {
		fmt.Println("GitHub集成未初始化")
		return nil, fmt.Errorf("GitHub集成未初始化")
	}

	fmt.Printf("准备获取仓库认证令牌: %s\n", repo)

	// 检查仓库格式
	repoName := repo
	if strings.HasPrefix(repo, "https://github.com/") {
		repoName = strings.TrimPrefix(repo, "https://github.com/")
	}
	fmt.Printf("处理后的仓库名称: %s\n", repoName)

	// 验证仓库名称格式
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		fmt.Printf("无效的仓库名称格式: %s，应为 'owner/repo' 格式\n", repoName)
		return nil, fmt.Errorf("无效的仓库名称格式: %s，应为 'owner/repo' 格式", repoName)
	}
	fmt.Printf("仓库所有者: %s, 仓库名称: %s\n", parts[0], parts[1])

	token, err := m.github.GetInstallationToken(repoName)
	if err != nil {
		fmt.Printf("获取安装令牌失败: %v\n", err)
		return nil, fmt.Errorf("获取安装令牌失败: %v", err)
	}

	if token == "" {
		fmt.Println("获取到的token为空")
		return nil, fmt.Errorf("获取到的token为空")
	}
	fmt.Printf("成功获取安装令牌 (长度: %d)\n", len(token))

	auth := &githttp.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}

	// 验证认证信息完整性
	if auth.Username == "" || auth.Password == "" {
		fmt.Printf("认证信息不完整: username=%v, token_length=%d\n",
			auth.Username != "", len(auth.Password))
		return nil, fmt.Errorf("认证信息不完整: username=%v, token_length=%d",
			auth.Username != "", len(auth.Password))
	}

	fmt.Println("认证信息构建成功")
	return auth, nil
}

// AddDeployer 添加新的部署器
func (m *DeployManager) AddDeployer(deployer contract.SuperDeployer) {
	fmt.Printf("添加新的部署器: %T\n", deployer)
	m.deployers = append(m.deployers, deployer)
}

// RemoveDeployer 移除指定类型的部署器
func (m *DeployManager) RemoveDeployer(deployerType string) {
	fmt.Printf("移除部署器: %s\n", deployerType)
	newDeployers := make([]contract.SuperDeployer, 0)
	for _, d := range m.deployers {
		if fmt.Sprintf("%T", d) != deployerType {
			newDeployers = append(newDeployers, d)
		}
	}
	m.deployers = newDeployers
}

// GetDeployers 获取所有部署器
func (m *DeployManager) GetDeployers() []contract.SuperDeployer {
	return m.deployers
}

// ClearDeployers 清空所有部署器
func (m *DeployManager) ClearDeployers() {
	fmt.Println("清空所有部署器")
	m.deployers = make([]contract.SuperDeployer, 0)
}

// getDeployer 根据项目类型选择合适的部署器
func (m *DeployManager) getDeployer(projectType string) contract.SuperDeployer {
	for _, deployer := range m.deployers {
		if deployer.GetName() == projectType {
			return deployer
		}
	}
	return nil
}
