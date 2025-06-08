package managers

import (
	"servon/components/cron_util"
)

var DefaultCronManager = newCronManager()

type CronManager struct {
	taskManager *cron_util.CronTaskManager
}

func newCronManager() *CronManager {
	return &CronManager{
		taskManager: cron_util.NewCronTaskManager(),
	}
}

// GetCronTasks 获取所有定时任务
func (p *CronManager) GetCronTasks() ([]*cron_util.CronTask, error) {
	PrintInfo("获取所有定时任务...")
	return p.taskManager.GetTasks(), nil
}

// CreateCronTask 创建定时任务
func (p *CronManager) CreateCronTask(task cron_util.CronTask) (*cron_util.CronTask, error) {
	PrintInfo("创建定时任务...")
	return p.taskManager.CreateTask(task)
}

// UpdateCronTask 更新定时任务
func (p *CronManager) UpdateCronTask(task cron_util.CronTask) (*cron_util.CronTask, error) {
	PrintInfo("更新定时任务...")
	return p.taskManager.UpdateTask(task)
}

// DeleteCronTask 删除定时任务
func (p *CronManager) DeleteCronTask(id int) error {
	PrintInfo("删除定时任务...")
	return p.taskManager.DeleteTask(id)
}

// ToggleCronTask 启用/禁用定时任务
func (p *CronManager) ToggleCronTask(id int) (*cron_util.CronTask, error) {
	PrintInfo("启用/禁用定时任务...")
	return p.taskManager.ToggleTask(id)
}
