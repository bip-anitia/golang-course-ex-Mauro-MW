# Esercizio 14: Graceful HTTP Server

## Obiettivo
Imparare a gestire segnali del sistema operativo e implementare graceful shutdown per un HTTP server.

## Descrizione
Creare un HTTP server che risponde correttamente a SIGINT (Ctrl+C) e SIGTERM, completando le richieste in corso prima di terminare.

## Requisiti

### 1. Signal Handling Base

Catturare segnali OS:

```go
import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Channel per ricevere segnali
    sigChan := make(chan os.Signal, 1)

    // Registra interesse per SIGINT e SIGTERM
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // Aspetta segnale
    sig := <-sigChan
    fmt.Printf("Received signal: %v\n", sig)
}
```

### 2. HTTP Server Graceful Shutdown

Server che shutdown correttamente:

```go
func main() {
    server := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }

    // Start server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Aspetta segnale shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan

    // Graceful shutdown con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }

    fmt.Println("Server stopped gracefully")
}
```

### 3. In-Flight Request Tracking

Tracciare richieste in corso:

```go
type Server struct {
    inFlightRequests sync.WaitGroup
}

func (s *Server) trackRequest(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        s.inFlightRequests.Add(1)
        defer s.inFlightRequests.Done()

        next(w, r)
    }
}

func (s *Server) Shutdown() {
    fmt.Println("Waiting for in-flight requests...")
    s.inFlightRequests.Wait()
    fmt.Println("All requests completed")
}
```

### 4. Multiple Signal Types

Gestire diversi segnali:

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan,
    os.Interrupt,    // SIGINT (Ctrl+C)
    syscall.SIGTERM, // SIGTERM (kill)
    syscall.SIGQUIT, // SIGQUIT (Ctrl+\)
)

sig := <-sigChan
switch sig {
case os.Interrupt:
    fmt.Println("Interrupt signal received")
case syscall.SIGTERM:
    fmt.Println("Termination signal received")
case syscall.SIGQUIT:
    fmt.Println("Quit signal received")
}
```

## Scenari da Implementare

### Scenario 1: Basic HTTP Server con Shutdown

```go
func main() {
    // Setup routes
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/slow", slowHandler)
    http.HandleFunc("/health", healthHandler)

    server := &http.Server{
        Addr:         ":8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Start server
    go func() {
        fmt.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    fmt.Println("\nShutting down server...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    fmt.Println("Server exited properly")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    // Simula operazione lenta
    time.Sleep(10 * time.Second)
    w.Write([]byte("Slow operation completed"))
}
```

### Scenario 2: Worker con Graceful Stop

```go
type Worker struct {
    jobs   chan Job
    done   chan struct{}
    quit   chan os.Signal
}

func NewWorker() *Worker {
    w := &Worker{
        jobs: make(chan Job),
        done: make(chan struct{}),
        quit: make(chan os.Signal, 1),
    }
    signal.Notify(w.quit, os.Interrupt, syscall.SIGTERM)
    return w
}

func (w *Worker) Start() {
    go func() {
        for {
            select {
            case job := <-w.jobs:
                w.processJob(job)
            case <-w.quit:
                fmt.Println("Worker received shutdown signal")
                w.shutdown()
                return
            }
        }
    }()
}

func (w *Worker) shutdown() {
    // Completa job corrente
    fmt.Println("Completing current job...")

    // Processa job rimanenti nel channel
    close(w.jobs)
    for job := range w.jobs {
        w.processJob(job)
    }

    close(w.done)
    fmt.Println("Worker stopped gracefully")
}
```

### Scenario 3: Server con Shutdown Hooks

```go
type Server struct {
    http       *http.Server
    shutdownFn []func() error
}

func (s *Server) RegisterShutdownHook(fn func() error) {
    s.shutdownFn = append(s.shutdownFn, fn)
}

func (s *Server) Shutdown(ctx context.Context) error {
    // Shutdown HTTP server
    if err := s.http.Shutdown(ctx); err != nil {
        return err
    }

    // Execute shutdown hooks
    for _, fn := range s.shutdownFn {
        if err := fn(); err != nil {
            log.Printf("Shutdown hook error: %v", err)
        }
    }

    return nil
}

func main() {
    server := &Server{
        http: &http.Server{Addr: ":8080"},
    }

    // Register hooks
    server.RegisterShutdownHook(func() error {
        fmt.Println("Closing database connections...")
        return db.Close()
    })

    server.RegisterShutdownHook(func() error {
        fmt.Println("Flushing logs...")
        return logger.Flush()
    })

    // ... start and shutdown
}
```

### Scenario 4: Timeout e Force Shutdown

```go
func shutdownWithTimeout(server *http.Server, timeout time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    done := make(chan error, 1)
    go func() {
        done <- server.Shutdown(ctx)
    }()

    select {
    case err := <-done:
        if err != nil {
            log.Printf("Graceful shutdown error: %v", err)
        } else {
            fmt.Println("Server shutdown gracefully")
        }
    case <-ctx.Done():
        fmt.Println("Shutdown timeout, forcing exit")
        // Force close
        server.Close()
    }
}
```

## Esempi di Utilizzo

```bash
# Start server
go run main.go

# In altro terminale, fai richieste
curl http://localhost:8080/
curl http://localhost:8080/slow

# Durante slow request, premi Ctrl+C nel primo terminale
# Server aspetta che /slow finisca prima di fermarsi
```

## Output Atteso

### Normal Operation
```
Server starting on :8080
Press Ctrl+C to shutdown

[2024-03-13 10:30:15] GET / - 200 OK (1ms)
[2024-03-13 10:30:20] GET /health - 200 OK (0ms)
[2024-03-13 10:30:25] GET /slow - Processing...
```

### Graceful Shutdown
```
^C
Shutdown signal received (SIGINT)

Shutting down server gracefully...
Waiting for 2 in-flight requests to complete...

[2024-03-13 10:30:35] GET /slow - 200 OK (10s)
All in-flight requests completed

Running shutdown hooks:
  ✓ Closing database connections
  ✓ Flushing logs
  ✓ Cleaning up temp files

Server stopped gracefully
Total uptime: 5m 20s
```

### Timeout Shutdown
```
^C
Shutdown signal received

Shutting down server gracefully...
Waiting for 5 in-flight requests to complete...

Shutdown timeout (30s exceeded)
⚠️  Forcing shutdown - 2 requests may be incomplete

Server stopped (forced)
```

### Multiple Signals
```
Server running on :8080

^C
First interrupt received - initiating graceful shutdown...
Press Ctrl+C again to force quit

Waiting for requests to complete... (10s elapsed)

^C
Second interrupt received - forcing immediate shutdown

Server stopped (forced after 10s)
```

## Concetti Go da Usare

- `os/signal` package
- `signal.Notify()` per registrare signals
- `os.Interrupt` (SIGINT - Ctrl+C)
- `syscall.SIGTERM` (terminate signal)
- `http.Server.Shutdown()` con context
- `context.WithTimeout()` per shutdown timeout
- `sync.WaitGroup` per tracking requests
- Goroutines per concurrent operations
- `select` statement per multiplexing signals

## Struttura Suggerita

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// Server con request tracking
type Server struct {
    http             *http.Server
    inFlightRequests sync.WaitGroup
    shutdownHooks    []func() error
}

func NewServer(addr string) *Server {
    s := &Server{
        http: &http.Server{
            Addr:         addr,
            ReadTimeout:  10 * time.Second,
            WriteTimeout: 10 * time.Second,
        },
    }
    return s
}

// Middleware per tracciare requests
func (s *Server) trackRequests(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s.inFlightRequests.Add(1)
        defer s.inFlightRequests.Done()

        log.Printf("Request started: %s %s", r.Method, r.URL.Path)
        start := time.Now()

        next.ServeHTTP(w, r)

        log.Printf("Request completed: %s %s (took %v)",
            r.Method, r.URL.Path, time.Since(start))
    })
}

// Registra shutdown hook
func (s *Server) OnShutdown(fn func() error) {
    s.shutdownHooks = append(s.shutdownHooks, fn)
}

// Start server
func (s *Server) Start() error {
    log.Printf("Starting server on %s", s.http.Addr)
    if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
        return err
    }
    return nil
}

// Graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
    log.Println("Shutting down server...")

    // Shutdown HTTP server (smette di accettare nuove connessioni)
    if err := s.http.Shutdown(ctx); err != nil {
        return fmt.Errorf("server shutdown error: %w", err)
    }

    // Aspetta che richieste in corso completino
    log.Println("Waiting for in-flight requests...")
    done := make(chan struct{})
    go func() {
        s.inFlightRequests.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Println("All requests completed")
    case <-ctx.Done():
        return fmt.Errorf("shutdown timeout: %w", ctx.Err())
    }

    // Esegui shutdown hooks
    log.Println("Running shutdown hooks...")
    for i, hook := range s.shutdownHooks {
        if err := hook(); err != nil {
            log.Printf("Shutdown hook %d error: %v", i, err)
        }
    }

    return nil
}

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello! Server is running.\n")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    // Simula operazione lenta
    log.Println("Starting slow operation...")
    time.Sleep(10 * time.Second)
    fmt.Fprintf(w, "Slow operation completed\n")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK\n")
}

func main() {
    // Crea server
    server := NewServer(":8080")

    // Setup routes con tracking middleware
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/slow", slowHandler)
    mux.HandleFunc("/health", healthHandler)

    server.http.Handler = server.trackRequests(mux)

    // Registra shutdown hooks
    server.OnShutdown(func() error {
        log.Println("Cleanup: closing database")
        time.Sleep(500 * time.Millisecond) // Simula cleanup
        return nil
    })

    server.OnShutdown(func() error {
        log.Println("Cleanup: flushing logs")
        time.Sleep(200 * time.Millisecond)
        return nil
    })

    // Start server in goroutine
    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    sig := <-quit
    log.Printf("\nReceived signal: %v", sig)

    // Graceful shutdown con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Shutdown failed: %v", err)
    }

    log.Println("Server exited properly")
}
```

## Suggerimenti

### Best Practices

1. **Buffered signal channel**
   ```go
   // ✅ GOOD: buffered channel
   sigChan := make(chan os.Signal, 1)

   // ❌ BAD: unbuffered (può perdere segnali)
   sigChan := make(chan os.Signal)
   ```

2. **Sempre context con timeout per shutdown**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   server.Shutdown(ctx)
   ```

3. **HTTP server in goroutine**
   ```go
   go func() {
       if err := server.ListenAndServe(); err != http.ErrServerClosed {
           log.Fatal(err)
       }
   }()
   ```

4. **Stop signal notification quando non più necessario**
   ```go
   signal.Stop(sigChan)
   ```

### Pattern Comuni

#### Double Ctrl+C per force quit
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt)

<-quit
fmt.Println("Shutting down gracefully... (press Ctrl+C again to force)")

go gracefulShutdown()

<-quit
fmt.Println("Force shutdown!")
os.Exit(1)
```

#### Cleanup con defer
```go
func main() {
    // Setup
    db := openDatabase()
    defer db.Close()

    cache := openCache()
    defer cache.Close()

    // ... rest of main
}
```

## Challenge Extra

- **Multi-Service Shutdown**: Coordina shutdown di HTTP server + worker pool + database
- **Shutdown Metrics**: Traccia tempo di shutdown, richieste in-flight, etc.
- **Rolling Restart**: Zero-downtime restart
- **Health Check durante Shutdown**: Endpoint `/health` ritorna "unhealthy"
- **Configurable Timeout**: Timeout configurabile per shutdown
- **Shutdown Progress**: Progress bar durante shutdown

## Testing

```go
func TestGracefulShutdown(t *testing.T) {
    server := NewServer(":8081")

    // Start server
    go server.Start()

    // Give server time to start
    time.Sleep(100 * time.Millisecond)

    // Start slow request
    done := make(chan bool)
    go func() {
        resp, err := http.Get("http://localhost:8081/slow")
        if err != nil {
            t.Errorf("Request failed: %v", err)
        }
        resp.Body.Close()
        done <- true
    }()

    // Give request time to start
    time.Sleep(100 * time.Millisecond)

    // Shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    go server.Shutdown(ctx)

    // Request should complete
    select {
    case <-done:
        // Success
    case <-time.After(12 * time.Second):
        t.Error("Request didn't complete during shutdown")
    }
}
```

## Risorse

- [os/signal documentation](https://pkg.go.dev/os/signal)
- [http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- [Graceful Shutdown in Go](https://go.dev/blog/context)
