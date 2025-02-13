package github_runner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"servon/core"
	"servon/core/contract"
	"strings"
	"time"
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

func Setup(core *core.Core) {
	plugin := NewGitHubRunner(core)
	core.RegisterSoftware("github-runner", plugin)
}

type GitHubRunner struct {
	info      contract.SoftwareInfo
	targetDir string
	*core.Core
}

func NewGitHubRunner(core *core.Core) contract.SuperSoft {
	return &GitHubRunner{
		info: contract.SoftwareInfo{
			Name:        "github-runner",
			Description: "GitHub Actions self-hosted runner",
		},
		targetDir: core.GetDataRootFolder() + "/actions-runner",
		Core:      core,
	}
}

// Install 安装 GitHub Runner
func (g *GitHubRunner) Install() error {
	// 先检查运行状态
	status, err := g.GetStatus()
	if err != nil {
		return fmt.Errorf("检查运行状态失败: %s", err)
	}
	if status["status"] == "running" {
		return fmt.Errorf("GitHub Runner 正在运行中，请先停止后再安装")
	}

	osType := g.GetOSType()
	g.PrintInfof("检测到操作系统: %s", osType)

	// 创建 runner 专用用户
	g.PrintInfo("创建 GitHub Runner 专用用户...")
	runnerUser := "github-runner"
	runnerPassword := "runner" + fmt.Sprint(time.Now().Unix()) // 生成随机密码

	exists, err := g.UserExists(runnerUser)
	if err != nil {
		return fmt.Errorf("检查用户失败: %s", err)
	}

	if !exists {
		if err := g.CreateUser(runnerUser, runnerPassword); err != nil {
			return fmt.Errorf("创建用户失败: %s", err)
		}
		g.PrintSuccessf("已创建专用用户: %s", runnerUser)
	} else {
		g.PrintInfo("专用用户已存在，跳过创建")
	}

	switch osType {
	case core.Ubuntu, core.Debian, core.CentOS, core.RedHat:
		g.PrintInfo("开始安装 GitHub Runner...")

		// 获取最新版本号
		resp, err := http.Get("https://api.github.com/repos/actions/runner/releases/latest")
		if err != nil {
			return fmt.Errorf("获取最新版本失败: %s", err)
		}
		defer resp.Body.Close()

		// 检查状态码
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusForbidden && strings.Contains(string(body), "API rate limit exceeded") {
				return fmt.Errorf("GitHub API 调用次数超限，请稍后再试或使用 API 令牌。状态码: %d", resp.StatusCode)
			}
			return fmt.Errorf("获取版本信息失败，GitHub API 返回状态码: %d，响应内容: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("读取响应失败: %s", err)
		}

		// 解析版本号
		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.Unmarshal(body, &release); err != nil {
			return fmt.Errorf("解析版本信息失败: %s", err)
		}

		if release.TagName == "" {
			return fmt.Errorf("获取版本号失败：API 返回的版本号为空")
		}

		version := strings.TrimPrefix(release.TagName, "v")
		if version == "" {
			return fmt.Errorf("获取版本号失败：无效的版本号格式")
		}

		g.PrintInfof("最新版本: %s", version)

		// 清理本地文件夹
		g.PrintInfo("清理本地文件夹...")
		err = os.RemoveAll(g.targetDir)
		if err != nil {
			return fmt.Errorf("清理本地文件夹失败: %s", err)
		}

		// 创建目标目录
		g.PrintInfo("创建目标目录...")
		err = os.MkdirAll(g.targetDir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %s", err)
		}

		// 获取系统架构
		arch, err := g.RunShellWithOutput("uname", "-m")
		if err != nil {
			return fmt.Errorf("获取系统架构失败: %s", err)
		}
		arch = strings.TrimSpace(arch)

		// 确定下载的架构版本
		var archInUrl string
		switch arch {
		case "x86_64":
			archInUrl = "x64"
		case "aarch64":
			archInUrl = "arm64"
		default:
			return fmt.Errorf("不支持的系统架构: %s", arch)
		}

		g.PrintInfof("检测到系统架构: %s", arch)
		g.PrintInfo("开始下载 GitHub Runner...")
		downloadUrl := fmt.Sprintf("https://github.com/actions/runner/releases/download/v%s/actions-runner-linux-%s-%s.tar.gz", version, archInUrl, version)

		err = g.Download(downloadUrl, g.targetDir+"/actions-runner-linux.tar.gz")
		if err != nil {
			return fmt.Errorf("下载 runner 失败: %s", err)
		}

		// 解压
		g.PrintInfo("开始解压 GitHub Runner...")
		err = g.RunShell("tar", "xzf", g.targetDir+"/actions-runner-linux.tar.gz", "-C", g.targetDir)
		if err != nil {
			return fmt.Errorf("解压失败: %s", err)
		}

		// 删除压缩包
		g.PrintInfo("删除压缩包，因为已解压到 " + g.targetDir)
		err = os.Remove(g.targetDir + "/actions-runner-linux.tar.gz")
		if err != nil {
			return fmt.Errorf("删除压缩包失败: %s", err)
		}

		// 安装依赖
		g.PrintInfo("开始安装依赖...")
		err = g.RunShell("/bin/bash", g.targetDir+"/bin/installdependencies.sh")
		if err != nil {
			return fmt.Errorf("安装依赖失败: %s", err)
		}

		// 修改目录所有权
		g.PrintInfo("修改目录所有权...")
		if err := g.RunShell("chown", "-R", runnerUser+":"+runnerUser, g.targetDir); err != nil {
			return fmt.Errorf("修改目录所有权失败: %s", err)
		}

		g.PrintSuccess("GitHub Runner 安装完成")
		g.PrintInfof("请使用 %s 用户运行 GitHub Runner", runnerUser)
		return nil

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		g.PrintErrorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}
}

// Uninstall 卸载 GitHub Runner
func (g *GitHubRunner) Uninstall() error {
	// 运行卸载脚本
	g.PrintInfo("正在卸载 GitHub Runner...")
	if _, err := os.Stat(g.targetDir + "/config.sh"); !os.IsNotExist(err) {
		err = g.RunShell(g.targetDir+"/config.sh", "remove", "--unattended")
		if err != nil {
			return fmt.Errorf("卸载失败:\n%s", err.Error())
		}
	}

	// 删除安装目录
	if err := os.RemoveAll(g.targetDir); err != nil {
		return fmt.Errorf("删除安装目录失败: %s", err)
	}

	g.PrintSuccess("GitHub Runner 卸载完成")

	return nil
}

func (g *GitHubRunner) GetStatus() (map[string]string, error) {
	g.PrintInfo("获取 GitHub Runner 状态...")
	status := "not_installed"

	// 检查安装目录是否存在
	if _, err := os.Stat(g.targetDir + "/run.sh"); !os.IsNotExist(err) {
		g.PrintInfo("GitHub Runner 目录存在")
		status = "stopped"

		// 检查进程是否运行
		cmd := exec.Command("pgrep", "-f", "run.sh")
		if err := cmd.Run(); err == nil {
			status = "running"
		}
	}

	g.PrintInfof("GitHub Runner 状态: %s", status)

	return map[string]string{
		"status":  status,
		"version": "unknown",
	}, nil
}

func (g *GitHubRunner) GetInfo() contract.SoftwareInfo {
	return g.info
}

// Start 启动 GitHub Runner
func (g *GitHubRunner) Start() error {
	// 检查是否为 root 用户
	if os.Geteuid() != 0 {
		err := fmt.Errorf("启动 GitHub Runner 需要 root 权限")
		g.PrintErrorf("%s", err.Error())
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
		g.PrintInfo("Runner 未配置，请输入以下信息：")

		var url, token string
		g.PrintInfo("请输入 GitHub 仓库/组织 URL：")
		fmt.Scanln(&url)

		g.PrintInfo("请输入 Runner 注册令牌：")
		fmt.Scanln(&token)

		if url == "" || token == "" {
			return fmt.Errorf("URL 和令牌不能为空")
		}

		// 以 github-runner 用户身份配置 runner
		err = g.RunShell("su", "-", "github-runner", "-c", fmt.Sprintf("%s/config.sh --url %s --token %s --unattended", g.targetDir, url, token))
		if err != nil {
			return g.PrintAndReturnErrorf("配置 runner 失败: %s", err.Error())
		}
	}

	// 以 github-runner 用户身份启动 runner
	err = g.RunShell("su", "-", "github-runner", "-c", fmt.Sprintf("cd %s && nohup ./run.sh > runner.log 2>&1 &", g.targetDir))
	if err != nil {
		return fmt.Errorf("启动失败: %s", err)
	}

	g.PrintSuccess("GitHub Runner 已启动")
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
