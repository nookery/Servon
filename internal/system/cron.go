package system

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

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

// GetCronTasks 获取所有定时任务
func GetCronTasks() ([]*CronTask, error) {
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

// CreateCronTask 创建定时任务
func CreateCronTask(task CronTask) (*CronTask, error) {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	lastID++
	task.ID = lastID
	task.Enabled = true

	// 验证cron表达式
	if _, err := cron.ParseStandard(task.Schedule); err != nil {
		return nil, fmt.Errorf("无效的定时表达式: %v", err)
	}

	// 添加到cron调度器
	entryID, err := cronInstance.AddFunc(task.Schedule, func() {
		executeTask(&task)
	})
	if err != nil {
		return nil, err
	}

	task.entryID = entryID
	tasks[task.ID] = &task

	return &task, nil
}

// UpdateCronTask 更新定时任务
func UpdateCronTask(task CronTask) (*CronTask, error) {
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

	// 验证新的cron表达式
	if _, err := cron.ParseStandard(task.Schedule); err != nil {
		return nil, fmt.Errorf("无效的定时表达式: %v", err)
	}

	// 如果任务需要启用，添加新的调度
	if task.Enabled {
		entryID, err := cronInstance.AddFunc(task.Schedule, func() {
			executeTask(&task)
		})
		if err != nil {
			return nil, err
		}
		task.entryID = entryID
	}

	tasks[task.ID] = &task
	return &task, nil
}

// DeleteCronTask 删除定时任务
func DeleteCronTask(id int) error {
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
func ToggleCronTask(id int) (*CronTask, error) {
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
