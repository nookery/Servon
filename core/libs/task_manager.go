package libs

import (
	"fmt"
)

type Task struct {
	ID      string       `json:"id"`
	Execute func() error // 可执行的函数
}

type TaskManager struct {
	tasks map[string]Task // 存储任务的映射
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]Task),
	}
}

// GetTasks 获取所有任务
func (tm *TaskManager) GetTasks() []Task {
	tasks := make([]Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// AddTask 添加任务
func (tm *TaskManager) AddTask(task Task) {
	PrintInfof("添加任务: %s", task.ID)
	if tm.tasks == nil {
		tm.tasks = make(map[string]Task)
	}
	tm.tasks[task.ID] = task
}

// RemoveTask 删除任务
func (tm *TaskManager) RemoveTask(id string) {
	delete(tm.tasks, id)
}

// ExecuteTask 执行任务
func (tm *TaskManager) ExecuteTask(id string) error {
	task, exists := tm.tasks[id]
	if !exists {
		return fmt.Errorf("任务未找到: %s", id)
	}
	return task.Execute()
}
