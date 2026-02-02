# Esercizio 6: Worker Pool

## Obiettivo
Implementare il pattern Worker Pool per processare task in parallelo con un numero controllato di goroutines worker.

## Descrizione
Creare un sistema di worker pool che distribuisce task a un pool fisso di worker goroutines, evitando di creare troppe goroutines e gestendo correttamente il ciclo di vita dei worker.

## Requisiti

### 1. Componenti del Worker Pool

```go
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
}
```

### 2. Funzionalità Base

- **Start()**: Avvia N worker goroutines
- **Submit(task)**: Invia un task al pool
- **Stop()**: Ferma gracefully il pool
- **Results()**: Channel per ricevere risultati

### 3. Features Avanzate

- Gestire panic nei worker (recovery)
- Timeout per task individuali
- Cancellazione con context
- Statistiche (task processati, errori, tempo medio)

## Implementazione Richiesta

### Pattern Base

```go
func (wp *WorkerPool) Start() {
    for i := 0; i < wp.numWorkers; i++ {
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    // TODO: Implementare loop del worker
    // - Riceve task dal channel
    // - Processa task
    // - Invia risultato
    // - Gestisce shutdown
}

func (wp *WorkerPool) Submit(task Task) {
    // TODO: Invia task al pool
}

func (wp *WorkerPool) Stop() {
    // TODO: Graceful shutdown
    // - Chiude task channel
    // - Aspetta che tutti i worker finiscano
    // - Chiude result channel
}
```

## Esempi di Utilizzo

```bash
# Test base
go run main.go

# Con numero di workers configurabile
go run main.go -workers=10 -tasks=100

# Test di performance
go run main.go -workers=5 -tasks=1000 -benchmark
```

## Caso d'Uso: Image Processing

Implementare un worker pool per processare immagini:

```go
func main() {
    pool := NewWorkerPool(5)
    pool.Start()
    defer pool.Stop()

    // Simula processing di 100 immagini
    for i := 0; i < 100; i++ {
        task := Task{
            ID: i,
            Data: fmt.Sprintf("image_%d.jpg", i),
            Process: processImage,
        }
        pool.Submit(task)
    }

    // Raccogli risultati
    for i := 0; i < 100; i++ {
        result := <-pool.Results()
        if result.Error != nil {
            fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Error)
        } else {
            fmt.Printf("Task %d completed: %v\n", result.TaskID, result.Value)
        }
    }
}

func processImage(data interface{}) (interface{}, error) {
    filename := data.(string)
    // Simula processing
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("processed_%s", filename), nil
}
```

## Output Atteso

```
Worker Pool started with 5 workers

Worker 0: Processing task 1 (image_0.jpg)
Worker 1: Processing task 2 (image_1.jpg)
Worker 2: Processing task 3 (image_2.jpg)
Worker 3: Processing task 4 (image_3.jpg)
Worker 4: Processing task 5 (image_4.jpg)
Worker 0: Task 1 completed in 102ms
Worker 1: Task 2 completed in 98ms
Worker 0: Processing task 6 (image_5.jpg)
...

Statistics:
  Total tasks: 100
  Successful: 98
  Failed: 2
  Total time: 2.1s
  Average time per task: 101ms
  Throughput: 47.6 tasks/sec
```

## Concetti Go da Usare

- Goroutines per worker pool
- Buffered channels per task queue
- `sync.WaitGroup` per aspettare completion
- `context.Context` per cancellazione
- `defer` e `recover()` per gestire panic
- `select` per timeout e cancellazione
- Channel idioms (close, range over channel)
- Graceful shutdown pattern

## Varianti da Implementare

### 1. Worker Pool Semplice
- N worker fissi
- Task queue unbuffered o buffered
- Results channel

### 2. Dynamic Worker Pool
- Scala il numero di worker in base al carico
- Min/max workers
- Idle timeout per worker

### 3. Priority Worker Pool
- Task con priorità diverse
- Multiple queue (high/medium/low priority)
- Worker processano prima task ad alta priorità

### 4. Worker Pool con Context
```go
func (wp *WorkerPool) StartWithContext(ctx context.Context) {
    for i := 0; i < wp.numWorkers; i++ {
        go wp.workerWithContext(ctx, i)
    }
}

func (wp *WorkerPool) workerWithContext(ctx context.Context, id int) {
    for {
        select {
        case task, ok := <-wp.tasks:
            if !ok {
                return // Channel closed
            }
            // Process task with context
            wp.processTaskWithContext(ctx, task)
        case <-ctx.Done():
            return // Context cancelled
        }
    }
}
```

## Suggerimenti

### Best Practices

1. **Buffered Channels**: Usa buffered channel per tasks per evitare blocking
2. **Graceful Shutdown**:
   - Close task channel per segnalare stop
   - Aspetta con WaitGroup che tutti i worker finiscano
   - Poi chiudi results channel
3. **Error Handling**: Cattura panic nei worker con recover
4. **Resource Cleanup**: Usa defer per cleanup
5. **Dimensione Pool**:
   - CPU-bound: num workers ≈ runtime.NumCPU()
   - I/O-bound: num workers > NumCPU()

### Pattern Comuni

```go
// Graceful shutdown
close(wp.tasks)      // No more tasks
wp.wg.Wait()         // Wait for workers
close(wp.results)    // Signal completion

// Panic recovery in worker
defer func() {
    if r := recover(); r != nil {
        wp.results <- Result{Error: fmt.Errorf("panic: %v", r)}
    }
}()

// Task con timeout
select {
case result := <-processWithTimeout(task, 5*time.Second):
    wp.results <- result
case <-time.After(5 * time.Second):
    wp.results <- Result{Error: errors.New("timeout")}
}
```

## Challenge Extra

- **Fan-out/Fan-in Pattern**: Multiple stage pipeline
- **Retry Logic**: Riprova task falliti con backoff
- **Task Dependencies**: Task che dipendono da altri
- **Metrics**: Prometheus metrics per monitoring
- **Rate Limiting**: Integra con rate limiter
- **Load Balancing**: Distribuisci task basato su worker load
- **Task Timeout**: Timeout individuale per task
- **Dead Letter Queue**: Gestisci task non processabili
- **Backpressure**: Gestisci quando producer è più veloce di consumer

## Testing

```go
func TestWorkerPool(t *testing.T) {
    pool := NewWorkerPool(3)
    pool.Start()
    defer pool.Stop()

    // Submit tasks
    numTasks := 10
    for i := 0; i < numTasks; i++ {
        pool.Submit(Task{
            ID: i,
            Process: func(data interface{}) (interface{}, error) {
                return data.(int) * 2, nil
            },
            Data: i,
        })
    }

    // Collect results
    results := make(map[int]int)
    for i := 0; i < numTasks; i++ {
        result := <-pool.Results()
        if result.Error != nil {
            t.Errorf("Task %d failed: %v", result.TaskID, result.Error)
        }
        results[result.TaskID] = result.Value.(int)
    }

    // Verify
    if len(results) != numTasks {
        t.Errorf("Expected %d results, got %d", numTasks, len(results))
    }
}
```

## Casi d'Uso Reali

1. **Batch Processing**: Processa grandi quantità di dati
2. **Web Scraping**: Scarica multipli URL
3. **Image/Video Processing**: Resize, convert, compress
4. **Data ETL**: Extract, Transform, Load pipelines
5. **API Requests**: Parallelize external API calls
6. **File Processing**: Processa directory di file
