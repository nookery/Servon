package libs

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type UserManager struct{}

func NewUserManager() *UserManager {
	return &UserManager{}
}

// GetUserList 获取系统用户列表
func (u *UserManager) GetUserList() ([]string, error) {
	osType := GetOSType()
	if osType != Ubuntu {
		return nil, fmt.Errorf("不支持的操作系统类型: %s", osType)
	}

	output, err := RunShellWithOutput("cat", "/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}

	var users []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) >= 1 {
			users = append(users, fields[0])
		}
	}
	return users, nil
}

// CreateUser 创建新用户
func (u *UserManager) CreateUser(username string, password string) error {
	exists, err := u.UserExists(username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("用户 %s 已存在", username)
	}

	// 创建用户
	err = RunShell("useradd", "-m", username)
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}

	// 设置密码
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", username, password))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置密码失败: %v", err)
	}

	return nil
}

// DeleteUser 删除用户
func (u *UserManager) DeleteUser(username string) error {
	exists, err := u.UserExists(username)
	if err != nil {
		return fmt.Errorf("检查用户是否存在失败: %v", err)
	}
	if !exists {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	err = RunShell("userdel", "-r", username)
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

// UserExists 检查用户是否存在
func (u *UserManager) UserExists(username string) (bool, error) {
	output, err := RunShellWithOutput("id", username)

	if err != nil {
		if strings.Contains(output, "no such user") {
			return false, nil
		}

		PrintError(err)

		return false, err
	}

	return true, nil
}

// --- 命令 ---

// GetUserRootCommand 获取用户管理命令
func (u *UserManager) GetUserRootCommand() *cobra.Command {
	rootCmd := NewCommand(CommandOptions{
		Use:   "user",
		Short: "用户管理",
	})

	rootCmd.AddCommand(u.GetUserListCommand())
	rootCmd.AddCommand(u.CreateUserCommand())
	rootCmd.AddCommand(u.DeleteUserCommand())

	return rootCmd
}

// GetUserListCommand 获取用户列表命令
func (u *UserManager) GetUserListCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "list",
		Short: "获取用户列表",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := u.GetUserList()
			if err != nil {
				DefaultPrinter.PrintError(err)

			}
			DefaultPrinter.PrintList(users, "用户列表")
		},
	})
}

// CreateUserCommand 创建用户命令
func (u *UserManager) CreateUserCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "create",
		Short: "创建用户",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				DefaultPrinter.PrintErrorMessage("请提供用户名和密码，例如：user create username password")
				return
			}
			username := args[0]
			password := args[1]
			err := u.CreateUser(username, password)
			if err != nil {
				PrintError(err)
			}
		},
	})
}

// DeleteUserCommand 删除用户命令
func (u *UserManager) DeleteUserCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "delete",
		Short:   "删除用户",
		Aliases: []string{"del", "d"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				DefaultPrinter.PrintErrorMessage("请提供用户名")
				return
			}
			username := args[0]
			err := u.DeleteUser(username)
			if err != nil {
				PrintError(err)
				return
			}

			PrintSuccess("用户 %s 已删除", username)
		},
	})
}
