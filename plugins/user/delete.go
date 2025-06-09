package user

import (
	"fmt"
	"servon/components/user"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// deleteCmd 删除用户
var deleteCmd = &cobra.Command{
	Use:   "delete [username]",
	Short: "删除用户",
	Long:  color.Success.Render("\r\n删除指定的系统用户，可选择是否删除用户主目录"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		removeHome, _ := cmd.Flags().GetBool("remove-home")
		force, _ := cmd.Flags().GetBool("force")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// 验证用户名
		if username == "" {
			return fmt.Errorf("用户名不能为空")
		}

		// 防止删除重要系统用户
		if isProtectedUser(username) && !force {
			return fmt.Errorf("用户 %s 是受保护的系统用户，使用 --force 参数强制删除", username)
		}

		userManager := user.NewUserManager()

		// 检查用户是否存在
		exists, err := userManager.UserExists(username)
		if err != nil {
			return fmt.Errorf("检查用户是否存在失败: %v", err)
		}
		if !exists {
			return fmt.Errorf("用户 %s 不存在", username)
		}

		// 获取用户信息（用于显示）
		if verbose {
			users, err := userManager.GetUserList()
			if err == nil {
				for _, u := range users {
					if u.Username == username {
						color.Info.Printf("📋 用户信息:\n")
						color.Gray.Printf("   用户名: %s\n", u.Username)
						color.Gray.Printf("   主目录: %s\n", u.HomeDir)
						color.Gray.Printf("   Shell: %s\n", u.Shell)
						color.Gray.Printf("   用户组: %s\n", strings.Join(u.Groups, ", "))
						break
					}
				}
			}
		}

		// 确认删除
		if !force {
			color.Warn.Printf("⚠️  确定要删除用户 %s 吗？", username)
			if removeHome {
				color.Warn.Print(" (包括主目录)")
			}
			color.Warn.Print(" [y/N]: ")

			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
				color.Info.Println("❌ 操作已取消")
				return nil
			}
		}

		color.Info.Printf("🗑️  正在删除用户 %s...\n", username)

		// 删除用户
		var deleteArgs []string
		if removeHome {
			deleteArgs = []string{"-r", username}
		} else {
			deleteArgs = []string{username}
		}

		err, output := user.RunShell("userdel", deleteArgs...)
		if err != nil {
			return fmt.Errorf("删除用户失败: %v\n输出: %s", err, output)
		}

		color.Success.Printf("✅ 用户 %s 删除成功", username)
		if removeHome {
			color.Success.Print(" (包括主目录)")
		}
		color.Success.Print("\n")

		// 检查是否还有相关进程
		if verbose {
			checkUserProcesses(username)
		}

		return nil
	},
}

func init() {
	// 添加命令行参数
	deleteCmd.Flags().BoolP("remove-home", "r", false, "同时删除用户主目录")
	deleteCmd.Flags().BoolP("force", "f", false, "强制删除（跳过确认和保护检查）")
	deleteCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
}

// isProtectedUser 检查是否为受保护的用户
func isProtectedUser(username string) bool {
	// 受保护的系统用户列表
	protectedUsers := []string{
		"root", "daemon", "bin", "sys", "sync", "games", "man", "lp",
		"mail", "news", "uucp", "proxy", "www-data", "backup", "list",
		"irc", "gnats", "nobody", "systemd-network", "systemd-resolve",
		"syslog", "messagebus", "_apt", "lxd", "uuidd", "dnsmasq",
		"landscape", "pollinate", "sshd", "mysql", "redis", "postgres",
		"nginx", "apache", "docker", "git", "jenkins", "mongodb",
	}

	for _, protectedUser := range protectedUsers {
		if username == protectedUser {
			return true
		}
	}

	// 检查用户名是否以下划线开头（通常是系统用户）
	if strings.HasPrefix(username, "_") {
		return true
	}

	return false
}

// checkUserProcesses 检查用户是否还有运行的进程
func checkUserProcesses(username string) {
	err, output := user.RunShell("ps", "-u", username)
	if err == nil && strings.TrimSpace(output) != "" {
		color.Warn.Printf("⚠️  警告: 用户 %s 可能还有运行的进程:\n", username)
		color.Gray.Println(output)
		color.Info.Println("💡 建议使用 'pkill -u username' 终止用户进程")
	} else {
		color.Info.Printf("✅ 用户 %s 没有运行的进程\n", username)
	}
}
