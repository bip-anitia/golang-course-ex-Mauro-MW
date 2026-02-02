# Esercizio 10: Pipeline Pattern

## Obiettivo
Implementare il pattern Pipeline per processare stream di dati attraverso una serie di stage collegati tramite channels, con supporto per concorrenza e cancellazione.

## Descrizione
Creare pipeline di processing dove ogni stage riceve dati da un channel, li processa, e invia i risultati al prossimo stage. Implementare vari tipi di pipeline: lineari, fan-out/fan-in, e con cancellazione.

## Requisiti

### 1. Pipeline Base (Lineare)

```
[Generator] -> [Stage1] -> [Stage2] -> [Stage3] -> [Consumer]
```

Ogni stage:
- Riceve input da channel
- Processa i dati
- Invia output a channel successivo
- Gestisce graceful shutdown

### 2. Fan-Out / Fan-In Pattern

```
                  -> [Worker1] -
[Generator] ----- -> [Worker2] -> [Merger] -> [Consumer]
                  -> [Worker3] -
```

- Fan-out: Distribuisce lavoro a multipli worker
- Fan-in: Combina risultati da multipli worker

### 3. Pipeline con Context

Supporto per:
- Cancellazione
- Timeout
- Graceful shutdown di tutta la pipeline

## Esempi di Pipeline

### Pipeline 1: Text Processing

```
[Read Lines] -> [Filter] -> [Transform] -> [Count] -> [Print]
```

1. Read Lines: Legge linee da file
2. Filter: Filtra linee che matchano pattern
3. Transform: Converte a uppercase
4. Count: Conta parole
5. Print: Stampa risultati

### Pipeline 2: Image Processing

```
[Load Images] -> [Resize] -> [Apply Filter] -> [Compress] -> [Save]
```

Con fan-out per processing parallelo.

### Pipeline 3: Data ETL

```
[Extract] -> [Validate] -> [Transform] -> [Aggregate] -> [Load]
```

## Implementazione Base

### Stage Function Signature

```go
// Stage riceve input channel, ritorna output channel
type Stage[In, Out any] func(ctx context.Context, in <-chan In) <-chan Out
```

### Generator Pattern

```go
func generator(ctx context.Context, items ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, item := range items {
            select {
            case out <- item:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

### Processing Stage Pattern

```go
func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

### Fan-Out Pattern

```go
func fanOut(ctx context.Context, in <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = worker(ctx, in, i)
    }
    return channels
}

func worker(ctx context.Context, in <-chan int, id int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            // Process
            result := heavyComputation(n)
            select {
            case out <- result:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

### Fan-In Pattern

```go
func fanIn(ctx context.Context, channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    // Funzione che copia da un channel al output
    multiplex := func(c <-chan int) {
        defer wg.Done()
        for val := range c {
            select {
            case out <- val:
            case <-ctx.Done():
                return
            }
        }
    }

    // Avvia goroutine per ogni input channel
    wg.Add(len(channels))
    for _, c := range channels {
        go multiplex(c)
    }

    // Chiudi output quando tutti finiti
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

## Esempi di Utilizzo

### Esempio 1: Simple Pipeline

```go
func main() {
    ctx := context.Background()

    // Pipeline: Generate -> Square -> Add10 -> Print
    nums := generator(ctx, 1, 2, 3, 4, 5)
    squared := square(ctx, nums)
    added := add10(ctx, squared)

    // Consume
    for result := range added {
        fmt.Println(result)
    }
}

// Output: 11, 14, 19, 26, 35
```

### Esempio 2: Fan-Out/Fan-In

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Generate numeri
    nums := generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

    // Fan-out a 3 workers
    workers := fanOut(ctx, nums, 3)

    // Fan-in per combinare risultati
    results := fanIn(ctx, workers...)

    // Consume
    for result := range results {
        fmt.Println(result)
    }
}
```

### Esempio 3: Pipeline con Timeout

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    nums := slowGenerator(ctx)
    processed := slowProcessor(ctx, nums)

    for result := range processed {
        fmt.Println(result)
    }

    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Pipeline timed out")
    }
}
```

### Esempio 4: Text Processing Pipeline

```go
func main() {
    ctx := context.Background()

    // Pipeline stages
    lines := readLines(ctx, "input.txt")
    filtered := filterLines(ctx, lines, regexp.MustCompile("error"))
    uppercased := toUpperCase(ctx, filtered)
    counted := countWords(ctx, uppercased)

    // Print results
    for count := range counted {
        fmt.Printf("Words: %d\n", count)
    }
}

func readLines(ctx context.Context, filename string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        file, err := os.Open(filename)
        if err != nil {
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            select {
            case out <- scanner.Text():
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func filterLines(ctx context.Context, in <-chan string, pattern *regexp.Regexp) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for line := range in {
            if pattern.MatchString(line) {
                select {
                case out <- line:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}

// ... altri stage
```

## Output Atteso

### Simple Pipeline
```
Pipeline: Generate -> Square -> Add10
Processing: [1, 2, 3, 4, 5]

Stage 1 (Square): 1 -> 1
Stage 2 (Add10):  1 -> 11
Output: 11

Stage 1 (Square): 2 -> 4
Stage 2 (Add10):  4 -> 14
Output: 14

...
Final results: [11, 14, 19, 26, 35]
```

### Fan-Out/Fan-In
```
Fan-Out Pipeline with 3 workers
Distributing work: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

Worker 0: Processing 1
Worker 1: Processing 2
Worker 2: Processing 3
Worker 0: Processing 4
Worker 1: Processing 5
...

Results (order may vary):
2, 4, 6, 8, 10, 12, 14, 16, 18, 20

Completed in 1.2s
```

### Text Processing
```
Processing file: input.txt (1000 lines)

Stage: Read Lines     (1000 lines)
Stage: Filter         (125 lines matched "error")
Stage: Uppercase      (125 lines transformed)
Stage: Count Words    (Total: 1543 words)

Pipeline completed successfully
```

## Concetti Go da Usare

- **Channels** per comunicazione tra stage
- **Goroutines** per eseguire stage in parallelo
- **context.Context** per cancellazione e timeout
- **sync.WaitGroup** per sincronizzazione in fan-in
- **select** statement per controllo cancellazione
- **defer close(ch)** per chiudere channels correttamente
- **range over channel** per consumare dati
- **Buffered channels** per ottimizzazione (opzionale)

## Struttura Suggerita

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Generic pipeline stage
type PipelineStage[In, Out any] func(context.Context, <-chan In) <-chan Out

// Generator: produce valori
func generator[T any](ctx context.Context, values ...T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for _, v := range values {
            select {
            case out <- v:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Map: trasforma ogni valore
func mapStage[In, Out any](ctx context.Context, in <-chan In, fn func(In) Out) <-chan Out {
    out := make(chan Out)
    go func() {
        defer close(out)
        for v := range in {
            select {
            case out <- fn(v):
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Filter: filtra valori
func filterStage[T any](ctx context.Context, in <-chan T, predicate func(T) bool) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for v := range in {
            if predicate(v) {
                select {
                case out <- v:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}

// Fan-out: distribuisce lavoro
func fanOut[T any](ctx context.Context, in <-chan T, numWorkers int, worker func(context.Context, <-chan T) <-chan T) []<-chan T {
    channels := make([]<-chan T, numWorkers)
    for i := 0; i < numWorkers; i++ {
        channels[i] = worker(ctx, in)
    }
    return channels
}

// Fan-in: combina risultati
func fanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
    var wg sync.WaitGroup
    out := make(chan T)

    multiplex := func(c <-chan T) {
        defer wg.Done()
        for v := range c {
            select {
            case out <- v:
            case <-ctx.Done():
                return
            }
        }
    }

    wg.Add(len(channels))
    for _, c := range channels {
        go multiplex(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

// Take: prende solo N elementi
func take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for i := 0; i < n; i++ {
            select {
            case v, ok := <-in:
                if !ok {
                    return
                }
                select {
                case out <- v:
                case <-ctx.Done():
                    return
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func main() {
    ctx := context.Background()

    // TODO: Implementare varie pipeline

    fmt.Println("Pipeline Pattern Examples")
}
```

## Suggerimenti

### Best Practices

1. **Sempre chiudi channels**: Usa `defer close(out)` nel producer
2. **Context ovunque**: Passa context a tutti gli stage per cancellazione
3. **Non chiudere channel di input**: Solo chi crea il channel lo chiude
4. **Range over channels**: Usa `for v := range ch` quando possibile
5. **Buffered channels**: Considera buffered channels per performance
6. **Error handling**: Usa struct con errore per propagare errori

### Error Handling in Pipeline

```go
type Result[T any] struct {
    Value T
    Error error
}

func processWithErrors(ctx context.Context, in <-chan int) <-chan Result[int] {
    out := make(chan Result[int])
    go func() {
        defer close(out)
        for v := range in {
            result, err := riskyOperation(v)
            select {
            case out <- Result[int]{Value: result, Error: err}:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

### Bounded Parallelism

```go
func boundedFanOut[T any](ctx context.Context, in <-chan T, maxWorkers int, worker func(T) T) <-chan T {
    out := make(chan T)
    sem := make(chan struct{}, maxWorkers)

    go func() {
        defer close(out)
        for v := range in {
            // Acquire semaphore
            select {
            case sem <- struct{}{}:
            case <-ctx.Done():
                return
            }

            go func(val T) {
                defer func() { <-sem }() // Release
                result := worker(val)
                select {
                case out <- result:
                case <-ctx.Done():
                }
            }(v)
        }

        // Wait for all workers
        for i := 0; i < maxWorkers; i++ {
            sem <- struct{}{}
        }
    }()

    return out
}
```

## Challenge Extra

- **Retry Logic**: Stage che riprova operazioni fallite
- **Batching**: Stage che raggruppa N elementi insieme
- **Throttling**: Limita rate di processing
- **Buffering**: Buffer intelligente che si adatta al carico
- **Monitoring**: Metriche per ogni stage (throughput, latency)
- **Dynamic Pipeline**: Modifica pipeline a runtime
- **Branching**: Pipeline con branch condizionali
- **Merging**: Combina multipli input stream con priorità
- **Ordered Fan-In**: Mantieni ordine originale dopo fan-out
- **Backpressure**: Gestisci quando consumer è lento

## Pattern Avanzati

### Tee (Duplicate Stream)

```go
func tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
    out1 := make(chan T)
    out2 := make(chan T)
    go func() {
        defer close(out1)
        defer close(out2)
        for v := range in {
            v1, v2 := v, v
            for i := 0; i < 2; i++ {
                select {
                case out1 <- v1:
                    out1 = nil
                case out2 <- v2:
                    out2 = nil
                case <-ctx.Done():
                    return
                }
            }
            out1, out2 = out1, out2
        }
    }()
    return out1, out2
}
```

### Bridge (Flatten Channel of Channels)

```go
func bridge[T any](ctx context.Context, chanStream <-chan <-chan T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for {
            var stream <-chan T
            select {
            case maybeStream, ok := <-chanStream:
                if !ok {
                    return
                }
                stream = maybeStream
            case <-ctx.Done():
                return
            }

            for val := range stream {
                select {
                case out <- val:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}
```

## Testing

```go
func TestPipeline(t *testing.T) {
    ctx := context.Background()

    // Input
    input := []int{1, 2, 3, 4, 5}

    // Pipeline: double -> add 10
    gen := generator(ctx, input...)
    doubled := mapStage(ctx, gen, func(n int) int { return n * 2 })
    added := mapStage(ctx, doubled, func(n int) int { return n + 10 })

    // Collect results
    var results []int
    for v := range added {
        results = append(results, v)
    }

    // Verify
    expected := []int{12, 14, 16, 18, 20}
    if !reflect.DeepEqual(results, expected) {
        t.Errorf("Expected %v, got %v", expected, results)
    }
}

func TestPipelineCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    gen := slowGenerator(ctx)
    processed := processor(ctx, gen)

    // Cancel dopo un po'
    go func() {
        time.Sleep(100 * time.Millisecond)
        cancel()
    }()

    count := 0
    for range processed {
        count++
    }

    if count > 10 {
        t.Errorf("Expected pipeline to stop, but processed %d items", count)
    }
}
```
