package github_runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"servon/core"
)

// ProgressReader 用于跟踪读取进度
type ProgressReader struct {
	io.Reader
	Total      int64
	Current    int64
	OnProgress func(current, total int64)
}

// Read 实现了 io.Reader 接口
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)
	if pr.OnProgress != nil {
		pr.OnProgress(pr.Current, pr.Total)
	}
	return n, err
}

func Setup(app *core.App) {
	plugin := NewGitHubRunner(app)
	app.RegisterSoftware("github-runner", plugin)
}

type GitHubRunner struct {
	info      core.SoftwareInfo
	targetDir string
	*core.App
}

func NewGitHubRunner(app *core.App) core.SuperSoft {
	return &GitHubRunner{
		info: core.SoftwareInfo{
			Name:        "github-runner",
			Description: "GitHub Actions self-hosted runner",
		},
		targetDir: app.GetSoftwareRootFolder("github-runner"),
		App:       app,
	}
}

// Uninstall 卸载 GitHub Runner
func (g *GitHubRunner) Uninstall() error {
	// 运行卸载脚本
	g.Infof("正在卸载 GitHub Runner...")
	if _, err := os.Stat(g.targetDir + "/config.sh"); !os.IsNotExist(err) {
		err, _ = g.RunShell(g.targetDir+"/config.sh", "remove", "--unattended")
		if err != nil {
			return fmt.Errorf("卸载失败:\n%s", err.Error())
		}
	}

	// 删除安装目录
	if err := os.RemoveAll(g.targetDir); err != nil {
		return fmt.Errorf("删除安装目录失败: %s", err)
	}

	g.Success("GitHub Runner 卸载完成")

	return nil
}

func (g *GitHubRunner) GetStatus() (map[string]string, error) {
	g.Infof("获取 GitHub Runner 状态...")
	status := "not_installed"

	// 检查安装目录是否存在
	if _, err := os.Stat(g.targetDir + "/run.sh"); !os.IsNotExist(err) {
		g.Infof("GitHub Runner 目录存在")
		status = "stopped"

		// 检查进程是否运行
		cmd := exec.Command("pgrep", "-f", "run.sh")
		if err := cmd.Run(); err == nil {
			status = "running"
		}
	}

	g.Infof("GitHub Runner 状态: %s", status)

	return map[string]string{
		"status":  status,
		"version": "unknown",
	}, nil
}

func (g *GitHubRunner) GetInfo() core.SoftwareInfo {
	return g.info
}

// Start 启动 GitHub Runner
func (g *GitHubRunner) Start() error {
	// 检查是否为 root 用户
	if os.Geteuid() != 0 {
		err := fmt.Errorf("启动 GitHub Runner 需要 root 权限")
		g.Errorf("%s", err.Error())
		return err
	}

	// 检查运行状态
	status, err := g.GetStatus()
	if err != nil {
		return fmt.Errorf("检查运行状态失败: %s", err)
	}
	if status["status"] == "running" {
		return fmt.Errorf("GitHub Runner 已经在运行中")
	}

	// 检查是否已配置
	if _, err := os.Stat(g.targetDir + "/.runner"); os.IsNotExist(err) {
		g.Infof("Runner 未配置，请输入以下信息：")

		var url, token string
		g.Infof("请输入 GitHub 仓库/组织 URL：")
		fmt.Scanln(&url)

		g.Infof("请输入 Runner 注册令牌：")
		fmt.Scanln(&token)

		if url == "" || token == "" {
			return fmt.Errorf("URL 和令牌不能为空")
		}

		// 以 github-runner 用户身份配置 runner
		err, _ = g.RunShell("su", "-", "github-runner", "-c", fmt.Sprintf("%s/config.sh --url %s --token %s --unattended", g.targetDir, url, token))
		if err != nil {
			return g.LogAndReturnErrorf("配置 runner 失败: %s", err.Error())
		}
	}

	// 以 github-runner 用户身份启动 runner
	err, _ = g.RunShell("su", "-", "github-runner", "-c", fmt.Sprintf("cd %s && nohup ./run.sh > runner.log 2>&1 &", g.targetDir))
	if err != nil {
		return fmt.Errorf("启动失败: %s", err)
	}

	g.Success("GitHub Runner 已启动")
	return nil
}

func (g *GitHubRunner) Stop() error {
	cmd := exec.Command("pkill", "-f", "run.sh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止 runner 失败: %s", err)
	}
	return nil
}

func (g *GitHubRunner) Reload() error {
	if err := g.Stop(); err != nil {
		return err
	}
	return g.Start()
}
