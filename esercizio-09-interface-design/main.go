package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	mu     sync.RWMutex
	data   map[string][]byte
	closed bool
}

type FileStorage struct {
	baseDir string
	closed  bool
	mu      sync.RWMutex
}

func NewFileStorage(baseDir string) (*FileStorage, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	return &FileStorage{baseDir: baseDir}, nil
}

func encodeKey(key string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(key))
}

func (f *FileStorage) pathForKey(key string) string {
	return filepath.Join(f.baseDir, encodeKey(key)+".dat")
}

func (m *MemoryStorage) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, errors.New("storage closed")
	}

	value, ok := m.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	// copy difensiva: evita che il caller modifichi lo stato interno
	copied := append([]byte(nil), value...)
	return copied, nil
}

func (m *MemoryStorage) Put(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("storage closed")
	}

	m.data[key] = append([]byte(nil), value...)
	return nil
}

func (m *MemoryStorage) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("storage closed")
	}

	if _, ok := m.data[key]; !ok {
		return ErrNotFound
	}
	delete(m.data, key)
	return nil
}

func (m *MemoryStorage) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, errors.New("storage closed")
	}

	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

func (m *MemoryStorage) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}
	m.closed = true
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

type CachedStorage struct {
	backend Storage
	cache   map[string][]byte
	mu      sync.RWMutex
}

func (f *FileStorage) Put(key string, value []byte) error {
	f.mu.RLock()
	if f.closed {
		f.mu.RUnlock()
		return errors.New("storage closed")
	}
	f.mu.RUnlock()

	finalPath := f.pathForKey(key)
	tmpPath := finalPath + ".tmp"

	if err := os.WriteFile(tmpPath, append([]byte(nil), value...), 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, finalPath)
}

func (f *FileStorage) Get(key string) ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.closed {
		return nil, errors.New("storage closed")
	}

	data, err := os.ReadFile(f.pathForKey(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (f *FileStorage) Delete(key string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.closed {
		return errors.New("storage closed")
	}

	err := os.Remove(f.pathForKey(key))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

func (f *FileStorage) List() ([]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.closed {
		return nil, errors.New("storage closed")
	}

	entries, err := os.ReadDir(f.baseDir)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".dat" {
			continue
		}
		encoded := strings.TrimSuffix(entry.Name(), ".dat")
		raw, err := base64.RawURLEncoding.DecodeString(encoded)
		if err != nil {
			continue
		}
		keys = append(keys, string(raw))
	}
	sort.Strings(keys)
	return keys, nil
}

func (f *FileStorage) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
	return nil
}

func createStorage(storageType string) (Storage, error) {
	switch storageType {
	case "memory":
		return NewMemoryStorage(), nil
	case "file":
		return NewFileStorage("./data")
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func runDemo(storage Storage) error {
	defer storage.Close()

	if err := storage.Put("user:1", []byte(`{"name":"Alice","age":30}`)); err != nil {
		return err
	}
	if err := storage.Put("user:2", []byte(`{"name":"Bob","age":35}`)); err != nil {
		return err
	}

	value, err := storage.Get("user:1")
	if err != nil {
		return err
	}
	fmt.Println("user:1 =", string(value))

	keys, err := storage.List()
	if err != nil {
		return err
	}
	fmt.Println("keys:", keys)

	if err := storage.Delete("user:1"); err != nil {
		return err
	}
	_, err = storage.Get("user:1")
	if err != nil {
		fmt.Println("after delete get user:1 ->", err)
	}

	return nil
}

func main() {
	for _, storageType := range []string{"memory", "file"} {
		fmt.Println("\n===", storageType, "===")
		storage, err := createStorage(storageType)
		if err != nil {
			panic(err)
		}
		if err := runDemo(storage); err != nil {
			panic(err)
		}
	}
}
