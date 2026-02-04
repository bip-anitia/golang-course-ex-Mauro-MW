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

func main() {
	// TODO: Implementare il worker pool
	fmt.Println("Worker Pool Pattern")
}
