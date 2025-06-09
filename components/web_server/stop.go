package web_server

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Stop 停止服务器
func (ws *WebServer) Stop() error {
	if ws.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return ws.server.Shutdown(ctx)
	}
	return nil
}

// StopBackground 停止后台运行的服务器 - 基于PID文件
func (ws *WebServer) StopBackground() error {
	var stopped bool

	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🛑 开始停止服务器进程...")
	}

	// 通过PID文件停止进程
	pidFile := PID_FILE
	if ws.config.Verbose && ws.logger != nil {
		ws.logger.Infof("🗂️  检查PID文件: %s", pidFile)
	}

	if _, err := os.Stat(pidFile); err == nil {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("📄 找到PID文件，正在读取...")
		}

		if content, err := os.ReadFile(pidFile); err == nil {
			if pidStr := strings.TrimSpace(string(content)); pidStr != "" {
				if pid, err := strconv.Atoi(pidStr); err == nil {
					if ws.config.Verbose && ws.logger != nil {
						ws.logger.Infof("🎯 从PID文件找到进程 %d", pid)
						ws.logger.Infof("📤 发送 SIGTERM 信号给进程 %d", pid)
					} else {
						fmt.Printf("从PID文件找到进程 %d\n", pid)
					}

					// 发送SIGTERM信号
					if err := syscall.Kill(pid, syscall.SIGTERM); err == nil {
						if ws.config.Verbose && ws.logger != nil {
							ws.logger.Infof("✅ 已发送终止信号给进程 %d", pid)
							ws.logger.Infof("⏳ 等待进程优雅退出...")
						} else {
							fmt.Printf("已发送终止信号给进程 %d\n", pid)
						}

						// 等待进程终止
						for i := 0; i < 10; i++ {
							time.Sleep(500 * time.Millisecond)
							if ws.config.Verbose && ws.logger != nil {
								ws.logger.Infof("🔍 检查进程是否已退出 (%d/10)...", i+1)
							}
							// 检查进程是否还存在
							if err := syscall.Kill(pid, 0); err != nil {
								if ws.config.Verbose && ws.logger != nil {
									ws.logger.Infof("✅ 进程已优雅退出")
								}
								stopped = true
								break
							}
						}

						// 如果SIGTERM无效，使用SIGKILL
						if !stopped {
							if ws.config.Verbose && ws.logger != nil {
								ws.logger.Infof("⚠️  进程未响应 SIGTERM，使用 SIGKILL 强制终止")
								ws.logger.Infof("💥 强制终止进程 %d", pid)
							} else {
								fmt.Printf("强制终止进程 %d\n", pid)
							}
							syscall.Kill(pid, syscall.SIGKILL)
							stopped = true
						}
					} else {
						if ws.config.Verbose && ws.logger != nil {
							ws.logger.Infof("⚠️  无法发送信号给进程 %d: %v", pid, err)
						} else {
							fmt.Printf("无法发送信号给进程 %d: %v\n", pid, err)
						}
					}
				} else {
					if ws.config.Verbose && ws.logger != nil {
						ws.logger.Infof("⚠️  PID文件内容无效: %s", pidStr)
					}
				}
			} else {
				if ws.config.Verbose && ws.logger != nil {
					ws.logger.Infof("⚠️  PID文件为空")
				}
			}
		} else {
			if ws.config.Verbose && ws.logger != nil {
				ws.logger.Infof("⚠️  无法读取PID文件: %v", err)
			}
		}

		// 删除PID文件
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("🗑️  删除PID文件: %s", pidFile)
		}
		os.Remove(pidFile)
	} else {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("❌ PID文件不存在")
		} else {
			fmt.Println("PID文件不存在，服务器可能未在运行")
		}
	}

	if stopped {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("✅ 服务器已成功关闭")
		} else {
			fmt.Println("服务器已关闭")
		}
	} else {
		if ws.config.Verbose && ws.logger != nil {
			ws.logger.Infof("ℹ️  服务器未在运行")
		} else {
			fmt.Println("服务器未在运行")
		}
	}

	return nil
}
