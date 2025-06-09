package user

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type UserManager struct{}

func NewUserManager() *UserManager {
	return &UserManager{}
}

// User 表示系统用户的结构体
type User struct {
	Username   string    `json:"username"`
	Groups     []string  `json:"groups"`
	Shell      string    `json:"shell"`
	HomeDir    string    `json:"home_dir"`
	CreateTime time.Time `json:"create_time"`
	LastLogin  time.Time `json:"last_login"`
	Sudo       bool      `json:"sudo"`
}

// GetUserList 获取系统用户列表
func (u *UserManager) GetUserList() ([]User, error) {
	// 检测操作系统类型
	if u.isMacOS() {
		return u.getUserListMacOS()
	} else {
		return u.getUserListLinux()
	}
}

// isMacOS 检测是否为macOS系统
func (u *UserManager) isMacOS() bool {
	_, err := os.Stat("/System/Library/CoreServices/SystemVersion.plist")
	return err == nil
}

// getUserListMacOS 获取macOS系统用户列表
func (u *UserManager) getUserListMacOS() ([]User, error) {
	// 使用dscl命令获取用户列表
	err, output := RunShellWithOutput("dscl", ".", "list", "/Users")
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}

	var users []User
	usernames := strings.Split(strings.TrimSpace(output), "\n")

	for _, username := range usernames {
		if username == "" {
			continue
		}

		// 获取用户详细信息
		userInfo, err := u.getUserInfoMacOS(username)
		if err != nil {
			continue // 跳过无法获取信息的用户
		}

		users = append(users, *userInfo)
	}

	return users, nil
}

// getUserInfoMacOS 获取macOS用户详细信息
func (u *UserManager) getUserInfoMacOS(username string) (*User, error) {
	// 获取用户主目录
	err, homeDir := RunShellWithOutput("dscl", ".", "read", "/Users/"+username, "NFSHomeDirectory")
	if err != nil {
		return nil, err
	}
	homeDir = strings.TrimPrefix(strings.TrimSpace(homeDir), "NFSHomeDirectory: ")

	// 获取用户Shell
	err, shell := RunShellWithOutput("dscl", ".", "read", "/Users/"+username, "UserShell")
	if err != nil {
		shell = "/bin/bash" // 默认shell
	} else {
		shell = strings.TrimPrefix(strings.TrimSpace(shell), "UserShell: ")
	}

	// 获取用户组信息
	groups, _ := u.getUserGroups(username)

	// 获取用户创建时间（通过 home 目录创建时间估算）
	createTime := u.getUserCreateTime(homeDir)

	// 获取最后登录时间
	lastLogin := u.getLastLogin(username)

	// 检查是否有 sudo 权限
	sudo := u.hasSudoPermission(username)

	return &User{
		Username:   username,
		Groups:     groups,
		Shell:      shell,
		HomeDir:    homeDir,
		CreateTime: createTime,
		LastLogin:  lastLogin,
		Sudo:       sudo,
	}, nil
}

// getUserListLinux 获取Linux系统用户列表
func (u *UserManager) getUserListLinux() ([]User, error) {
	// 使用 os/user 包读取 /etc/passwd
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) >= 7 {
			// 获取用户信息
			username := fields[0]
			homeDir := fields[5]
			shell := fields[6]

			// 获取用户组信息
			groups, _ := u.getUserGroups(username)

			// 获取用户创建时间（通过 home 目录创建时间估算）
			createTime := u.getUserCreateTime(homeDir)

			// 获取最后登录时间
			lastLogin := u.getLastLogin(username)

			// 检查是否有 sudo 权限
			sudo := u.hasSudoPermission(username)

			users = append(users, User{
				Username:   username,
				Groups:     groups,
				Shell:      shell,
				HomeDir:    homeDir,
				CreateTime: createTime,
				LastLogin:  lastLogin,
				Sudo:       sudo,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// 获取用户组信息
func (u *UserManager) getUserGroups(username string) ([]string, error) {
	err, output := RunShellWithOutput("groups", username)
	if err != nil {
		return nil, err
	}
	// 解析输出格式 "username : group1 group2 group3"
	parts := strings.Split(output, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected groups output format")
	}
	groups := strings.Fields(parts[1])
	return groups, nil
}

// 获取用户创建时间
func (u *UserManager) getUserCreateTime(homeDir string) time.Time {
	info, err := os.Stat(homeDir)
	if err != nil {
		return time.Time{} // 返回零值表示未知
	}
	return info.ModTime()
}

// 获取最后登录时间
func (u *UserManager) getLastLogin(username string) time.Time {
	err, _ := RunShell("last", "-1", username)
	if err != nil {
		return time.Time{} // 返回零值表示未知
	}
	// TODO: 解析 last 命令输出获取最后登录时间
	return time.Now() // 临时返回当前时间
}

// 检查是否有 sudo 权限
func (u *UserManager) hasSudoPermission(username string) bool {
	// 检查用户是否在 sudo 组中
	err, output := RunShellWithOutput("groups", username)
	if err != nil {
		return false
	}
	return strings.Contains(output, "sudo") || strings.Contains(output, "wheel")
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
	err, _ = RunShell("useradd", "-m", username)
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

	err, _ = RunShell("userdel", "-r", username)
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

// UserExists 检查用户是否存在
func (u *UserManager) UserExists(username string) (bool, error) {
	err, output := RunShellWithOutput("id", username)

	if err != nil {
		if strings.Contains(output, "no such user") {
			return false, nil
		}

		return false, PrintAndReturnErrorf("检查用户是否存在失败: %v", err)
	}

	return true, nil
}
