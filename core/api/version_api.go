package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servon/core/libs"
	"strings"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version = "0.1.0"
)

type VersionApi struct {
	version string
}

func NewVersion() VersionApi {
	return VersionApi{
		version: Version,
	}
}

// GetVersionCommand 返回版本命令
func (c *VersionApi) GetVersionCommand() *cobra.Command {
	return libs.NewCommand(libs.CommandOptions{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(c.GetVersion())
		},
	})
}

func (c *VersionApi) GetVersion() string {
	return c.version
}

func (c *VersionApi) GetLatestVersion() (string, error) {
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
