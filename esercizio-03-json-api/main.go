package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	PublishYear int       `json:"publish_year"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookStore struct {
	mu     sync.RWMutex
	books  map[string]Book
	nextID int64
}

func main() {
	store := &BookStore{
		books: make(map[string]Book),
	}

}

func (s *BookStore) Get(id string) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

func (s *BookStore) List() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		list = append(list, b)
	}
	return list
}

func (s *BookStore) Create(b Book) Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	b.ID = strconv.FormatInt(s.nextID, 10)
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	s.books[b.ID] = b
	return b

}

func (s *BookStore) Update(id string, b Book) (Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.books[id]
	if !ok {
		return Book{}, false
	}

	b.ID = old.ID
	b.CreatedAt = old.CreatedAt
	s.books[id] = b
	return b, true

}

func (s *BookStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; !ok {
		return false
	}
	delete(s.books, id)
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
