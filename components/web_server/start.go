package web_server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

// start 启动服务器
func (ws *WebServer) start() error {
	addr := fmt.Sprintf("%s:%d", ws.config.Host, ws.config.Port)

	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🚀 正在启动服务器...")
		ws.logger.Infof("📍 监听地址: %s", addr)
		ws.logger.Infof("🌐 主机: %s", ws.config.Host)
		ws.logger.Infof("🔌 端口: %d", ws.config.Port)
	}

	ws.server = &http.Server{
		Addr:    addr,
		Handler: ws.Engine,
	}

	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("✅ HTTP 服务器已配置完成，开始监听...")
	}

	return ws.server.ListenAndServe()
}

// RunInBackground 在后台运行服务器（作为独立进程）
func (ws *WebServer) RunInBackground() error {
	return ws.RunInBackgroundWithOptions(false)
}

// RunInBackgroundWithOptions 在后台运行服务器（支持跳过端口检查）
func (ws *WebServer) RunInBackgroundWithOptions(skipPortCheck bool) error {
	// 如果未设置端口，则使用默认端口
	if ws.config.Port == 0 {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("📍 未设置端口，使用默认端口")
		}
		ws.config.Port = DEFAULT_PORT
	}

	// 输出配置信息
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("📍 主机: %s", ws.config.Host)
		ws.logger.Infof("🔌 端口: %d", ws.config.Port)
	}

	// 检查PID文件是否存在
	if _, err := os.Stat(PID_FILE); err == nil {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("🔍 发现PID文件: %s", PID_FILE)
		}
		// 读取PID文件内容
		if pidData, err := os.ReadFile(PID_FILE); err == nil && len(pidData) > 0 {
			pidStr := string(pidData)
			if ws.config.Verbose && ws.logger != nil {
				ws.logger.Infof("📄 PID文件内容: %s", pidStr)
			}
			return fmt.Errorf("服务器可能已在运行中 (PID文件: %s)\n提示：请先执行 'stop' 命令关闭服务器，或删除PID文件后重试", PID_FILE)
		} else if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("⚠️  PID文件为空或无法读取，将继续启动")
		}
	}

	// 检查服务器是否已经在运行（除非跳过端口检查）
	if !skipPortCheck {
		if pid, err := findProcessByPortWithVerbose(ws.config.Port, ws.config.Verbose, ws.logger); err == nil && pid > 0 {
			return fmt.Errorf("服务器已在运行中 (PID: %d)\n提示：如需重启，请使用 'restart' 命令", pid)
		}
	}

	// 设置守护进程的上下文
	ctx := &daemon.Context{
		PidFileName: PID_FILE,
		PidFilePerm: 0644,
		LogFileName: LOG_FILE,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	// 启动守护进程
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🌐 正在将服务器作为守护进程运行...")
	}
	d, err := ctx.Reborn()
	if err != nil {
		return fmt.Errorf("创建守护进程失败: %v", err)
	}
	if d != nil {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("🌐 服务器已作为守护进程运行，PID: %d", d.Pid)
		}
		return nil // 父进程退出
	}

	// 子进程继续执行
	defer ctx.Release()

	// 启动服务器
	if err := ws.start(); err != nil {
		return err
	}

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	return ws.stop()
}
