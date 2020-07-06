package main

import "sync"

//Task what to do
type Task interface {
	Do()
}

//TaskPool goroutine pool
type TaskPool struct {
	task chan Task
	wg   sync.WaitGroup
}

//NewTaskPool task pool constructor
func NewTaskPool(size int) *TaskPool {
	tp := &TaskPool{}
	tp.task = make(chan Task)
	tp.wg.Add(size)

	for i := 0; i < size; i++ {
		go func() {
			defer tp.wg.Done()
			for t := range tp.task {
				t.Do()
			}
		}()
	}

	return tp
}

//Add a task
func (tp *TaskPool) Add(t Task) {
	go func() {
		tp.task <- t
	}()
}

//Close task pool
func (tp *TaskPool) Close() {
	close(tp.task)
	tp.wg.Wait()
}
