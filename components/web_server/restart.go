package web_server

import (
	"os"
	"time"
)

// RestartBackground 重启后台运行的服务器
func (ws *WebServer) RestartBackground() error {
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🔄 开始重启服务器...")
	}

	// 先停止服务器
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🛑 正在停止当前服务器...")
	}
	if err := ws.StopBackground(); err != nil {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Warnf("⚠️  停止服务器时出错: %v", err)
		}
	}

	// 确保PID文件被删除（防止重启时检查失败）
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🗑️  确保PID文件已清理")
	}
	os.Remove(PID_FILE)

	// 等待一小段时间确保资源释放
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("⏳ 等待资源释放...")
	}
	time.Sleep(time.Second)

	// 重新启动服务器（跳过端口检查）
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🚀 正在启动服务器...")
	}
	if err := ws.RunInBackgroundWithOptions(true); err != nil {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Errorf("❌ 启动服务器失败: %v", err)
		}
		return err
	}

	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("✅ 服务器重启成功 -> http://%s:%d", ws.config.Host, ws.config.Port)
	}

	return nil
}
