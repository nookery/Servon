package managers

import (
	"servon/core/internal/contract"
)

type ProjectManager struct {
	softManager *SoftManager
}

func NewTopologyManager(softManager *SoftManager) *ProjectManager {
	return &ProjectManager{
		softManager: softManager,
	}
}

func (m *ProjectManager) GetAllGateways() []string {
	return m.softManager.GetAllGateways()
}

func (m *ProjectManager) GetProjects(gatewayName string) ([]contract.Project, error) {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return nil, err
	}
	return gateway.GetProjects()
}

func (m *ProjectManager) AddProject(gatewayName string, project contract.Project) error {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.AddProject(project)
}

func (m *ProjectManager) RemoveProject(gatewayName string, projectName string) error {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.RemoveProject(projectName)
}

// GetGateway 获取指定的网关软件
func (m *ProjectManager) GetGateway(gatewayName string) (contract.SuperGateway, error) {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return nil, err
	}
	return gateway, nil
}
