# Esercizio 8: Context Propagation

## Obiettivo
Utilizzare `context.Context` per gestire cancellazione, timeout, deadline e propagazione di valori attraverso chiamate di funzioni e goroutines.

## Descrizione
Implementare vari scenari che dimostrano l'uso corretto di context per controllare il ciclo di vita di operazioni, gestire timeout, e propagare valori request-scoped.

## Requisiti

### 1. Tipi di Context

#### Background & TODO
```go
ctx := context.Background()  // Root context
ctx := context.TODO()        // Placeholder quando non sei sicuro quale context usare
```

#### WithCancel
```go
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()  // Cleanup
```

#### WithTimeout
```go
ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
defer cancel()
```

#### WithDeadline
```go
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(parentCtx, deadline)
defer cancel()
```

#### WithValue
```go
ctx := context.WithValue(parentCtx, key, value)
```

### 2. Scenari da Implementare

## Scenario 1: HTTP Request con Timeout

Simulare una richiesta HTTP che deve completare entro un timeout:

```go
func fetchDataWithTimeout(ctx context.Context, url string) (string, error) {
    // TODO:
    // - Creare richiesta HTTP con context
    // - Gestire timeout
    // - Gestire cancellazione
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    data, err := fetchDataWithTimeout(ctx, "https://example.com/api/data")
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Request timed out")
        }
        return
    }
    fmt.Println("Data:", data)
}
```

## Scenario 2: Cancellazione Propagata a Worker

Worker che risponde a cancellazione:

```go
func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
    for {
        select {
        case job := <-jobs:
            // Processa job
            result := processJob(ctx, job)
            results <- result
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled: %v\n", id, ctx.Err())
            return
        }
    }
}

func processJob(ctx context.Context, job int) int {
    // Simulazione: verifica cancellazione durante processing
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():
            return -1  // Job cancelled
        default:
            time.Sleep(100 * time.Millisecond)
            // Do work
        }
    }
    return job * 2
}
```

## Scenario 3: Pipeline con Context

Pipeline di processing che può essere cancellata:

```go
func pipeline(ctx context.Context, nums []int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n * 2:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    for result := range pipeline(ctx, nums) {
        fmt.Println(result)
    }
}
```

## Scenario 4: Context Values (Request ID, User Info)

Propagare metadata attraverso call chain:

```go
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Crea context con request ID
    requestID := generateRequestID()
    ctx := context.WithValue(r.Context(), requestIDKey, requestID)

    // Simula autenticazione
    userID := authenticate(r)
    ctx = context.WithValue(ctx, userIDKey, userID)

    // Passa context alla business logic
    result, err := processRequest(ctx, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(result))
}

func processRequest(ctx context.Context, r *http.Request) (string, error) {
    // Recupera valori dal context
    requestID := ctx.Value(requestIDKey).(string)
    userID := ctx.Value(userIDKey).(int)

    log.Printf("[%s] Processing request for user %d", requestID, userID)

    // Business logic con context
    return fetchUserData(ctx, userID)
}

func fetchUserData(ctx context.Context, userID int) (string, error) {
    requestID := ctx.Value(requestIDKey).(string)
    log.Printf("[%s] Fetching data for user %d", requestID, userID)

    // Simula database query con timeout dal context
    select {
    case <-time.After(1 * time.Second):
        return fmt.Sprintf("Data for user %d", userID), nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

## Scenario 5: Graceful Shutdown

Server che gestisce graceful shutdown:

```go
func startServer(ctx context.Context) error {
    server := &http.Server{Addr: ":8080"}

    // Handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Usa context della richiesta
        handleWithContext(r.Context(), w, r)
    })

    // Start server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()

    // Aspetta segnale di shutdown
    <-ctx.Done()

    // Graceful shutdown con timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return server.Shutdown(shutdownCtx)
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // Gestisci segnali OS
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        fmt.Println("\nShutdown signal received")
        cancel()
    }()

    if err := startServer(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Output Atteso

### Timeout Example
```
Fetching data from https://slow-api.com/data
Request timed out after 2s
Error: context deadline exceeded
```

### Cancellation Example
```
Starting 5 workers...
Worker 1: Processing job 1
Worker 2: Processing job 2
Worker 3: Processing job 3
^C
Shutdown signal received
Worker 1 cancelled: context canceled
Worker 2 cancelled: context canceled
Worker 3 cancelled: context canceled
Cleanup completed
```

### Context Values Example
```
[req-123456] Processing request for user 42
[req-123456] Fetching data for user 42
[req-123456] Query database...
[req-123456] Request completed in 1.2s
```

## Concetti Go da Usare

- `context.Context` interface
- `context.Background()` e `context.TODO()`
- `context.WithCancel()` per cancellazione
- `context.WithTimeout()` per timeout
- `context.WithDeadline()` per deadline assoluta
- `context.WithValue()` per propagare valori
- `ctx.Done()` channel per controllare cancellazione
- `ctx.Err()` per ottenere errore (Canceled o DeadlineExceeded)
- `ctx.Value()` per recuperare valori
- `select` statement per ascoltare cancellazione

## Struttura Suggerita

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Scenario 1: Simple timeout
func withTimeoutExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    result := make(chan string)

    go func() {
        // Simula operazione lenta
        time.Sleep(3 * time.Second)
        result <- "completed"
    }()

    select {
    case res := <-result:
        fmt.Println("Result:", res)
    case <-ctx.Done():
        fmt.Println("Timeout:", ctx.Err())
    }
}

// Scenario 2: Manual cancellation
func withCancellationExample() {
    ctx, cancel := context.WithCancel(context.Background())

    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Cancelling...")
        cancel()
    }()

    longRunningTask(ctx)
}

func longRunningTask(ctx context.Context) {
    ticker := time.NewTicker(200 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("Working...")
        case <-ctx.Done():
            fmt.Println("Task cancelled:", ctx.Err())
            return
        }
    }
}

// TODO: Implementare altri scenari

func main() {
    fmt.Println("=== Timeout Example ===")
    withTimeoutExample()

    fmt.Println("\n=== Cancellation Example ===")
    withCancellationExample()

    // TODO: Altri esempi
}
```

## Suggerimenti

### Best Practices

1. **Sempre defer cancel()**: Anche se context scade, chiama cancel per rilasciare risorse
2. **Context come primo parametro**: Convenzione Go è `func DoSomething(ctx context.Context, ...)`
3. **Non salvare context in struct**: Context è per request-scoped, non per long-lived objects
4. **Non passare nil context**: Usa `context.TODO()` se non sei sicuro
5. **Context values solo per request-scoped data**: Non usare per parametri opzionali
6. **Type-safe keys**: Usa type unexported per keys in WithValue

### Keys Type-Safe
```go
type contextKey string  // Unexported type

const userIDKey contextKey = "userID"

// Helpers
func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}
```

### Propagazione Corretta
```go
// CORRETTO: Propaga context
func handler(ctx context.Context) {
    result := fetchData(ctx)  // Passa context
}

// SBAGLIATO: Ignora context
func handler(ctx context.Context) {
    result := fetchData(context.Background())  // Non usare nuovo context
}
```

## Challenge Extra

- **Context Hierarchy**: Crea tree di context con timeout diversi
- **Custom Context**: Implementa context custom che logga quando viene cancellato
- **Metrics**: Traccia quante operazioni vengono cancellate vs completate
- **Retry with Context**: Implementa retry logic che rispetta context
- **Fan-out with Context**: Lancia N goroutines e cancella tutte se una fallisce
- **Circuit Breaker**: Combina context con circuit breaker pattern
- **Distributed Tracing**: Usa context per propagare trace IDs

## Errori Comuni da Evitare

```go
// ❌ NON fare: Ignorare errore context
select {
case result := <-ch:
    return result
case <-ctx.Done():
    return nil  // Ignora ctx.Err()
}

// ✅ FARE: Gestire errore
select {
case result := <-ch:
    return result, nil
case <-ctx.Done():
    return nil, ctx.Err()
}

// ❌ NON fare: Context in struct
type Server struct {
    ctx context.Context  // Wrong!
}

// ✅ FARE: Context come parametro
func (s *Server) HandleRequest(ctx context.Context, req Request) {
}

// ❌ NON fare: Bloccare senza check context
for _, item := range items {
    process(item)  // Può bloccare a lungo
}

// ✅ FARE: Controlla context periodicamente
for _, item := range items {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        process(item)
    }
}
```

## Testing

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    err := slowOperation(ctx)

    if err != context.DeadlineExceeded {
        t.Errorf("Expected DeadlineExceeded, got %v", err)
    }
}

func TestWithCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    done := make(chan error)
    go func() {
        done <- longTask(ctx)
    }()

    time.Sleep(50 * time.Millisecond)
    cancel()

    err := <-done
    if err != context.Canceled {
        t.Errorf("Expected Canceled, got %v", err)
    }
}
```
