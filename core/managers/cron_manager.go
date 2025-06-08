package managers

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var DefaultCronManager = newCronManager()

type CronTask struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Command     string    `json:"command"`
	Schedule    string    `json:"schedule"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	LastRun     time.Time `json:"last_run,omitempty"`
	NextRun     time.Time `json:"next_run,omitempty"`
	entryID     cron.EntryID
}

var (
	cronInstance *cron.Cron
	tasks        = make(map[int]*CronTask)
	tasksMutex   sync.RWMutex
	lastID       = 0
)

func init() {
	cronInstance = cron.New(cron.WithSeconds())
	cronInstance.Start()
}

type CronManager struct {
}

func newCronManager() *CronManager {
	return &CronManager{}
}

// GetCronTasks 获取所有定时任务
func (p *CronManager) GetCronTasks() ([]*CronTask, error) {
	PrintInfo("获取所有定时任务...")
	tasksMutex.RLock()
	defer tasksMutex.RUnlock()

	result := make([]*CronTask, 0, len(tasks))
	for _, task := range tasks {
		if task.Enabled {
			entry := cronInstance.Entry(task.entryID)
			task.NextRun = entry.Next
		}
		result = append(result, task)
	}
	return result, nil
}

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

// validateTask 验证任务的各个字段
func (p *CronManager) validateTask(task CronTask) error {
	var errors []ValidationError

	// 验证名称
	if task.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "任务名称不能为空",
		})
	}

	// 验证命令
	if task.Command == "" {
		errors = append(errors, ValidationError{
			Field:   "command",
			Message: "执行命令不能为空",
		})
	}

	// 验证cron表达式
	if task.Schedule == "" {
		errors = append(errors, ValidationError{
			Field:   "schedule",
			Message: "定时表达式不能为空",
		})
	} else {
		// 使用支持秒的解析器
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

// CreateCronTask 创建定时任务
func (p *CronManager) CreateCronTask(task CronTask) (*CronTask, error) {
	PrintInfo("创建定时任务...")
	// 先进行字段验证
	if err := p.validateTask(task); err != nil {
		return nil, err
	}

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	lastID++
	task.ID = lastID
	task.Enabled = true

	// 添加到cron调度器
	entryID, err := cronInstance.AddFunc(task.Schedule, func() {
		executeTask(&task)
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
	tasks[task.ID] = &task

	return &task, nil
}

// UpdateCronTask 更新定时任务
func (p *CronManager) UpdateCronTask(task CronTask) (*CronTask, error) {
	// 先进行字段验证
	if err := p.validateTask(task); err != nil {
		return nil, err
	}

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	existingTask, exists := tasks[task.ID]
	if !exists {
		return nil, fmt.Errorf("任务不存在")
	}

	// 如果任务已启用，先移除旧的调度
	if existingTask.Enabled {
		cronInstance.Remove(existingTask.entryID)
	}

	// 如果任务需要启用，添加新的调度
	if task.Enabled {
		entryID, err := cronInstance.AddFunc(task.Schedule, func() {
			executeTask(&task)
		})
		if err != nil {
			return nil, fmt.Errorf("添加任务失败: %v", err)
		}
		task.entryID = entryID
	}

	tasks[task.ID] = &task
	return &task, nil
}

// DeleteCronTask 删除定时任务
func (p *CronManager) DeleteCronTask(id int) error {
	PrintInfo("删除定时任务...")
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	task, exists := tasks[id]
	if !exists {
		return fmt.Errorf("任务不存在")
	}

	if task.Enabled {
		cronInstance.Remove(task.entryID)
	}

	delete(tasks, id)
	return nil
}

// ToggleCronTask 启用/禁用定时任务
func (p *CronManager) ToggleCronTask(id int) (*CronTask, error) {
	PrintInfo("启用/禁用定时任务...")
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	task, exists := tasks[id]
	if !exists {
		return nil, fmt.Errorf("任务不存在")
	}

	if task.Enabled {
		cronInstance.Remove(task.entryID)
		task.Enabled = false
	} else {
		entryID, err := cronInstance.AddFunc(task.Schedule, func() {
			executeTask(task)
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
func executeTask(task *CronTask) {
	tasksMutex.Lock()
	task.LastRun = time.Now()
	tasksMutex.Unlock()

	// 执行命令
	// TODO: 根据实际需求实现命令执行逻辑
	fmt.Printf("执行任务 %s: %s\n", task.Name, task.Command)
}
