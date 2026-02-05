package main

import (
	"fmt"
	"sync"
)

type Task struct {
	ID      int
	Data    interface{}
	Process func(interface{}) (interface{}, error)
}

type Result struct {
	TaskID int
	Value  interface{}
	Error  error
}

type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	results    chan Result
	done       chan struct{}
	wg         sync.WaitGroup
}

func NewWorkerPool(n int) *WorkerPool {
	return &WorkerPool{
		numWorkers: n,
		tasks:      make(chan Task, n),
		results:    make(chan Result, n),
		done:       make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			wp.results <- Result{TaskID: -1, Error: fmt.Errorf("panic: %v", r)}
		}
	}()
	for task := range wp.tasks {
		val, err := task.Process(task.Data)
		wp.results <- Result{TaskID: task.ID, Value: val, Error: err}
	}
}

func main() {
	// TODO: Implementare il worker pool
	fmt.Println("Worker Pool Pattern")
}
