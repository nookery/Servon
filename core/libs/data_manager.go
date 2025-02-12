package libs

type DataManager struct{}

func NewDataManager() *DataManager {
	return &DataManager{}
}

func (c *DataManager) GetDataRootFolder() string {
	return "/data"
}
