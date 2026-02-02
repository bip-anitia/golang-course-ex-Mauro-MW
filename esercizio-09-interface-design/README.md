# Esercizio 9: Interface Design

## Obiettivo
Progettare un sistema flessibile ed estensibile utilizzando interfacce Go, applicando i principi SOLID e i pattern di design comuni in Go.

## Descrizione
Creare un sistema di storage generico con multiple implementazioni (in-memory, file-based, database) utilizzando interfacce per permettere flessibilità e testabilità.

## Requisiti

### 1. Storage Interface

Progettare un'interfaccia per un sistema di storage chiave-valore:

```go
type Storage interface {
    Get(key string) ([]byte, error)
    Put(key string, value []byte) error
    Delete(key string) error
    List() ([]string, error)
    Close() error
}
```

### 2. Implementazioni Multiple

Implementare almeno 3 storage backend:

#### A. In-Memory Storage
```go
type MemoryStorage struct {
    data map[string][]byte
    mu   sync.RWMutex
}
```

#### B. File Storage
```go
type FileStorage struct {
    baseDir string
}
```

#### C. Cache Storage (Decorator)
```go
type CachedStorage struct {
    backend Storage
    cache   map[string][]byte
    mu      sync.RWMutex
    ttl     time.Duration
}
```

### 3. Interface Composition

Estendere con interfacce opzionali:

```go
// Interfaccia base
type Reader interface {
    Get(key string) ([]byte, error)
}

type Writer interface {
    Put(key string, value []byte) error
    Delete(key string) error
}

// Composizione
type Storage interface {
    Reader
    Writer
    Lister
    io.Closer
}

type Lister interface {
    List() ([]string, error)
}

// Funzionalità opzionali
type BatchWriter interface {
    PutBatch(items map[string][]byte) error
}

type Transactional interface {
    BeginTx() (Transaction, error)
}

type Transaction interface {
    Storage
    Commit() error
    Rollback() error
}
```

### 4. Serialization Layer

Aggiungere layer di serializzazione per oggetti typed:

```go
type Repository[T any] struct {
    storage Storage
}

func (r *Repository[T]) Get(key string) (*T, error) {
    // TODO: Get + deserialize
}

func (r *Repository[T]) Put(key string, obj *T) error {
    // TODO: Serialize + put
}
```

## Esempi di Utilizzo

### Basic Usage

```go
func main() {
    // In-memory storage
    storage := NewMemoryStorage()
    defer storage.Close()

    // Put
    err := storage.Put("user:1", []byte(`{"name":"Alice","age":30}`))
    if err != nil {
        log.Fatal(err)
    }

    // Get
    data, err := storage.Get("user:1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(data))

    // List
    keys, err := storage.List()
    fmt.Println("Keys:", keys)

    // Delete
    storage.Delete("user:1")
}
```

### Swapping Implementations

```go
func createStorage(storageType string) (Storage, error) {
    switch storageType {
    case "memory":
        return NewMemoryStorage(), nil
    case "file":
        return NewFileStorage("/tmp/storage"), nil
    case "cached":
        backend := NewFileStorage("/tmp/storage")
        return NewCachedStorage(backend, 5*time.Minute), nil
    default:
        return nil, fmt.Errorf("unknown storage type: %s", storageType)
    }
}

func main() {
    storage, err := createStorage("cached")
    if err != nil {
        log.Fatal(err)
    }
    defer storage.Close()

    // Usa storage indipendentemente dall'implementazione
    useStorage(storage)
}

func useStorage(s Storage) {
    s.Put("key1", []byte("value1"))
    data, _ := s.Get("key1")
    fmt.Println(string(data))
}
```

### Type-Safe Repository

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    storage := NewMemoryStorage()
    userRepo := NewRepository[User](storage)

    // Put user
    user := &User{
        ID:    "1",
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   30,
    }
    err := userRepo.Put(user.ID, user)

    // Get user
    retrieved, err := userRepo.Get("1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", retrieved)

    // List all users
    users, err := userRepo.List()
    for _, u := range users {
        fmt.Printf("- %s (%s)\n", u.Name, u.Email)
    }
}
```

### Decorator Pattern (Caching)

```go
func main() {
    // Base storage
    fileStorage := NewFileStorage("/tmp/data")

    // Wrap con cache
    cachedStorage := NewCachedStorage(fileStorage, 5*time.Minute)

    // Wrap con logging
    loggedStorage := NewLoggingStorage(cachedStorage)

    // Wrap con metrics
    storage := NewMetricsStorage(loggedStorage)

    // Usa storage con tutti i decorator
    storage.Put("key", []byte("value"))
    // Output: [METRICS] Put called
    //         [LOG] Put(key) = nil
    //         [CACHE] Cache miss, fetching from backend
}
```

## Output Atteso

```
=== In-Memory Storage ===
Put user:1 = OK
Get user:1 = {"name":"Alice","age":30}
List keys = [user:1, user:2, user:3]
Delete user:1 = OK

=== File Storage ===
Put product:100 = OK (saved to /tmp/storage/product/100)
Get product:100 = {"name":"Laptop","price":999.99}
List keys = [product:100, product:101]

=== Cached Storage ===
Put item:1 = OK
Get item:1 = HIT (from cache) {"data":"cached"}
Get item:1 = HIT (from cache) {"data":"cached"}
[After 5min TTL]
Get item:1 = MISS (fetching from backend)

=== Type-Safe Repository ===
User Repository:
  Stored: User{ID:1, Name:Alice, Email:alice@example.com}
  Retrieved: User{ID:1, Name:Alice, Email:alice@example.com}
  All users: [Alice, Bob, Charlie]
```

## Concetti Go da Usare

- **Interfaces**: Definizione e implementazione
- **Interface Composition**: Embedded interfaces
- **Empty Interface**: `interface{}` o `any`
- **Type Assertions**: `value.(Type)`
- **Type Switches**: `switch v := i.(type)`
- **Generics**: `Repository[T any]` (Go 1.18+)
- **Embedding**: Struct embedding per composizione
- **Decorator Pattern**: Wrapper interfaces
- **Factory Pattern**: Funzioni constructor
- **Strategy Pattern**: Swappable implementations

## Struttura Suggerita

```go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

var (
    ErrNotFound = errors.New("key not found")
    ErrInvalidKey = errors.New("invalid key")
)

// Core interfaces
type Storage interface {
    Reader
    Writer
    Lister
    io.Closer
}

type Reader interface {
    Get(key string) ([]byte, error)
}

type Writer interface {
    Put(key string, value []byte) error
    Delete(key string) error
}

type Lister interface {
    List() ([]string, error)
}

// MemoryStorage implementation
type MemoryStorage struct {
    data map[string][]byte
    mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        data: make(map[string][]byte),
    }
}

func (m *MemoryStorage) Get(key string) ([]byte, error) {
    // TODO: Implement
    return nil, nil
}

func (m *MemoryStorage) Put(key string, value []byte) error {
    // TODO: Implement
    return nil
}

func (m *MemoryStorage) Delete(key string) error {
    // TODO: Implement
    return nil
}

func (m *MemoryStorage) List() ([]string, error) {
    // TODO: Implement
    return nil, nil
}

func (m *MemoryStorage) Close() error {
    return nil
}

// FileStorage implementation
type FileStorage struct {
    baseDir string
}

func NewFileStorage(baseDir string) *FileStorage {
    // TODO: Create baseDir if not exists
    return &FileStorage{baseDir: baseDir}
}

// TODO: Implement Storage interface methods

// CachedStorage decorator
type CachedStorage struct {
    backend Storage
    cache   map[string]cacheEntry
    mu      sync.RWMutex
    ttl     time.Duration
}

type cacheEntry struct {
    value     []byte
    expiresAt time.Time
}

func NewCachedStorage(backend Storage, ttl time.Duration) *CachedStorage {
    // TODO: Implement
    return nil
}

// TODO: Implement Storage interface with caching

// Generic Repository
type Repository[T any] struct {
    storage Storage
}

func NewRepository[T any](storage Storage) *Repository[T] {
    return &Repository[T]{storage: storage}
}

func (r *Repository[T]) Get(key string) (*T, error) {
    // TODO: Get + JSON unmarshal
    return nil, nil
}

func (r *Repository[T]) Put(key string, obj *T) error {
    // TODO: JSON marshal + Put
    return nil
}

func (r *Repository[T]) List() ([]*T, error) {
    // TODO: List + unmarshal all
    return nil, nil
}

// Main
func main() {
    fmt.Println("=== Testing Storage Implementations ===")

    // TODO: Test each implementation
}
```

## Suggerimenti

### Design Principles

1. **Accept interfaces, return structs**: Le funzioni dovrebbero accettare interfacce (flessibilità) ma ritornare tipi concreti (chiarezza)

```go
// Good
func ProcessData(r io.Reader) *Result {
    // ...
}

// Less flexible
func ProcessData(f *os.File) *Result {
    // ...
}
```

2. **Small interfaces**: Preferisci molte piccole interfacce a poche grandi (Interface Segregation Principle)

```go
// Good - composable
type Reader interface { Read() }
type Writer interface { Write() }
type ReadWriter interface { Reader; Writer }

// Less flexible - monolithic
type Storage interface {
    Read()
    Write()
    Delete()
    List()
    Backup()
    Restore()
    // ... molti altri metodi
}
```

3. **Interface discovery**: Definisci interfacce dove vengono usate, non dove vengono implementate

```go
// In package consumer
type DataFetcher interface {
    Fetch(id string) ([]byte, error)
}

func ProcessData(fetcher DataFetcher) {
    // usa fetcher
}

// In package provider - implementa implicitamente
type APIClient struct {}
func (c *APIClient) Fetch(id string) ([]byte, error) { ... }
```

### Patterns Comuni

#### Decorator Pattern
```go
type LoggingStorage struct {
    Storage
    logger *log.Logger
}

func (ls *LoggingStorage) Put(key string, value []byte) error {
    ls.logger.Printf("Put(%s)", key)
    return ls.Storage.Put(key, value)
}
```

#### Adapter Pattern
```go
// Adatta io.Reader a nostro Storage
type ReaderAdapter struct {
    reader io.Reader
}

func (ra *ReaderAdapter) Get(key string) ([]byte, error) {
    return io.ReadAll(ra.reader)
}
```

#### Factory Pattern
```go
type StorageConfig struct {
    Type    string
    Options map[string]string
}

func NewStorage(config StorageConfig) (Storage, error) {
    switch config.Type {
    case "memory":
        return NewMemoryStorage(), nil
    case "file":
        return NewFileStorage(config.Options["path"]), nil
    default:
        return nil, errors.New("unknown type")
    }
}
```

## Challenge Extra

- **Plugin System**: Carica storage implementations da plugin
- **Middleware Chain**: Chain di middleware per storage operations
- **Observer Pattern**: Notifiche quando dati cambiano
- **Composite Storage**: Combina multiple storage (es. primary + replica)
- **Encryption Layer**: Decorator che cripta/decripta automaticamente
- **Compression**: Decorator per compressione trasparente
- **Versioning**: Storage che mantiene versioni storiche
- **Schema Validation**: Valida dati contro schema prima di salvare
- **Indexing**: Aggiungi secondary indexes per query efficienti
- **Transactions**: Implementa transazioni ACID

## Testing

```go
func TestStorage(t *testing.T) {
    // Test suite che funziona con qualsiasi Storage implementation
    testCases := []struct {
        name    string
        storage Storage
    }{
        {"Memory", NewMemoryStorage()},
        {"File", NewFileStorage(t.TempDir())},
        {"Cached", NewCachedStorage(NewMemoryStorage(), time.Minute)},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            testStorageImplementation(t, tc.storage)
        })
    }
}

func testStorageImplementation(t *testing.T, s Storage) {
    defer s.Close()

    // Test Put
    err := s.Put("test-key", []byte("test-value"))
    if err != nil {
        t.Fatalf("Put failed: %v", err)
    }

    // Test Get
    value, err := s.Get("test-key")
    if err != nil {
        t.Fatalf("Get failed: %v", err)
    }
    if string(value) != "test-value" {
        t.Errorf("Expected 'test-value', got '%s'", string(value))
    }

    // Test Delete
    err = s.Delete("test-key")
    if err != nil {
        t.Fatalf("Delete failed: %v", err)
    }

    // Verify deleted
    _, err = s.Get("test-key")
    if err != ErrNotFound {
        t.Errorf("Expected ErrNotFound, got %v", err)
    }
}

// Mock per testing
type MockStorage struct {
    GetFunc    func(string) ([]byte, error)
    PutFunc    func(string, []byte) error
    DeleteFunc func(string) error
    ListFunc   func() ([]string, error)
}

func (m *MockStorage) Get(key string) ([]byte, error) {
    if m.GetFunc != nil {
        return m.GetFunc(key)
    }
    return nil, errors.New("not implemented")
}

// ... altri metodi
```

## Best Practices

- Mantieni interfacce piccole e focalizzate
- Implementa implicitamente le interfacce
- Usa composizione invece di ereditarietà
- Restituisci errori espliciti, non panic
- Documenta le interfacce chiaramente
- Usa generics per type-safety quando appropriato
- Testa ogni implementazione con gli stessi test
- Considera performance e memory allocation
