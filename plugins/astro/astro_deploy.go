package astro

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const DefaultPort = 8080
const DefaultBranch = "main"
const DefaultHost = "0.0.0.0"

// deploy 部署 Astro 项目
func (a *AstroPlugin) deploy(repo string, branch string, host string, port int) error {
	a.Info("开始部署 Astro 项目")

	projectFolder := a.DataManager.GetProjectsRootFolder() + "/" + getProjectNameFromRepo(repo)
	targetFolder := projectFolder + "/" + time.Now().Format("20060102150405")

	err := a.GitClone(repo, branch, targetFolder)
	if err != nil {
		return err
	}

	// 判断是不是 Astro 项目
	if !isAstroProject(targetFolder) {
		return fmt.Errorf("项目不是 Astro 项目")
	}

	err = a.build(targetFolder)
	if err != nil {
		return err
	}

	// 计算 current 目录
	currentFolder := projectFolder + "/current"

	// 如果项目目录下的 current 目录存在，则删除
	if _, err := os.Stat(currentFolder); err == nil {
		err = os.Remove(currentFolder)
		if err != nil {
			return err
		}
	}

	// 将构建好的项目软链接到项目目录下的 current 目录
	err = os.Symlink(targetFolder, currentFolder)
	if err != nil {
		return err
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
			return err
		}
	} else {
		serviceFilePath = a.GetServiceFilePath(getProjectNameFromRepo(repo))
	}

	// 成功提示
	fmt.Println()
	color.New(color.FgGreen, color.Bold).Printf("✨ Astro项目部署成功！\n")
	fmt.Println()
	color.New(color.FgWhite).Print("📦 仓库地址: ")
	color.New(color.FgHiWhite).Printf("%s\n", repo)
	color.New(color.FgWhite).Print("📦 分支: ")
	color.New(color.FgHiWhite).Printf("%s\n", branch)
	color.New(color.FgWhite).Print("📁 项目路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", projectFolder)
	color.New(color.FgWhite).Print("📁 目标路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", targetFolder)
	color.New(color.FgWhite).Print("📁 current（软链接） 路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", currentFolder)
	color.New(color.FgWhite).Print("📁 服务文件路径: ")
	color.New(color.FgHiWhite).Printf("%s\n", serviceFilePath)
	color.New(color.FgWhite).Print("🌐 服务端口: ")
	color.New(color.FgHiWhite).Printf("%d\n", port)
	color.New(color.FgWhite).Print("🌐 服务Host: ")
	color.New(color.FgHiWhite).Printf("%s\n", host)
	color.New(color.FgWhite).Print("🌐 快速打开: ")
	color.New(color.FgHiWhite).Printf("http://%s:%d\n", host, port)
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

func (a *AstroPlugin) build(path string) error {
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
