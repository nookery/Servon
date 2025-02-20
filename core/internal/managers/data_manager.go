package managers

import "os"

var DefaultDataManager = newDataManager()

type DataManager struct{}

func newDataManager() *DataManager {
	return &DataManager{}
}

// GetDataRootFolder 获取数据根目录，如果不存在会自动创建
func (c *DataManager) GetDataRootFolder() string {
	folder := "/data"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	return folder
}

// GetProjectsRootFolder 获取项目根目录，如果不存在会自动创建
func (c *DataManager) GetProjectsRootFolder() string {
	folder := c.GetDataRootFolder() + "/projects"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	return folder
}

// GetSoftwaresRootFolder 获取软件根目录，如果不存在会自动创建
func (c *DataManager) GetSoftwaresRootFolder() string {
	folder := c.GetDataRootFolder() + "/softwares"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	return folder
}

// GetSoftwareRootFolder 获取软件根目录
func (c *DataManager) GetSoftwareRootFolder(softwareName string) string {
	folder := c.GetSoftwaresRootFolder() + "/" + softwareName

	return folder
}

// GetAndCreateSoftwareRootFolder 获取软件根目录，如果不存在会自动创建
func (c *DataManager) GetAndCreateSoftwareRootFolder(softwareName string) string {
	folder := c.GetSoftwareRootFolder(softwareName)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	return folder
}

// GetConfigRootFolder 获取配置根目录
func (c *DataManager) GetConfigRootFolder() string {
	folder := c.GetDataRootFolder() + "/config"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, 0755)
	}

	return folder
}
