package state_machine

import (
	"context"
	"fmt"
	"time"
)

type TaskPool struct {
	tasks chan *RetryTask
}

type RetryTask struct {
	ctx  context.Context
	node *Node
}

func NewTaskPool(size int) *TaskPool {
	return &TaskPool{tasks: make(chan *RetryTask, size)}
}

func (p *TaskPool) AddNode(ctx context.Context, node *Node) {
	p.tasks <- &RetryTask{ctx: ctx, node: node}
}

func (p *TaskPool) Start() {
	for {
		task := <-p.tasks
		go func() {
			defer func() {
				if it := recover(); it != nil {
					fmt.Printf("[AlarmPanic] panic %v", it)
				}
			}()
			_ = NewMainProcessor(task.node).Process(task.ctx)
		}()

		time.Sleep(10 * time.Millisecond)
	}
}
