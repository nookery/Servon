package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"servon/cmd/internal/softwares"
	"servon/cmd/internal/utils"
	"sync"
)

var (
	projects   = make(map[int]*Project)
	projectsMu sync.RWMutex
	lastID     = 0
	dataFile   = "data/projects.json"
)

// 加载项目数据
func loadProjects() error {
	utils.Debug("开始加载项目数据")
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Debug("项目数据文件不存在，将创建新文件")
			return nil
		}
		return err
	}

	var projectsList []*Project
	if err := json.Unmarshal(data, &projectsList); err != nil {
		return err
	}

	projectsMu.Lock()
	defer projectsMu.Unlock()

	for _, p := range projectsList {
		projects[p.ID] = p
		if p.ID > lastID {
			lastID = p.ID
		}
	}

	utils.Info("成功加载 %d 个项目", len(projectsList))
	return nil
}

// 保存项目数据
func saveProjects() error {
	projectsMu.RLock()
	projectsList := make([]*Project, 0, len(projects))
	for _, p := range projects {
		projectsList = append(projectsList, p)
	}
	projectsMu.RUnlock()

	data, err := json.MarshalIndent(projectsList, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, data, 0644)
}

// GetProjects 获取所有项目
func GetProjects() ([]*Project, error) {
	projectsMu.RLock()
	defer projectsMu.RUnlock()

	result := make([]*Project, 0, len(projects))
	for _, p := range projects {
		result = append(result, p)
	}
	return result, nil
}

// CreateProject 创建项目
func CreateProject(project Project) (*Project, error) {
	utils.Info("创建新项目: %s", project.Name)

	if err := validateProject(project); err != nil {
		utils.Error("项目验证失败: %v", err)
		return nil, err
	}

	projectsMu.Lock()
	defer projectsMu.Unlock()

	lastID++
	project.ID = lastID
	project.Status = "stopped"

	projects[project.ID] = &project
	if err := saveProjects(); err != nil {
		utils.Error("保存项目数据失败: %v", err)
		return nil, err
	}

	utils.Info("项目创建成功: [%d] %s", project.ID, project.Name)
	return &project, nil
}

// UpdateProject 更新项目
func UpdateProject(project Project) (*Project, error) {
	utils.Info("更新项目: [%d] %s", project.ID, project.Name)

	if err := validateProject(project); err != nil {
		utils.Error("项目验证失败: %v", err)
		return nil, err
	}

	projectsMu.Lock()
	defer projectsMu.Unlock()

	existing, exists := projects[project.ID]
	if !exists {
		utils.Error("项目不存在: %d", project.ID)
		return nil, fmt.Errorf("项目不存在")
	}

	// 保持原有状态和最后部署时间
	project.Status = existing.Status
	project.LastDeploy = existing.LastDeploy

	projects[project.ID] = &project
	if err := saveProjects(); err != nil {
		utils.Error("保存项目数据失败: %v", err)
		return nil, err
	}

	utils.Info("项目更新成功: [%d] %s", project.ID, project.Name)
	return &project, nil
}

// DeleteProject 删除项目
func DeleteProject(id int) error {
	utils.Info("删除项目: %d", id)

	projectsMu.Lock()
	defer projectsMu.Unlock()

	if _, exists := projects[id]; !exists {
		utils.Error("项目不存在: %d", id)
		return fmt.Errorf("项目不存在")
	}

	// 删除项目文件
	if err := cleanupProject(id); err != nil {
		utils.Error("清理项目文件失败: %v", err)
		return err
	}

	delete(projects, id)
	if err := saveProjects(); err != nil {
		utils.Error("保存项目数据失败: %v", err)
		return err
	}

	utils.Info("项目删除成功: %d", id)
	return nil
}

// validateProject 验证项目配置
func validateProject(project Project) error {
	if project.Name == "" {
		return fmt.Errorf("项目名称不能为空")
	}
	if project.GitRepo == "" {
		return fmt.Errorf("Git仓库地址不能为空")
	}
	if project.Branch == "" {
		return fmt.Errorf("分支不能为空")
	}
	if project.Domain == "" {
		return fmt.Errorf("域名不能为空")
	}
	if project.Type != "static" && project.Port == 0 {
		return fmt.Errorf("非静态项目必须指定端口")
	}
	return nil
}

// cleanupProject 清理项目文件
func cleanupProject(id int) error {
	projectDir := filepath.Join("data", "projects", fmt.Sprintf("%d", id))
	return os.RemoveAll(projectDir)
}

// ServeStatic 创建静态文件服务
func ServeStatic(name string, path string, domain string) error {
	utils.Info("创建静态文件服务: %s -> %s", path, domain)

	// 验证路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", path)
	}

	// 配置 Caddy
	caddy := softwares.NewCaddy()
	err := caddy.UpdateConfig(&softwares.Project{
		Name:      name,
		Domain:    domain,
		Type:      "static",
		OutputDir: path, // 直接使用提供的路径
	})
	if err != nil {
		return fmt.Errorf("配置 Caddy 失败: %v", err)
	}

	utils.Info("静态文件服务创建成功: %s", domain)
	return nil
}
