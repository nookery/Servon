// Package components 提供了 Servon 项目的独立功能组件集合
//
// 设计理念：
// components 包遵循组件化架构设计，每个子包都是一个相对独立的功能模块。
// 这种设计有助于代码的模块化、可维护性和可重用性。
//
// 设计特点：
// 1. 独立性：每个组件尽可能减少对其他模块的依赖
// 2. 可重用性：组件可以在不同的上下文中使用
// 3. 统一入口：通过 components 包提供统一的组件访问入口
//
// 当前包含的组件：
// - events: 事件总线系统，提供发布-订阅模式的事件处理机制
// - env_manager: 环境变量管理器，提供环境变量的读取和管理功能
// - github: GitHub 集成组件，提供 GitHub API 交互和 Webhook 处理功能
// - log_util: 日志工具组件，提供统一的日志记录和管理功能
// - command_util: 命令行工具组件，提供命令执行和选项管理功能
// - shell_util: 提供Shell命令执行功能
// - file_util: 提供文件和目录操作功能
// - web_server_util: 提供Web服务器功能
// - string_util: 提供字符串处理功能
// - project_util: 提供项目检测和管理功能
// - git_util: 提供Git操作功能
// - process_util: 提供进程管理功能
// - package_util: 提供包管理功能
// - dev_util: 提供开发环境检测功能
//
// 使用示例：
//
//	import (
//		"servon/components/log_util"
//		"servon/components/command_util"
//		"servon/components/shell_util"
//		"servon/components/file_util"
//		"servon/components/web_server_util"
//		"servon/components/string_util"
//		"servon/components/project_util"
//		"servon/components/git_util"
//		"servon/components/process_util"
//		"servon/components/package_util"
//		"servon/components/dev_util"
//	)
//
//	// 使用日志工具
//	logger := log_util.NewLogUtil("MyApp")
//	logger.Info("应用启动")
//
//	// 使用命令工具
//	cmd := command_util.NewCommand(command_util.CommandOptions{
//		Use:   "test",
//		Short: "测试命令",
//	})
//
//	// 使用Shell工具
//	shell := &shell_util.ShellUtil{}
//	output, err := shell.RunShellWithOutput("ls -la")
//
//	// 使用文件工具
//	fileUtil := &file_util.FileUtil{}
//	exists := fileUtil.FileExists("/path/to/file")
//
//	// 使用Web服务器
//	server := web_server_util.NewWebServer(web_server_util.WebServerConfig{
//		Host: "localhost",
//		Port: 8080,
//	})
//
//	// 使用字符串工具
//	stringUtil := &string_util.StringUtil{}
//	projectName := stringUtil.GetProjectNameFromRepoURL("https://github.com/user/repo.git")
//
//	// 使用项目工具
//	projectUtil := project_util.NewProjectUtil()
//	projectType := projectUtil.DetectProjectType("/path/to/project")
//
//	// 使用Git工具
//	gitUtil := git_util.NewGitUtil(logger)
//	err = gitUtil.CloneRepo("https://github.com/user/repo.git", "main", "/path/to/clone", nil)
//
//	// 使用进程工具
//	processUtil := process_util.NewProcessUtil()
//	err = processUtil.AutoStopPortProcess(8080)
//
//	// 使用包工具
//	version := package_util.ReadPackageVersion()
//
//	// 使用开发环境检测
//	devUtil := &dev_util.DevUtil{}
//	isDev := devUtil.IsDev()
//
// 添加新组件的指导原则：
// 1. 确保组件具有明确的职责边界
// 2. 最小化外部依赖，特别是对 core 模块的依赖
// 3. 提供清晰的 API 接口
// 4. 包含必要的文档和使用示例
// 5. 遵循项目的代码规范和命名约定
package components

import "servon/components/events"

// EventBus 全局事件总线实例
// 提供系统级的事件发布-订阅和请求-响应功能
var EventBus = events.GetEventBusInstance()
