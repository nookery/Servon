package api

type Data struct{}

func NewData() Data {
	return Data{}
}

func (c *Data) GetDataRootFolder() string {
	return "/data"
}
