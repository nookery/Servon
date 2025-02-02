package version

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Version 将在构建时通过 ldflags 注入
var Version string

type PackageJSON struct {
	Version string `json:"version"`
}

// GetVersion 返回应用程序版本
func GetVersion() string {
	if Version == "" {
		// 尝试从 package.json 读取版本
		if ver := getPackageVersion(); ver != "" {
			return ver + " (development)"
		}
		return "development"
	}
	return Version
}

func getPackageVersion() string {
	// 尝试读取项目根目录的 package.json
	data, err := os.ReadFile(filepath.Join("package.json"))
	if err != nil {
		return ""
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}

	return pkg.Version
}
