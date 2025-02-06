package task

import (
	"atlas-toolkit/task/worker"
)

type Task struct {
	Worker *worker.Worker
}

// 一次性任务
func (tk Task) Once(options ...func(*worker.RunOptions)) error {
	return tk.Worker.Once(options...)
}

// 定时任务
func (tk Task) Cron(options ...func(*worker.RunOptions)) error {
	return tk.Worker.Cron(options...)
}
