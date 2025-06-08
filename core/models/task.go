package models

type Task struct {
	ID      string       `json:"id"`
	Execute func() error // 可执行的函数
}
