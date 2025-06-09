package user

import (
	"fmt"
	"servon/components/user"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// listCmd 列出系统用户
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出系统用户",
	Long:  color.Success.Render("\r\n列出系统中的所有用户，包括用户名、用户组、Shell、主目录等信息"),
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		showSystem, _ := cmd.Flags().GetBool("system")

		userManager := user.NewUserManager()
		users, err := userManager.GetUserList()
		if err != nil {
			return fmt.Errorf("获取用户列表失败: %v", err)
		}

		color.Info.Printf("📋 系统用户列表 (共 %d 个用户)\n\n", len(users))

		for _, u := range users {
			// 过滤系统用户（UID < 1000 的用户通常是系统用户）
			if !showSystem && isSystemUser(u.Username) {
				continue
			}

			// 显示用户基本信息
			color.Primary.Printf("👤 %s", u.Username)
			if u.Sudo {
				color.Warn.Print(" [SUDO]")
			}
			color.Print("\n")

			if verbose {
				// 详细信息
				color.Gray.Printf("   主目录: %s\n", u.HomeDir)
				color.Gray.Printf("   Shell: %s\n", u.Shell)
				color.Gray.Printf("   用户组: %s\n", strings.Join(u.Groups, ", "))
				if !u.CreateTime.IsZero() {
					color.Gray.Printf("   创建时间: %s\n", u.CreateTime.Format("2006-01-02 15:04:05"))
				}
				if !u.LastLogin.IsZero() {
					color.Gray.Printf("   最后登录: %s\n", u.LastLogin.Format("2006-01-02 15:04:05"))
				}
				color.Print("\n")
			}
		}

		return nil
	},
}

func init() {
	// 添加命令行参数
	listCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
	listCmd.Flags().BoolP("system", "s", false, "包含系统用户")
}

// isSystemUser 判断是否为系统用户
// 通常 UID < 1000 的用户是系统用户
func isSystemUser(username string) bool {
	// 常见的系统用户名
	systemUsers := []string{
		"root", "daemon", "bin", "sys", "sync", "games", "man", "lp",
		"mail", "news", "uucp", "proxy", "www-data", "backup", "list",
		"irc", "gnats", "nobody", "systemd-network", "systemd-resolve",
		"syslog", "messagebus", "_apt", "lxd", "uuidd", "dnsmasq",
		"landscape", "pollinate", "sshd", "mysql", "redis", "postgres",
		"nginx", "apache", "docker", "git", "jenkins", "mongodb",
	}

	for _, sysUser := range systemUsers {
		if username == sysUser {
			return true
		}
	}

	// 检查用户名是否以下划线开头（通常是系统用户）
	if strings.HasPrefix(username, "_") {
		return true
	}

	return false
}
