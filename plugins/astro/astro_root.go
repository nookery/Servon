package astro

import (
	"fmt"
	"path/filepath"
	"servon/core"
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

func (d *AstroDeployer) GetName() string {
	return "astro"
}

func (d *AstroDeployer) Deploy(projectName string, workDir string, targetDir string) error {
	fmt.Println("开始部署 Astro 项目，工作目录：" + workDir)
	fmt.Println("开始部署 Astro 项目，目标目录：" + targetDir)
	fmt.Println("开始部署 Astro 项目，项目名称：" + projectName)

	err := d.Build(workDir)
	if err != nil {
		fmt.Printf("构建失败: %v\n", err)
		return fmt.Errorf("构建失败: %v", err)
	}

	// 获取工作目录的名字
	workDirName := filepath.Base(workDir)

	// 计算 current 目录，将来会被软链接
	currentDir := targetDir + "/" + workDirName

	// 软链接
	currentLink := targetDir + "/current"

	// 如果项目目录下的软链接存在，则删除
	err = d.RemoveFileOrDir(currentLink)
	if err != nil {
		fmt.Printf("删除 current 目录失败: %v\n", err)
		return fmt.Errorf("删除 current 目录失败: %v", err)
	}

	// 将构建好的项目复制到项目目录下
	err = d.CopyDir(workDir, currentDir)
	if err != nil {
		fmt.Printf("复制项目失败: %v\n", err)
		return fmt.Errorf("复制项目失败: %v", err)
	}

	// 将构建好的项目软链接到项目目录下的 current 目录
	err = d.SymlinkForce(currentDir, currentLink)
	if err != nil {
		fmt.Printf("创建软链接失败: %v\n", err)
		return fmt.Errorf("创建软链接失败: %v", err)
	}

	// 设置Host
	host := DefaultHost
	port := DefaultPort
	serviceFilePath := ""

	// 检查服务配置文件是否存在，不存在则需要创建
	if !d.ServiceManager.HasServiceConf(projectName) {
		serviceFilePath, err = d.AddBackgroundService(projectName, "node", []string{currentLink + "/dist/server/entry.mjs"}, []string{
			fmt.Sprintf("HOST=%s", host),
			fmt.Sprintf("PORT=%d", port),
		})
		if err != nil {
			fmt.Printf("添加后台服务失败: %v\n", err)
			return fmt.Errorf("添加后台服务失败: %v", err)
		}
	} else {
		serviceFilePath = d.GetServiceFilePath(projectName)
	}

	// 成功提示
	fmt.Println()
	fmt.Println("✨ Astro项目部署成功！")
	fmt.Println()
	fmt.Printf("📦 工作目录: %s\n", workDir)
	fmt.Printf("📦 目标目录: %s\n", targetDir)
	fmt.Printf("📁 current（软链接） 路径: %s\n", currentLink)
	fmt.Printf("📁 服务文件路径: %s\n", serviceFilePath)
	fmt.Printf("🌐 服务端口: %d\n", port)
	fmt.Printf("🌐 服务Host: %s\n", host)
	fmt.Printf("🌐 快速打开: http://%s:%d\n", host, port)
	fmt.Println()
	return nil
}

func (d *AstroDeployer) Build(workDir string) error {
	fmt.Println("开始构建 Astro 项目，工作目录：" + workDir)
	// 确保保存路径存在
	if err := d.MakeDir(workDir); err != nil {
		fmt.Printf("创建工作目录失败: %v\n", err)
		return fmt.Errorf("创建工作目录失败: %v", err)
	}

	// pnpm install
	fmt.Println("开始安装 pnpm 依赖")
	err, output := d.RunShellInFolder(workDir, "pnpm", "install")
	fmt.Println(output)
	if err != nil {
		fmt.Printf("pnpm install 失败: %v\n", err)
		return fmt.Errorf("pnpm install 失败: %v", err)
	}

	fmt.Println("pnpm install 成功")

	// pnpm build
	fmt.Println("开始构建 Astro 项目")
	err, output = d.RunShellInFolder(workDir, "pnpm", "build")
	fmt.Println(output)
	if err != nil {
		fmt.Printf("pnpm build 失败: %v\n", err)
		return fmt.Errorf("pnpm build 失败: %v", err)
	}

	fmt.Println("pnpm build 成功")

	return nil
}
