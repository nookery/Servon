package providers

import (
	"servon/components/command_util"
	"servon/core/internal/utils"
)

type UtilProvider struct {
	*command_util.CommandUtil
	*utils.FileUtil
	*utils.DevUtil
	*utils.StringUtil
	*utils.ProjectUtil
}

func NewUtilProvider() *UtilProvider {
	return &UtilProvider{
		CommandUtil: command_util.DefaultCommandUtil,
		FileUtil:    utils.DefaultFileUtil,
		DevUtil:     utils.DefaultDevUtil,
		StringUtil:  utils.DefaultStringUtil,
		ProjectUtil: utils.DefaultProjectUtil,
	}
}
