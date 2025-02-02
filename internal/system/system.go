package system

import (
	"fmt"
	"os/user"
)

// GetCurrentUser 获取当前系统用户名
func GetCurrentUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("获取当前用户失败: %v", err)
	}
	return currentUser.Username, nil
}
