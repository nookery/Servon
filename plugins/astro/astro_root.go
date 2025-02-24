package astro

import (
	"fmt"
	"os"
	"path/filepath"
	"servon/core"
	"strings"
	"time"
)

type AstroPlugin struct {
	*core.App
}

func Setup(app *core.App) {
	deployer := NewAstroDeployer(app)

	// 添加 Astro 部署器到部署管理器
	app.AddDeployer(deployer)
}

const DefaultPort = 8080
const DefaultBranch = "main"
const DefaultHost = "0.0.0.0"

type AstroDeployer struct {
	*core.App
}

func NewAstroDeployer(app *core.App) *AstroDeployer {
	return &AstroDeployer{
		App: app,
	}
}

// deploy 部署 Astro 项目
func (a *AstroDeployer) deploy(repo string, branch string, host string, port int, logger *core.LogUtil) error {
	logger.Info("开始部署 Astro 项目")

	projectFolder := a.DataManager.GetProjectsRootFolder() + "/" + getProjectNameFromRepo(repo)
	targetFolder := projectFolder + "/" + time.Now().Format("20060102150405")

	err := a.GitClone(repo, branch, targetFolder)
	if err != nil {
		return logger.LogAndReturnErrorf("拉取代码失败: %v", err)
	}

	// 判断是不是 Astro 项目
	if !isAstroProject(targetFolder) {
		return logger.LogAndReturnErrorf("项目不是 Astro 项目")
	}

	err = a.build(targetFolder)
	if err != nil {
		return logger.LogAndReturnErrorf("构建失败: %v", err)
	}

	// 计算 current 目录
	currentFolder := projectFolder + "/current"

	// 如果项目目录下的 current 目录存在，则删除
	if _, err := os.Stat(currentFolder); err == nil {
		err = os.Remove(currentFolder)
		if err != nil {
			return logger.LogAndReturnErrorf("删除 current 目录失败: %v", err)
		}
	}

	// 将构建好的项目软链接到项目目录下的 current 目录
	err = os.Symlink(targetFolder, currentFolder)
	if err != nil {
		return logger.LogAndReturnErrorf("创建软链接失败: %v", err)
	}

	// 设置Host
	if host == "" {
		host = DefaultHost
	}

	// 设置端口
	if port == 0 {
		port = DefaultPort // Astro 的默认端口
	}

	serviceFilePath := ""

	// 检查服务配置文件是否存在，不存在则需要创建
	if !a.ServiceManager.HasServiceConf(getProjectNameFromRepo(repo)) {
		serviceFilePath, err = a.AddBackgroundService(getProjectNameFromRepo(repo), "node", []string{currentFolder + "/dist/server/entry.mjs"}, []string{
			fmt.Sprintf("HOST=%s", host),
			fmt.Sprintf("PORT=%d", port),
		})
		if err != nil {
			return logger.LogAndReturnErrorf("添加背景服务失败: %v", err)
		}
	} else {
		serviceFilePath = a.GetServiceFilePath(getProjectNameFromRepo(repo))
	}

	// 成功提示
	fmt.Println()
	logger.Info("✨ Astro项目部署成功！")
	fmt.Println()
	logger.Infof("📦 仓库地址: %s", repo)
	logger.Infof("📦 分支: %s", branch)
	logger.Infof("📁 项目路径: %s", projectFolder)
	logger.Infof("📁 目标路径: %s", targetFolder)
	logger.Infof("📁 current（软链接） 路径: %s", currentFolder)
	logger.Infof("📁 服务文件路径: %s", serviceFilePath)
	logger.Infof("🌐 服务端口: %d", port)
	logger.Infof("🌐 服务Host: %s", host)
	logger.Infof("🌐 快速打开: http://%s:%d", host, port)
	fmt.Println()
	return nil
}

// isAstroProject 判断是否是 Astro 项目
func isAstroProject(projectFolder string) bool {
	if _, err := os.Stat(projectFolder + "/astro.config.mjs"); os.IsNotExist(err) {
		return false
	}

	return true
}

// getProjectNameFromRepo 从仓库地址中获取项目名称
// 比如：https://github.com/user/project.git 返回 project
// 比如：git@github.com:user/project.git 返回 project
// 比如：ssh://git@github.com/user/project.git 返回 project
// 比如：git+ssh://git@github.com/user/project.git 返回 project
// 比如：git+https://github.com/user/project.git 返回 project
// 比如：git+http://github.com/user/project.git 返回 project
// 如果不能获取到项目名称，则返回随机字符串（根据当前时间生成）
func getProjectNameFromRepo(repo string) string {
	repo = strings.TrimSuffix(repo, ".git")
	repo = strings.TrimPrefix(repo, "https://")
	repo = strings.TrimPrefix(repo, "http://")
	repo = strings.TrimPrefix(repo, "git@")
	repo = strings.TrimPrefix(repo, "ssh://")
	repo = strings.TrimPrefix(repo, "git+")
	repo = strings.TrimPrefix(repo, "git+ssh://")

	parts := strings.Split(repo, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return time.Now().Format("20060102150405")
}

func (a *AstroDeployer) build(path string) error {
	// 确保保存路径存在
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// pnpm install
	if err := a.RunShellInFolder(path, "pnpm", "install"); err != nil {
		return err
	}

	a.Info("pnpm install 成功")

	// pnpm build
	if err := a.RunShellInFolder(path, "pnpm", "build"); err != nil {
		return err
	}

	a.Info("pnpm build 成功")

	return nil
}

func (d *AstroDeployer) GetName() string {
	return "astro"
}
func (d *AstroDeployer) CanHandle(workDir string) bool {
	// 检查是否存在 astro.config.mjs 文件
	configPath := filepath.Join(workDir, "astro.config.mjs")
	return d.FileUtil.IsFileExists(configPath)
}

func (d *AstroDeployer) Deploy(workDir string, targetDir string, logger *core.LogUtil) error {
	logger.Info("开始部署 Astro 项目")
	// 使用现有的 deploy 函数，但需要调整参数
	repo := filepath.Base(workDir) // 使用目录名作为项目名
	return d.deploy(repo, DefaultBranch, DefaultHost, DefaultPort, logger)
}
func (d *AstroDeployer) Build(workDir string, logger *core.LogUtil) error {
	return d.build(workDir)
}
