package providers

import (
	"servon/core/internal/utils"
)

type UtilProvider struct {
	*utils.CommandUtil
	*utils.FileUtil
	*utils.DevUtil
	*utils.StringUtil
}

func NewUtilProvider() *UtilProvider {
	return &UtilProvider{
		CommandUtil: utils.DefaultCommandUtil,
		FileUtil:    utils.DefaultFileUtil,
		DevUtil:     utils.DefaultDevUtil,
		StringUtil:  utils.DefaultStringUtil,
	}
}
