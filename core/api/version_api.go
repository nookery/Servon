package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Version information
var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type VersionApi struct {
	version   string
	buildTime string
	gitCommit string
}

func NewVersion() VersionApi {
	return VersionApi{
		version:   Version,
		buildTime: BuildTime,
		gitCommit: GitCommit,
	}
}

func (c *VersionApi) GetVersion() string {
	return c.version
}

func (c *VersionApi) GetFullVersionInfo() string {
	return fmt.Sprintf("Version: %s\nBuild Time: %s\nGit Commit: %s",
		c.version,
		c.buildTime,
		c.gitCommit,
	)
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
