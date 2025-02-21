package github

import "servon/core/internal/utils"

const (
	WebhookDir      = "/data/github/webhook"
	installationDir = "/data/github/installations" // GitHub安装数据目录
	configDir       = "/data/github/config"        // GitHub配置目录
	githubLogDir    = "/data/github/logs"          // GitHub集成日志目录
	timeFormat      = "2006-01-02"                 // 日期格式
)

var logger = utils.NewLogUtil(githubLogDir)
