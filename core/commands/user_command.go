package commands

import (
	"fmt"
	"servon/core/managers"
	"strings"

	"github.com/spf13/cobra"
)

// GetUserRootCommand 获取用户管理命令
func GetUserRootCommand(u *managers.UserManager) *cobra.Command {
	rootCmd := NewCommand(CommandOptions{
		Use:   "user",
		Short: "用户管理",
	})

	rootCmd.AddCommand(GetUserListCommand(u))
	rootCmd.AddCommand(CreateUserCommand(u))
	rootCmd.AddCommand(DeleteUserCommand(u))

	return rootCmd
}

// GetUserListCommand 获取用户列表命令
func GetUserListCommand(u *managers.UserManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "list",
		Short: "获取用户列表",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := u.GetUserList()
			if err != nil {
				logger.Error(err)
				return
			}

			// Convert []User to []string
			userStrings := make([]string, len(users))
			for i, user := range users {
				userStrings[i] = fmt.Sprintf("%s (%s)", user.Username, strings.Join(user.Groups, ","))
			}

			logger.ListWithTitle("用户列表", userStrings)
		},
	})
}

// CreateUserCommand 创建用户命令
func CreateUserCommand(u *managers.UserManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "create",
		Short: "创建用户",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				logger.ErrorMessage("请提供用户名和密码，例如：user create username password")
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
func DeleteUserCommand(u *managers.UserManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "delete",
		Short:   "删除用户",
		Aliases: []string{"del", "d"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logger.ErrorMessage("请提供用户名")
				return
			}
			username := args[0]
			err := u.DeleteUser(username)
			if err != nil {
				PrintError(err)
				return
			}

			PrintSuccessf("用户 %s 已删除", username)
		},
	})
}
