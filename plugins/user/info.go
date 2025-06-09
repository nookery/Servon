package user

import (
	"fmt"
	"os"
	"servon/components/user"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// infoCmd 显示用户详细信息
var infoCmd = &cobra.Command{
	Use:   "info [username]",
	Short: "显示用户详细信息",
	Long:  color.Success.Render("\r\n显示指定用户的详细信息，包括基本信息、权限、登录历史等"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		showProcesses, _ := cmd.Flags().GetBool("processes")
		showLoginHistory, _ := cmd.Flags().GetBool("login-history")

		userManager := user.NewUserManager()

		// 检查用户是否存在
		exists, err := userManager.UserExists(username)
		if err != nil {
			return fmt.Errorf("检查用户是否存在失败: %v", err)
		}
		if !exists {
			return fmt.Errorf("用户 %s 不存在", username)
		}

		// 获取用户列表并找到目标用户
		users, err := userManager.GetUserList()
		if err != nil {
			return fmt.Errorf("获取用户信息失败: %v", err)
		}

		var targetUser *user.User
		for _, u := range users {
			if u.Username == username {
				targetUser = &u
				break
			}
		}

		if targetUser == nil {
			return fmt.Errorf("未找到用户 %s 的详细信息", username)
		}

		// 显示用户基本信息
		displayUserBasicInfo(targetUser)

		// 显示用户权限信息
		displayUserPermissions(targetUser)

		// 显示用户统计信息
		displayUserStats(targetUser)

		// 显示用户进程（如果请求）
		if showProcesses {
			displayUserProcesses(username)
		}

		// 显示登录历史（如果请求）
		if showLoginHistory {
			displayLoginHistory(username)
		}

		return nil
	},
}

func init() {
	// 添加命令行参数
	infoCmd.Flags().BoolP("processes", "p", false, "显示用户当前运行的进程")
	infoCmd.Flags().BoolP("login-history", "l", false, "显示用户登录历史")
}

// displayUserBasicInfo 显示用户基本信息
func displayUserBasicInfo(u *user.User) {
	color.Primary.Printf("👤 用户信息: %s\n", u.Username)
	color.Print("\n")

	// 基本信息
	color.Info.Print("📋 基本信息:\n")
	color.Gray.Printf("   用户名: %s\n", u.Username)
	color.Gray.Printf("   主目录: %s\n", u.HomeDir)
	color.Gray.Printf("   Shell: %s\n", u.Shell)
	color.Gray.Printf("   用户组: %s\n", strings.Join(u.Groups, ", "))

	// 获取UID和GID
	uid, gid := getUserIDs(u.Username)
	if uid != "" {
		color.Gray.Printf("   UID: %s\n", uid)
	}
	if gid != "" {
		color.Gray.Printf("   GID: %s\n", gid)
	}

	// 时间信息
	if !u.CreateTime.IsZero() {
		color.Gray.Printf("   创建时间: %s\n", u.CreateTime.Format("2006-01-02 15:04:05"))
	}
	if !u.LastLogin.IsZero() {
		color.Gray.Printf("   最后登录: %s\n", u.LastLogin.Format("2006-01-02 15:04:05"))
	}

	color.Print("\n")
}

// displayUserPermissions 显示用户权限信息
func displayUserPermissions(u *user.User) {
	color.Info.Print("🔐 权限信息:\n")

	// Sudo权限
	if u.Sudo {
		color.Success.Print("   ✅ 拥有sudo权限\n")
	} else {
		color.Gray.Print("   ❌ 无sudo权限\n")
	}

	// 检查主目录权限
	homeDirPerms := getHomeDirPermissions(u.HomeDir)
	if homeDirPerms != "" {
		color.Gray.Printf("   主目录权限: %s\n", homeDirPerms)
	}

	// 检查Shell是否有效
	shellValid := isValidShell(u.Shell)
	if shellValid {
		color.Success.Printf("   ✅ Shell有效: %s\n", u.Shell)
	} else {
		color.Warn.Printf("   ⚠️  Shell可能无效: %s\n", u.Shell)
	}

	color.Print("\n")
}

// displayUserStats 显示用户统计信息
func displayUserStats(u *user.User) {
	color.Info.Print("📊 统计信息:\n")

	// 主目录大小
	homeDirSize := getHomeDirSize(u.HomeDir)
	if homeDirSize != "" {
		color.Gray.Printf("   主目录大小: %s\n", homeDirSize)
	}

	// 文件数量
	fileCount := getHomeFileCount(u.HomeDir)
	if fileCount != "" {
		color.Gray.Printf("   主目录文件数: %s\n", fileCount)
	}

	color.Print("\n")
}

// displayUserProcesses 显示用户当前运行的进程
func displayUserProcesses(username string) {
	color.Info.Printf("🔄 用户 %s 的运行进程:\n", username)

	err, output := user.RunShell("ps", "-u", username, "-o", "pid,ppid,pcpu,pmem,time,comm")
	if err != nil {
		color.Warn.Printf("   获取进程信息失败: %v\n", err)
		return
	}

	if strings.TrimSpace(output) == "" {
		color.Gray.Print("   无运行进程\n")
	} else {
		color.Gray.Print("   ")
		color.Gray.Print(strings.ReplaceAll(output, "\n", "\n   "))
		color.Print("\n")
	}

	color.Print("\n")
}

// displayLoginHistory 显示用户登录历史
func displayLoginHistory(username string) {
	color.Info.Printf("📅 用户 %s 的登录历史:\n", username)

	err, output := user.RunShell("last", "-n", "10", username)
	if err != nil {
		color.Warn.Printf("   获取登录历史失败: %v\n", err)
		return
	}

	if strings.TrimSpace(output) == "" {
		color.Gray.Print("   无登录记录\n")
	} else {
		color.Gray.Print("   ")
		color.Gray.Print(strings.ReplaceAll(output, "\n", "\n   "))
		color.Print("\n")
	}

	color.Print("\n")
}

// getUserIDs 获取用户的UID和GID
func getUserIDs(username string) (string, string) {
	err, output := user.RunShell("id", username)
	if err != nil {
		return "", ""
	}

	// 解析id命令输出: uid=1000(username) gid=1000(groupname) groups=...
	parts := strings.Fields(output)
	var uid, gid string

	for _, part := range parts {
		if strings.HasPrefix(part, "uid=") {
			uidPart := strings.Split(part, "=")[1]
			uid = strings.Split(uidPart, "(")[0]
		} else if strings.HasPrefix(part, "gid=") {
			gidPart := strings.Split(part, "=")[1]
			gid = strings.Split(gidPart, "(")[0]
		}
	}

	return uid, gid
}

// getHomeDirPermissions 获取主目录权限
func getHomeDirPermissions(homeDir string) string {
	info, err := os.Stat(homeDir)
	if err != nil {
		return ""
	}
	return info.Mode().String()
}

// isValidShell 检查Shell是否有效
func isValidShell(shell string) bool {
	_, err := os.Stat(shell)
	return err == nil
}

// getHomeDirSize 获取主目录大小
func getHomeDirSize(homeDir string) string {
	err, output := user.RunShell("du", "-sh", homeDir)
	if err != nil {
		return ""
	}
	parts := strings.Fields(output)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// getHomeFileCount 获取主目录文件数量
func getHomeFileCount(homeDir string) string {
	// 限制搜索深度和时间，避免长时间运行
	err, output := user.RunShell("sh", "-c", fmt.Sprintf("timeout 5 find %s -maxdepth 2 -type f 2>/dev/null | wc -l", homeDir))
	if err != nil {
		// 如果timeout命令不可用，尝试简单计数
		err, output = user.RunShell("sh", "-c", fmt.Sprintf("ls -la %s 2>/dev/null | grep '^-' | wc -l", homeDir))
		if err != nil {
			return "未知"
		}
	}
	return strings.TrimSpace(output)
}
