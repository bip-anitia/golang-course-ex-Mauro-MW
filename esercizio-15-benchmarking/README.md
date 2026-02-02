# Esercizio 15: Algorithm Benchmark

## Obiettivo
Imparare a misurare e confrontare performance di diverse implementazioni usando benchmarks, profiling, e analisi delle allocazioni di memoria.

## Descrizione
Creare programmi che confrontano diverse implementazioni dello stesso algoritmo o operazione, usando il package `testing` per benchmarking e `runtime/pprof` per profiling.

## Requisiti

### 1. Basic Benchmarking

Scrivere benchmark con `testing.B`:

```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Codice da benchmarkare
        result := functionToTest()
        _ = result // Previeni ottimizzazione
    }
}
```

### 2. Benchmark con Setup

Setup che non viene misurato:

```go
func BenchmarkWithSetup(b *testing.B) {
    // Setup (non misurato)
    data := generateTestData(1000)

    b.ResetTimer() // Reset timer dopo setup

    for i := 0; i < b.N; i++ {
        process(data)
    }
}
```

### 3. Memory Benchmarking

Misurare allocazioni:

```go
func BenchmarkMemory(b *testing.B) {
    b.ReportAllocs() // Report allocazioni

    for i := 0; i < b.N; i++ {
        result := createLargeSlice()
        _ = result
    }
}
```

### 4. Sub-Benchmarks

Comparare multiple implementazioni:

```go
func BenchmarkSearch(b *testing.B) {
    data := generateData()

    b.Run("Linear", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            linearSearch(data, target)
        }
    })

    b.Run("Binary", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            binarySearch(data, target)
        }
    })
}
```

## Scenari da Implementare

### Scenario 1: String Concatenation Comparison

Confrontare diversi modi di concatenare stringhe:

```go
// Metodo 1: += operator
func concatPlus(n int) string {
    s := ""
    for i := 0; i < n; i++ {
        s += "hello"
    }
    return s
}

// Metodo 2: strings.Builder
func concatBuilder(n int) string {
    var builder strings.Builder
    for i := 0; i < n; i++ {
        builder.WriteString("hello")
    }
    return builder.String()
}

// Metodo 3: bytes.Buffer
func concatBuffer(n int) string {
    var buffer bytes.Buffer
    for i := 0; i < n; i++ {
        buffer.WriteString("hello")
    }
    return buffer.String()
}

// Metodo 4: Preallocated slice
func concatSlice(n int) string {
    parts := make([]string, n)
    for i := 0; i < n; i++ {
        parts[i] = "hello"
    }
    return strings.Join(parts, "")
}

// Benchmarks
func BenchmarkConcat(b *testing.B) {
    sizes := []int{10, 100, 1000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("Plus_%d", size), func(b *testing.B) {
            b.ReportAllocs()
            for i := 0; i < b.N; i++ {
                _ = concatPlus(size)
            }
        })

        b.Run(fmt.Sprintf("Builder_%d", size), func(b *testing.B) {
            b.ReportAllocs()
            for i := 0; i < b.N; i++ {
                _ = concatBuilder(size)
            }
        })

        // ... altri metodi
    }
}
```

### Scenario 2: Search Algorithm Comparison

Confrontare algoritmi di ricerca:

```go
// Linear search
func linearSearch(data []int, target int) int {
    for i, v := range data {
        if v == target {
            return i
        }
    }
    return -1
}

// Binary search (richiede slice ordinato)
func binarySearch(data []int, target int) int {
    left, right := 0, len(data)-1

    for left <= right {
        mid := (left + right) / 2
        if data[mid] == target {
            return mid
        }
        if data[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}

// Map lookup
func mapSearch(data map[int]bool, target int) bool {
    return data[target]
}

func BenchmarkSearch(b *testing.B) {
    sizes := []int{100, 1000, 10000}

    for _, size := range sizes {
        // Generate data
        slice := make([]int, size)
        for i := range slice {
            slice[i] = i
        }

        mapData := make(map[int]bool)
        for _, v := range slice {
            mapData[v] = true
        }

        target := size / 2 // Middle element

        b.Run(fmt.Sprintf("Linear_%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                linearSearch(slice, target)
            }
        })

        b.Run(fmt.Sprintf("Binary_%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                binarySearch(slice, target)
            }
        })

        b.Run(fmt.Sprintf("Map_%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                mapSearch(mapData, target)
            }
        })
    }
}
```

### Scenario 3: Memory Allocation Analysis

Analizzare allocazioni di memoria:

```go
// Versione che alloca
func processSliceNaive(data []int) []int {
    result := []int{} // No capacity specified
    for _, v := range data {
        result = append(result, v*2)
    }
    return result
}

// Versione ottimizzata
func processSliceOptimized(data []int) []int {
    result := make([]int, 0, len(data)) // Preallocated
    for _, v := range data {
        result = append(result, v*2)
    }
    return result
}

// Versione in-place (se possibile)
func processSliceInPlace(data []int) []int {
    result := make([]int, len(data)) // Exact size
    for i, v := range data {
        result[i] = v * 2
    }
    return result
}

func BenchmarkMemoryAllocation(b *testing.B) {
    data := make([]int, 1000)
    for i := range data {
        data[i] = i
    }

    b.Run("Naive", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = processSliceNaive(data)
        }
    })

    b.Run("Optimized", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = processSliceOptimized(data)
        }
    })

    b.Run("InPlace", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = processSliceInPlace(data)
        }
    })
}
```

### Scenario 4: Simple Profiling

Programma principale che può essere profilato:

```go
func main() {
    // CPU profiling
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    // Run workload
    runWorkload()

    // Memory profiling
    if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal(err)
        }
        runtime.GC()
        pprof.WriteHeapProfile(f)
        f.Close()
    }
}

func runWorkload() {
    // Heavy computation
    for i := 0; i < 1000000; i++ {
        _ = fibonacci(20)
    }
}
```

## Comandi di Utilizzo

```bash
# Run benchmarks
go test -bench=.

# Con memory stats
go test -bench=. -benchmem

# Specifica pattern
go test -bench=BenchmarkConcat

# Run multiple volte per precisione
go test -bench=. -benchtime=10s

# CPU profiling
go test -bench=. -cpuprofile=cpu.prof

# Memory profiling
go test -bench=. -memprofile=mem.prof

# Analizza profile
go tool pprof cpu.prof

# Web UI per profiling
go tool pprof -http=:8080 cpu.prof
```

## Output Atteso

### Benchmark Results
```
goos: darwin
goarch: arm64
pkg: example.com/benchmark

BenchmarkConcat/Plus_10-8              2847130    421.3 ns/op    496 B/op    9 allocs/op
BenchmarkConcat/Plus_100-8               34423  34782.0 ns/op  48128 B/op   99 allocs/op
BenchmarkConcat/Plus_1000-8                330 3627845.0 ns/op 4983040 B/op  999 allocs/op

BenchmarkConcat/Builder_10-8          16537251     72.34 ns/op     64 B/op    2 allocs/op
BenchmarkConcat/Builder_100-8          2611633    459.2 ns/op     528 B/op    4 allocs/op
BenchmarkConcat/Builder_1000-8          285115   4197.0 ns/op    8224 B/op    6 allocs/op

BenchmarkConcat/Buffer_10-8           13644391     87.94 ns/op     64 B/op    1 allocs/op
BenchmarkConcat/Buffer_100-8           2284783    525.3 ns/op     560 B/op    2 allocs/op
BenchmarkConcat/Buffer_1000-8           254898   4722.0 ns/op    8320 B/op    3 allocs/op

PASS
ok      example.com/benchmark   15.847s
```

### Performance Comparison Table
```
String Concatenation Performance (1000 iterations):

Method          Time/op       Allocs/op    B/op      Speedup
--------------------------------------------------------------
+=              3.6ms         999          4.9MB     1x (baseline)
strings.Builder 4.2µs         6            8.2KB     857x faster ⚡
bytes.Buffer    4.7µs         3            8.3KB     766x faster ⚡
Join+Slice      5.1µs         2            8.0KB     706x faster ⚡

Winner: strings.Builder
- 857x faster than +=
- 63% fewer allocations than bytes.Buffer
- Recommended for string concatenation
```

### Search Algorithm Comparison
```
Search Performance (10000 elements):

Algorithm       Time/op    Complexity
-------------------------------------
Linear          52.3µs     O(n)
Binary          124ns      O(log n)      ⚡ 422x faster
Map Lookup      89ns       O(1)          ⚡ 587x faster

Winner: Map lookup (if memory permits)
Binary search excellent for sorted data
Linear search only for small datasets
```

## Concetti Go da Usare

- `testing.B` type
- `b.N` iterations
- `b.ResetTimer()` per setup
- `b.ReportAllocs()` per memory stats
- `b.Run()` per sub-benchmarks
- `runtime/pprof` per profiling
- `runtime.MemStats` per memory info
- `testing.AllocsPerRun()` helper
- Benchmark flags: `-bench`, `-benchmem`, `-benchtime`

## Struttura Suggerita

```go
package benchmark

import (
    "bytes"
    "strings"
    "testing"
)

// Implementation 1: Naive
func concatNaive(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s
    }
    return result
}

// Implementation 2: Optimized
func concatOptimized(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

// Benchmark comparison
func BenchmarkStringConcat(b *testing.B) {
    // Test data
    sizes := []int{10, 100, 1000}

    for _, size := range sizes {
        // Generate test data
        strs := make([]string, size)
        for i := range strs {
            strs[i] = "hello"
        }

        b.Run(fmt.Sprintf("Naive_%d", size), func(b *testing.B) {
            b.ReportAllocs()
            for i := 0; i < b.N; i++ {
                _ = concatNaive(strs)
            }
        })

        b.Run(fmt.Sprintf("Optimized_%d", size), func(b *testing.B) {
            b.ReportAllocs()
            for i := 0; i < b.N; i++ {
                _ = concatOptimized(strs)
            }
        })
    }
}

// Helper per stampare risultati
func printBenchmarkResults() {
    // Parsing output di benchmark e formattazione
}
```

### Main con Profiling

```go
// main.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
)

var (
    cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
    memprofile = flag.String("memprofile", "", "write mem profile to file")
)

func main() {
    flag.Parse()

    // CPU profiling
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()

        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal(err)
        }
        defer pprof.StopCPUProfile()
    }

    // Run workload
    start := time.Now()
    runWorkload()
    elapsed := time.Since(start)

    fmt.Printf("Workload completed in %s\n", elapsed)

    // Memory profiling
    if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()

        runtime.GC() // Get up-to-date stats
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal(err)
        }
    }

    // Print memory stats
    printMemStats()
}

func runWorkload() {
    // Example: string concatenation
    iterations := 100000
    for i := 0; i < iterations; i++ {
        _ = buildString(100)
    }
}

func buildString(n int) string {
    var builder strings.Builder
    for i := 0; i < n; i++ {
        builder.WriteString("x")
    }
    return builder.String()
}

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Println("\nMemory Statistics:")
    fmt.Printf("Alloc = %v MB", bToMb(m.Alloc))
    fmt.Printf("\tTotalAlloc = %v MB", bToMb(m.TotalAlloc))
    fmt.Printf("\tSys = %v MB", bToMb(m.Sys))
    fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
```

## Suggerimenti

### Best Practices

1. **Usa b.ReportAllocs()**
   ```go
   func BenchmarkFunc(b *testing.B) {
       b.ReportAllocs() // Mostra allocazioni
       for i := 0; i < b.N; i++ {
           // ...
       }
   }
   ```

2. **ResetTimer dopo setup**
   ```go
   func BenchmarkWithSetup(b *testing.B) {
       data := setupExpensiveData()
       b.ResetTimer() // Reset dopo setup

       for i := 0; i < b.N; i++ {
           process(data)
       }
   }
   ```

3. **Evita ottimizzazioni del compiler**
   ```go
   var result int // Package level

   func BenchmarkFunc(b *testing.B) {
       var r int // Local
       for i := 0; i < b.N; i++ {
           r = expensiveFunc()
       }
       result = r // Assegna a variabile package-level
   }
   ```

4. **Sub-benchmarks per comparazione**
   ```go
   func BenchmarkAll(b *testing.B) {
       b.Run("Method1", benchMethod1)
       b.Run("Method2", benchMethod2)
   }
   ```

### Interpretare Output

```
BenchmarkFunc-8   1000000   1234 ns/op   456 B/op   7 allocs/op
              │         │            │         │              │
              │         │            │         │              └─ Allocazioni per op
              │         │            │         └──────────────── Bytes allocati per op
              │         │            └────────────────────────── Nanoseconds per op
              │         └─────────────────────────────────────── Numero iterazioni
              └───────────────────────────────────────────────── GOMAXPROCS
```

### Pattern Comuni

#### Benchmark con diverse input size
```go
func BenchmarkProcess(b *testing.B) {
    for _, size := range []int{10, 100, 1000, 10000} {
        b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
            data := generateData(size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                process(data)
            }
        })
    }
}
```

#### Comparazione con baseline
```go
func BenchmarkOptimized(b *testing.B) {
    // Baseline
    b.Run("Baseline", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            baselineFunc()
        }
    })

    // Optimized version
    b.Run("Optimized", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            optimizedFunc()
        }
    })
}
```

## Challenge Extra

- **Benchmark Regression**: Detectare performance regression automaticamente
- **Benchmark Comparison Tool**: Tool per comparare risultati di benchmark
- **Automated Reports**: Generare report HTML da benchmark
- **Continuous Benchmarking**: Integrare benchmark in CI/CD
- **Benchmark Different Platforms**: Comparare performance su OS/arch diversi
- **Memory Leak Detection**: Usare profiling per trovare memory leak

## Testing

```go
// Verifica che benchmark funzioni
func TestBenchmarkRuns(t *testing.T) {
    result := testing.Benchmark(BenchmarkMyFunc)
    if result.N == 0 {
        t.Error("Benchmark didn't run")
    }
    if result.NsPerOp() == 0 {
        t.Error("Benchmark result is zero")
    }
}

// Helper per testing allocations
func TestAllocations(t *testing.T) {
    allocs := testing.AllocsPerRun(100, func() {
        _ = myFunction()
    })

    if allocs > 5 {
        t.Errorf("Too many allocations: %f", allocs)
    }
}
```

## Risorse

- [How to write benchmarks in Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [testing package](https://pkg.go.dev/testing)
- [runtime/pprof package](https://pkg.go.dev/runtime/pprof)
