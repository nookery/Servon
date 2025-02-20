package providers

import (
	"servon/core/internal/libs/utils"
)

type UtilProvider struct {
	*utils.Printer
	*utils.CommandUtil
	*utils.FileUtil
	*utils.DevUtil
	*utils.StringUtil
}

func NewUtilProvider() *UtilProvider {
	return &UtilProvider{
		Printer:     utils.DefaultPrinter,
		CommandUtil: utils.DefaultCommandUtil,
		FileUtil:    utils.DefaultFileUtil,
		DevUtil:     utils.DefaultDevUtil,
		StringUtil:  utils.DefaultStringUtil,
	}
}
