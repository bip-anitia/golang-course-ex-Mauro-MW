package main

import (
	"errors"
	"fmt"
	"sync"
)

var ErrNotFound = errors.New("key not found")

type Storage interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Delete(key string) error
	List() ([]string, error)
	Close() error
}

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

type CachedStorage struct {
	backend Storage
	cache   map[string][]byte
	mu      sync.RWMutex
}

func main() {
	// TODO: Implementare il sistema di storage con interfacce
	fmt.Println("Interface Design - Storage System")
}
