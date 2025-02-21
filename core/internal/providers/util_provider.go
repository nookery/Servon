package providers

import (
	"servon/core/internal/utils"
)

type UtilProvider struct {
	*utils.LogUtil
	*utils.CommandUtil
	*utils.FileUtil
	*utils.DevUtil
	*utils.StringUtil
}

func NewUtilProvider() *UtilProvider {
	return &UtilProvider{
		LogUtil:     utils.DefaultLogUtil,
		CommandUtil: utils.DefaultCommandUtil,
		FileUtil:    utils.DefaultFileUtil,
		DevUtil:     utils.DefaultDevUtil,
		StringUtil:  utils.DefaultStringUtil,
	}
}
