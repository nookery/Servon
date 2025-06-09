package user

import (
	"fmt"
	"servon/components/user"
	"strings"
	"syscall"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// createCmd 创建新用户
var createCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "创建新用户",
	Long:  color.Success.Render("\r\n创建新的系统用户，可以指定密码和其他选项"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		password, _ := cmd.Flags().GetString("password")
		shell, _ := cmd.Flags().GetString("shell")
		// homeDir, _ := cmd.Flags().GetString("home") // TODO: 实现自定义主目录功能
		groups, _ := cmd.Flags().GetStringSlice("groups")
		sudo, _ := cmd.Flags().GetBool("sudo")
		force, _ := cmd.Flags().GetBool("force")

		// 验证用户名
		if err := validateUsername(username); err != nil {
			return err
		}

		userManager := user.NewUserManager()

		// 检查用户是否已存在
		exists, err := userManager.UserExists(username)
		if err != nil {
			return fmt.Errorf("检查用户是否存在失败: %v", err)
		}
		if exists && !force {
			return fmt.Errorf("用户 %s 已存在，使用 --force 参数强制覆盖", username)
		}

		// 如果没有提供密码，提示输入
		if password == "" {
			password, err = promptPassword()
			if err != nil {
				return fmt.Errorf("获取密码失败: %v", err)
			}
		}

		color.Info.Printf("🔨 正在创建用户 %s...\n", username)

		// 创建用户
		err = userManager.CreateUser(username, password)
		if err != nil {
			return fmt.Errorf("创建用户失败: %v", err)
		}

		color.Success.Printf("✅ 用户 %s 创建成功\n", username)

		// 设置额外选项
		if shell != "" {
			if err := setUserShell(username, shell); err != nil {
				color.Warn.Printf("⚠️  设置Shell失败: %v\n", err)
			} else {
				color.Info.Printf("🐚 已设置Shell为: %s\n", shell)
			}
		}

		if len(groups) > 0 {
			if err := addUserToGroups(username, groups); err != nil {
				color.Warn.Printf("⚠️  添加到用户组失败: %v\n", err)
			} else {
				color.Info.Printf("👥 已添加到用户组: %s\n", strings.Join(groups, ", "))
			}
		}

		if sudo {
			if err := addUserToSudo(username); err != nil {
				color.Warn.Printf("⚠️  添加sudo权限失败: %v\n", err)
			} else {
				color.Info.Printf("🔐 已添加sudo权限\n")
			}
		}

		return nil
	},
}

func init() {
	// 添加命令行参数
	createCmd.Flags().StringP("password", "p", "", "用户密码（如果不提供将提示输入）")
	createCmd.Flags().StringP("shell", "s", "", "用户Shell（如 /bin/bash）")
	createCmd.Flags().StringP("home", "d", "", "用户主目录")
	createCmd.Flags().StringSliceP("groups", "g", []string{}, "添加到的用户组列表")
	createCmd.Flags().BoolP("sudo", "S", false, "添加sudo权限")
	createCmd.Flags().BoolP("force", "f", false, "强制创建（覆盖已存在的用户）")
}

// validateUsername 验证用户名格式
func validateUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("用户名不能为空")
	}
	if len(username) > 32 {
		return fmt.Errorf("用户名长度不能超过32个字符")
	}
	// 用户名只能包含字母、数字、下划线和连字符
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_' || char == '-') {
			return fmt.Errorf("用户名只能包含字母、数字、下划线和连字符")
		}
	}
	// 用户名不能以数字开头
	if username[0] >= '0' && username[0] <= '9' {
		return fmt.Errorf("用户名不能以数字开头")
	}
	return nil
}

// promptPassword 提示用户输入密码
func promptPassword() (string, error) {
	color.Info.Print("请输入密码: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	color.Print("\n")

	password := string(passwordBytes)
	if len(password) == 0 {
		return "", fmt.Errorf("密码不能为空")
	}

	// 确认密码
	color.Info.Print("请确认密码: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	color.Print("\n")

	confirm := string(confirmBytes)
	if password != confirm {
		return "", fmt.Errorf("两次输入的密码不一致")
	}

	return password, nil
}

// setUserShell 设置用户Shell
func setUserShell(username, shell string) error {
	err, _ := user.RunShell("chsh", "-s", shell, username)
	return err
}

// addUserToGroups 将用户添加到指定用户组
func addUserToGroups(username string, groups []string) error {
	for _, group := range groups {
		err, _ := user.RunShell("usermod", "-a", "-G", group, username)
		if err != nil {
			return fmt.Errorf("添加到用户组 %s 失败: %v", group, err)
		}
	}
	return nil
}

// addUserToSudo 添加用户到sudo组
func addUserToSudo(username string) error {
	// 尝试添加到sudo组
	err, _ := user.RunShell("usermod", "-a", "-G", "sudo", username)
	if err != nil {
		// 如果sudo组不存在，尝试wheel组（CentOS/RHEL）
		err, _ = user.RunShell("usermod", "-a", "-G", "wheel", username)
		if err != nil {
			return fmt.Errorf("添加sudo权限失败，请手动添加用户到sudo或wheel组")
		}
	}
	return nil
}
