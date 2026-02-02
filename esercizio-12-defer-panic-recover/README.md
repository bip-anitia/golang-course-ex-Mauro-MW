# Esercizio 12: Safe File Processor

## Obiettivo
Imparare ad usare `defer`, `panic`, e `recover` per gestire risorse e errori in modo sicuro.

## Descrizione
Creare un programma che processa file di dati (CSV, JSON, o testo) garantendo la corretta chiusura delle risorse con `defer` e gestendo situazioni di panic con `recover`.

## Requisiti

### 1. Defer per Resource Management

Usare `defer` per garantire che le risorse vengano sempre rilasciate:

```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // Garantisce chiusura anche in caso di errore

    // Processa il file
    return parseAndProcess(file)
}
```

### 2. Multiple Defer Statements

Capire l'ordine di esecuzione (LIFO - Last In First Out):

```go
func multipleDeferExample() {
    defer fmt.Println("1 - First defer")
    defer fmt.Println("2 - Second defer")
    defer fmt.Println("3 - Third defer")
    fmt.Println("Function body")
}
// Output:
// Function body
// 3 - Third defer
// 2 - Second defer
// 1 - First defer
```

### 3. Panic e Recover

Gestire panic per evitare crash del programma:

```go
func safeCalculate(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    result = divide(a, b) // Potrebbe causare panic se b == 0
    return result, nil
}

func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}
```

### 4. Defer con Named Return Values

Pattern per modificare return value in defer:

```go
func processWithCleanup(filename string) (err error) {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }

    defer func() {
        closeErr := file.Close()
        if err == nil {
            err = closeErr // Se nessun errore, ritorna errore di chiusura
        }
    }()

    return processData(file)
}
```

## Scenari da Implementare

### Scenario 1: CSV File Processor

Processare file CSV con gestione errori completa:

```go
type Record struct {
    Name  string
    Age   int
    Email string
}

func processCSV(filename string) ([]Record, error) {
    // TODO:
    // 1. Apri file con os.Open
    // 2. Usa defer per chiudere
    // 3. Crea csv.Reader
    // 4. Parse righe
    // 5. Gestisci errori di parsing con recover
}
```

### Scenario 2: Safe Calculator

Calculator che non crasha mai:

```go
func calculator() {
    operations := []struct {
        a, b int
        op   string
    }{
        {10, 2, "/"},
        {10, 0, "/"}, // Division by zero!
        {5, 3, "+"},
    }

    for _, op := range operations {
        result, err := safeOperation(op.a, op.b, op.op)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        fmt.Printf("%d %s %d = %d\n", op.a, op.op, op.b, result)
    }
}
```

### Scenario 3: Resource Pool Manager

Gestire pool di risorse con cleanup garantito:

```go
type Resource struct {
    id   int
    data []byte
}

func (r *Resource) Close() error {
    fmt.Printf("Closing resource %d\n", r.id)
    r.data = nil
    return nil
}

func useMultipleResources() error {
    r1 := acquireResource(1)
    defer r1.Close()

    r2 := acquireResource(2)
    defer r2.Close()

    r3 := acquireResource(3)
    defer r3.Close()

    // Anche se panic qui, tutti i defer vengono eseguiti
    return processResources(r1, r2, r3)
}
```

### Scenario 4: Transaction-like Pattern

Pattern per operazioni che richiedono rollback:

```go
func writeFiles(files []string, data []byte) error {
    written := []string{}

    defer func() {
        if r := recover(); r != nil {
            // Rollback: cancella file scritti
            fmt.Println("Rolling back due to panic:", r)
            for _, f := range written {
                os.Remove(f)
            }
        }
    }()

    for _, filename := range files {
        if err := os.WriteFile(filename, data, 0644); err != nil {
            panic(err) // Trigger rollback
        }
        written = append(written, filename)
    }

    return nil
}
```

## Esempi di Utilizzo

```bash
# Processa CSV file
go run main.go process data.csv

# Safe calculator
go run main.go calc

# Test defer order
go run main.go defer-demo

# Test panic recovery
go run main.go panic-demo
```

## Output Atteso

### CSV Processing
```
Processing file: data.csv
Opening file...
Reading records...
Record 1: {Name: Alice, Age: 30, Email: alice@example.com}
Record 2: {Name: Bob, Age: 25, Email: bob@example.com}
Error on line 3: invalid age value
Record 3 skipped due to error
Closing file...
Successfully processed 2 records
```

### Safe Calculator
```
Safe Calculator Demo
10 / 2 = 5
Error: panic recovered: division by zero
5 + 3 = 8
100 / 0 = Error: panic recovered: division by zero
All operations completed (2 errors)
```

### Defer Order Demo
```
Defer Order Demonstration:
Entering function
Defer 1 registered
Defer 2 registered
Defer 3 registered
Function body executed
Exiting function...
Defer 3 executed
Defer 2 executed
Defer 1 executed
```

### Panic Recovery
```
Testing panic recovery:
Opening resource 1
Opening resource 2
Opening resource 3
Processing...
Panic occurred: unexpected error
Closing resource 3
Closing resource 2
Closing resource 1
All resources cleaned up successfully
```

## Concetti Go da Usare

- `defer` statement per cleanup
- LIFO (Last In First Out) execution order
- `panic()` per errori irrecuperabili
- `recover()` per catturare panic
- Named return values con defer
- Resource management pattern
- `os.Open()` e `file.Close()`
- Error handling con defer

## Struttura Suggerita

```go
package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
)

// Scenario 1: File processing with defer
func processFileWithDefer(filename string) (err error) {
    fmt.Println("Opening file:", filename)

    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open: %w", err)
    }

    defer func() {
        fmt.Println("Closing file:", filename)
        if closeErr := file.Close(); closeErr != nil && err == nil {
            err = closeErr
        }
    }()

    // Process file
    // TODO: Implementare
    return nil
}

// Scenario 2: Panic recovery
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
            result = 0
        }
    }()

    if b == 0 {
        panic("division by zero")
    }

    return a / b, nil
}

// Scenario 3: Multiple defer order
func deferOrder() {
    defer fmt.Println("Defer 1")
    defer fmt.Println("Defer 2")
    defer fmt.Println("Defer 3")

    fmt.Println("Function body")
}

// Scenario 4: Resource cleanup
type Resource struct {
    name string
}

func (r *Resource) Close() {
    fmt.Printf("Closing resource: %s\n", r.name)
}

func useResources() {
    r1 := &Resource{"Database"}
    defer r1.Close()

    r2 := &Resource{"File"}
    defer r2.Close()

    r3 := &Resource{"Network"}
    defer r3.Close()

    fmt.Println("Using resources...")
    // Risorse vengono chiuse in ordine inverso: Network, File, Database
}

func main() {
    fmt.Println("=== Defer, Panic, Recover Examples ===\n")

    // Test 1: Defer order
    fmt.Println("1. Defer Order:")
    deferOrder()
    fmt.Println()

    // Test 2: Safe division
    fmt.Println("2. Safe Division:")
    result, err := safeDivide(10, 2)
    fmt.Printf("10 / 2 = %d, error: %v\n", result, err)

    result, err = safeDivide(10, 0)
    fmt.Printf("10 / 0 = %d, error: %v\n", result, err)
    fmt.Println()

    // Test 3: Resource cleanup
    fmt.Println("3. Resource Cleanup:")
    useResources()
    fmt.Println()

    // Test 4: File processing
    fmt.Println("4. File Processing:")
    if err := processFileWithDefer("test.txt"); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Suggerimenti

### Best Practices

1. **Defer subito dopo acquisizione risorsa**
   ```go
   file, err := os.Open(filename)
   if err != nil {
       return err
   }
   defer file.Close() // Subito dopo Open
   ```

2. **Non usare defer in loop (performance)**
   ```go
   // ❌ BAD: defer in loop
   for _, f := range files {
       file, _ := os.Open(f)
       defer file.Close() // Tutti i defer eseguiti alla fine della funzione!
   }

   // ✅ GOOD: usa funzione helper
   for _, f := range files {
       processFile(f) // defer dentro processFile
   }
   ```

3. **Recover solo in defer**
   ```go
   // ✅ CORRECT
   defer func() {
       if r := recover(); r != nil {
           // handle panic
       }
   }()

   // ❌ WRONG: recover fuori defer non funziona
   if r := recover(); r != nil { // Sempre nil!
       // ...
   }
   ```

4. **Panic per bug, error per condizioni attese**
   ```go
   // ✅ GOOD: usa error per condizioni attese
   if user == nil {
       return errors.New("user not found")
   }

   // ❌ BAD: non usare panic per flow control
   if user == nil {
       panic("user not found") // NO!
   }
   ```

### Pattern Comuni

#### Cleanup con errore
```go
func processFile(filename string) (err error) {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }

    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr
        }
    }()

    return process(f)
}
```

#### Panic recovery con logging
```go
func safeHandler() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered: %v\n", r)
            debug.PrintStack()
        }
    }()

    // Risky operations
}
```

#### Lock/Unlock pattern
```go
func criticalSection() {
    mu.Lock()
    defer mu.Unlock()

    // Critical section code
    // unlock automatico anche se panic
}
```

## Challenge Extra

- **Nested Panic**: Gestire panic dentro defer
- **Defer Performance**: Misurare overhead di defer
- **Defer vs Finally**: Differenze con try-finally di altri linguaggi
- **Stack Trace**: Stampare stack trace quando recover
- **Custom Panic Value**: Usare struct custom in panic
- **Goroutine Panic**: Panic in goroutine (non catturato da defer parent!)

## Errori Comuni da Evitare

```go
// ❌ Dimenticare defer Close
func bad() error {
    file, _ := os.Open("file.txt")
    data, _ := io.ReadAll(file)
    file.Close() // Se ReadAll fa panic, file non viene chiuso!
    return nil
}

// ✅ Usare defer
func good() error {
    file, _ := os.Open("file.txt")
    defer file.Close()
    data, _ := io.ReadAll(file)
    return nil
}

// ❌ Recover senza rethrow panic se necessario
defer func() {
    recover() // Sopprime TUTTI i panic, anche quelli seri!
}()

// ✅ Recover selettivo
defer func() {
    if r := recover(); r != nil {
        if err, ok := r.(error); ok {
            // Handle expected errors
        } else {
            panic(r) // Re-panic se non gestibile
        }
    }
}()
```

## Testing

```go
func TestSafeDivide(t *testing.T) {
    tests := []struct {
        a, b     int
        wantErr  bool
        expected int
    }{
        {10, 2, false, 5},
        {10, 0, true, 0},
        {100, 10, false, 10},
    }

    for _, tt := range tests {
        result, err := safeDivide(tt.a, tt.b)

        if tt.wantErr && err == nil {
            t.Errorf("Expected error for %d/%d", tt.a, tt.b)
        }

        if !tt.wantErr && err != nil {
            t.Errorf("Unexpected error: %v", err)
        }

        if result != tt.expected {
            t.Errorf("Got %d, want %d", result, tt.expected)
        }
    }
}

func TestDeferOrder(t *testing.T) {
    var order []int

    func() {
        defer func() { order = append(order, 1) }()
        defer func() { order = append(order, 2) }()
        defer func() { order = append(order, 3) }()
    }()

    expected := []int{3, 2, 1}
    if !reflect.DeepEqual(order, expected) {
        t.Errorf("Got %v, want %v", order, expected)
    }
}
```

## Risorse

- [Defer, Panic, and Recover - Go Blog](https://go.dev/blog/defer-panic-and-recover)
- [Effective Go - Defer](https://go.dev/doc/effective_go#defer)
