package astro

import (
	"os"
	"strings"
	"time"
)

// isAstroProject 判断是否是 Astro 项目
func isAstroProject(projectFolder string) bool {
	if _, err := os.Stat(projectFolder + "/astro.config.mjs"); os.IsNotExist(err) {
		return false
	}

	return true
}

// getProjectNameFromRepo 从仓库地址中获取项目名称
// 比如：https://github.com/user/project.git 返回 project
// 比如：git@github.com:user/project.git 返回 project
// 比如：ssh://git@github.com/user/project.git 返回 project
// 比如：git+ssh://git@github.com/user/project.git 返回 project
// 比如：git+https://github.com/user/project.git 返回 project
// 比如：git+http://github.com/user/project.git 返回 project
// 如果不能获取到项目名称，则返回随机字符串（根据当前时间生成）
func getProjectNameFromRepo(repo string) string {
	repo = strings.TrimSuffix(repo, ".git")
	repo = strings.TrimPrefix(repo, "https://")
	repo = strings.TrimPrefix(repo, "http://")
	repo = strings.TrimPrefix(repo, "git@")
	repo = strings.TrimPrefix(repo, "ssh://")
	repo = strings.TrimPrefix(repo, "git+")
	repo = strings.TrimPrefix(repo, "git+ssh://")

	parts := strings.Split(repo, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return time.Now().Format("20060102150405")
}
