package cron_util

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// ValidationError 表示字段验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors 表示多个字段的验证错误
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return ""
	}
	return ve.Errors[0].Message
}

// CronTaskManager 管理定时任务的核心结构
type CronTaskManager struct {
	cronInstance *cron.Cron
	tasks        map[int]*CronTask
	tasksMutex   sync.RWMutex
	lastID       int
}

// NewCronTaskManager 创建一个新的CronTaskManager实例
func NewCronTaskManager() *CronTaskManager {
	manager := &CronTaskManager{
		cronInstance: cron.New(cron.WithSeconds()),
		tasks:        make(map[int]*CronTask),
		lastID:       0,
	}
	manager.cronInstance.Start()
	return manager
}

// GetTasks 获取所有定时任务
func (m *CronTaskManager) GetTasks() []*CronTask {
	m.tasksMutex.RLock()
	defer m.tasksMutex.RUnlock()

	result := make([]*CronTask, 0, len(m.tasks))
	for _, task := range m.tasks {
		if task.Enabled {
			entry := m.cronInstance.Entry(task.entryID)
			task.NextRun = entry.Next
		}
		result = append(result, task)
	}
	return result
}

// validateTask 验证任务的各个字段
func (m *CronTaskManager) validateTask(task CronTask) error {
	var errors []ValidationError

	if task.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "任务名称不能为空",
		})
	}

	if task.Command == "" {
		errors = append(errors, ValidationError{
			Field:   "command",
			Message: "执行命令不能为空",
		})
	}

	if task.Schedule == "" {
		errors = append(errors, ValidationError{
			Field:   "schedule",
			Message: "定时表达式不能为空",
		})
	} else {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err := parser.Parse(task.Schedule); err != nil {
			errors = append(errors, ValidationError{
				Field:   "schedule",
				Message: "无效的定时表达式: " + err.Error(),
			})
		}
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}
	return nil
}

// CreateTask 创建定时任务
func (m *CronTaskManager) CreateTask(task CronTask) (*CronTask, error) {
	if err := m.validateTask(task); err != nil {
		return nil, err
	}

	m.tasksMutex.Lock()
	defer m.tasksMutex.Unlock()

	m.lastID++
	task.ID = m.lastID
	task.Enabled = true

	entryID, err := m.cronInstance.AddFunc(task.Schedule, func() {
		m.executeTask(&task)
	})
	if err != nil {
		return nil, ValidationErrors{
			Errors: []ValidationError{
				{
					Field:   "schedule",
					Message: fmt.Sprintf("添加任务失败: %v", err),
				},
			},
		}
	}

	task.entryID = entryID
	m.tasks[task.ID] = &task

	return &task, nil
}

// UpdateTask 更新定时任务
func (m *CronTaskManager) UpdateTask(task CronTask) (*CronTask, error) {
	if err := m.validateTask(task); err != nil {
		return nil, err
	}

	m.tasksMutex.Lock()
	defer m.tasksMutex.Unlock()

	existingTask, exists := m.tasks[task.ID]
	if !exists {
		return nil, fmt.Errorf("任务不存在")
	}

	if existingTask.Enabled {
		m.cronInstance.Remove(existingTask.entryID)
	}

	if task.Enabled {
		entryID, err := m.cronInstance.AddFunc(task.Schedule, func() {
			m.executeTask(&task)
		})
		if err != nil {
			return nil, fmt.Errorf("添加任务失败: %v", err)
		}
		task.entryID = entryID
	}

	m.tasks[task.ID] = &task
	return &task, nil
}

// DeleteTask 删除定时任务
func (m *CronTaskManager) DeleteTask(id int) error {
	m.tasksMutex.Lock()
	defer m.tasksMutex.Unlock()

	task, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("任务不存在")
	}

	if task.Enabled {
		m.cronInstance.Remove(task.entryID)
	}

	delete(m.tasks, id)
	return nil
}

// ToggleTask 启用/禁用定时任务
func (m *CronTaskManager) ToggleTask(id int) (*CronTask, error) {
	m.tasksMutex.Lock()
	defer m.tasksMutex.Unlock()

	task, exists := m.tasks[id]
	if !exists {
		return nil, fmt.Errorf("任务不存在")
	}

	if task.Enabled {
		m.cronInstance.Remove(task.entryID)
		task.Enabled = false
	} else {
		entryID, err := m.cronInstance.AddFunc(task.Schedule, func() {
			m.executeTask(task)
		})
		if err != nil {
			return nil, err
		}
		task.entryID = entryID
		task.Enabled = true
	}

	return task, nil
}

// executeTask 执行定时任务
func (m *CronTaskManager) executeTask(task *CronTask) {
	m.tasksMutex.Lock()
	task.LastRun = time.Now()
	m.tasksMutex.Unlock()

	// 执行命令
	// TODO: 根据实际需求实现命令执行逻辑
	fmt.Printf("执行任务 %s: %s\n", task.Name, task.Command)
}
