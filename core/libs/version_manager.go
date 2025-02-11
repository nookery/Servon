package libs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version = "0.1.0"
)

type VersionManager struct {
	Version string
}

func NewVersionManager() *VersionManager {
	return &VersionManager{}
}

type VersionApi struct {
	version string
}

func NewVersion() VersionApi {
	return VersionApi{
		version: Version,
	}
}

// GetVersionCommand 返回版本命令
func (c *VersionManager) GetVersionCommand() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			printer.Printf("Version: %s\n", c.GetVersion())
		},
	})
}

func (c *VersionManager) GetVersion() string {
	return c.Version
}

func (c *VersionManager) GetLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/nookery/servon/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
}
