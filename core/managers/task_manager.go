package managers

import (
	"fmt"
	"servon/core/models"
)

var DefaultTaskManager = newTaskManager()

type Task = models.Task

type TaskManager struct {
	tasks map[string]Task // 存储任务的映射
}

func newTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]models.Task),
	}
}

// GetTasks 获取所有任务
func (tm *TaskManager) GetTasks() []models.Task {
	tasks := make([]models.Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// AddTask 添加任务
func (tm *TaskManager) AddTask(task Task, reason string) {
	PrintInfof("添加任务: %s, 原因: %s", task.ID, reason)
	if tm.tasks == nil {
		tm.tasks = make(map[string]Task)
	}
	tm.tasks[task.ID] = task
}

// AddTaskAndExecute 添加任务并执行
func (tm *TaskManager) AddTaskAndExecute(task Task, reason string) {
	tm.AddTask(task, reason)
	tm.ExecuteTask(task.ID)
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
