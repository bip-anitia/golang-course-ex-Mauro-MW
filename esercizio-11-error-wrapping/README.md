# Esercizio 11: Error Wrapping

## Obiettivo
Implementare un sistema robusto di gestione errori utilizzando error wrapping, custom error types, e le funzionalità di `errors` package introdotte in Go 1.13+.

## Descrizione
Creare un'applicazione che dimostra le best practice per error handling in Go: wrapping, unwrapping, error types, sentinel errors, e error chains.

## Requisiti

### 1. Error Wrapping (Go 1.13+)

```go
// Wrapping con fmt.Errorf
err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Unwrapping
errors.Unwrap(err)
errors.Is(err, target)
errors.As(err, &target)
```

### 2. Custom Error Types

```go
// Error type con context
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Error con stack trace
type DetailedError struct {
    Op      string    // Operazione
    Path    string    // File/resource path
    Err     error     // Wrapped error
    Time    time.Time
}
```

### 3. Sentinel Errors

```go
var (
    ErrNotFound         = errors.New("not found")
    ErrPermissionDenied = errors.New("permission denied")
    ErrInvalidInput     = errors.New("invalid input")
    ErrTimeout          = errors.New("operation timed out")
)
```

### 4. Error Chain

```
Root Cause Error
    ↓ wrapped by
Validation Error
    ↓ wrapped by
Service Error
    ↓ wrapped by
API Error
```

## Scenari da Implementare

### Scenario 1: File Processing con Error Context

```go
func processFile(filename string) error {
    data, err := readFile(filename)
    if err != nil {
        return fmt.Errorf("process file %s: %w", filename, err)
    }

    if err := validateData(data); err != nil {
        return fmt.Errorf("validate data from %s: %w", filename, err)
    }

    if err := saveResults(data); err != nil {
        return fmt.Errorf("save results for %s: %w", filename, err)
    }

    return nil
}

func readFile(filename string) ([]byte, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("file not found: %w", ErrNotFound)
        }
        if os.IsPermission(err) {
            return nil, fmt.Errorf("permission denied: %w", ErrPermissionDenied)
        }
        return nil, fmt.Errorf("read file: %w", err)
    }
    return data, nil
}
```

### Scenario 2: API Client con Error Types

```go
type APIError struct {
    StatusCode int
    Method     string
    URL        string
    Message    string
    Err        error
}

func (e *APIError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("API %s %s failed (status %d): %s: %v",
            e.Method, e.URL, e.StatusCode, e.Message, e.Err)
    }
    return fmt.Sprintf("API %s %s failed (status %d): %s",
        e.Method, e.URL, e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
    return e.Err
}

type Client struct {
    baseURL string
}

func (c *Client) Get(path string) (*Response, error) {
    resp, err := http.Get(c.baseURL + path)
    if err != nil {
        return nil, &APIError{
            Method:  "GET",
            URL:     path,
            Message: "request failed",
            Err:     err,
        }
    }

    if resp.StatusCode != http.StatusOK {
        return nil, &APIError{
            StatusCode: resp.StatusCode,
            Method:     "GET",
            URL:        path,
            Message:    "unexpected status code",
        }
    }

    // ... parse response
    return parseResponse(resp)
}
```

### Scenario 3: Multi-Error (Aggregating Errors)

```go
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    if len(m.Errors) == 0 {
        return "no errors"
    }
    if len(m.Errors) == 1 {
        return m.Errors[0].Error()
    }
    return fmt.Sprintf("%d errors occurred: %v", len(m.Errors), m.Errors[0])
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.Errors = append(m.Errors, err)
    }
}

func (m *MultiError) HasErrors() bool {
    return len(m.Errors) > 0
}

// Uso: processare multipli file
func processFiles(files []string) error {
    var multiErr MultiError

    for _, file := range files {
        if err := processFile(file); err != nil {
            multiErr.Add(fmt.Errorf("process %s: %w", file, err))
        }
    }

    if multiErr.HasErrors() {
        return &multiErr
    }
    return nil
}
```

### Scenario 4: Error Handling con errors.Is e errors.As

```go
func handleError(err error) {
    // Check sentinel error
    if errors.Is(err, ErrNotFound) {
        fmt.Println("Resource not found, creating new one...")
        return
    }

    // Check error type
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation failed on field '%s': %s\n",
            validationErr.Field, validationErr.Message)
        return
    }

    var apiErr *APIError
    if errors.As(err, &apiErr) {
        if apiErr.StatusCode == 429 {
            fmt.Println("Rate limited, retrying...")
            return
        }
        fmt.Printf("API error: %s\n", apiErr.Message)
        return
    }

    // Generic error
    fmt.Printf("Unexpected error: %v\n", err)
}
```

### Scenario 5: Error con Stack Trace

```go
type StackError struct {
    Err   error
    Stack []uintptr
}

func NewStackError(err error) *StackError {
    const depth = 32
    var pcs [depth]uintptr
    n := runtime.Callers(2, pcs[:])

    return &StackError{
        Err:   err,
        Stack: pcs[0:n],
    }
}

func (e *StackError) Error() string {
    return e.Err.Error()
}

func (e *StackError) Unwrap() error {
    return e.Err
}

func (e *StackError) StackTrace() string {
    frames := runtime.CallersFrames(e.Stack)
    var buf strings.Builder

    for {
        frame, more := frames.Next()
        fmt.Fprintf(&buf, "%s\n\t%s:%d\n",
            frame.Function, frame.File, frame.Line)
        if !more {
            break
        }
    }

    return buf.String()
}
```

## Esempi di Output

### Error Chain
```
Error chain:
  → save results for users.json: database connection failed: connection timeout

Details:
  Operation: save results
  File: users.json
  Cause: database connection failed
  Root cause: connection timeout
```

### Validation Error
```
Validation Error:
  Field: email
  Value: "invalid-email"
  Message: must be a valid email address

Suggestion: Please provide an email in format user@domain.com
```

### API Error
```
API Error:
  Method: POST
  URL: /api/users
  Status: 422 Unprocessable Entity
  Message: Invalid request body

Response body:
{
  "error": "email field is required"
}
```

### Multi-Error
```
Multiple errors occurred (3 total):
  1. process file1.txt: validation failed: invalid format
  2. process file2.txt: file not found
  3. process file3.txt: permission denied

2 files processed successfully
3 files failed
```

## Concetti Go da Usare

- `error` interface
- `fmt.Errorf()` con `%w` verb per wrapping
- `errors.New()` per sentinel errors
- `errors.Is()` per comparare con sentinel
- `errors.As()` per type assertion su error chain
- `errors.Unwrap()` per unwrapping manuale
- Custom error types con `Error() string` method
- `Unwrap() error` method per custom wrappers
- Error composition e chaining

## Struttura Suggerita

```go
package main

import (
    "errors"
    "fmt"
    "os"
    "time"
)

// Sentinel errors
var (
    ErrNotFound      = errors.New("not found")
    ErrInvalidInput  = errors.New("invalid input")
    ErrUnauthorized  = errors.New("unauthorized")
)

// Custom error types
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s (got: %v)",
        e.Field, e.Message, e.Value)
}

type OpError struct {
    Op   string
    Path string
    Err  error
    Time time.Time
}

func (e *OpError) Error() string {
    return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
}

func (e *OpError) Unwrap() error {
    return e.Err
}

// Example functions
func validateEmail(email string) error {
    if email == "" {
        return &ValidationError{
            Field:   "email",
            Value:   email,
            Message: "cannot be empty",
        }
    }
    if !strings.Contains(email, "@") {
        return &ValidationError{
            Field:   "email",
            Value:   email,
            Message: "must contain @",
        }
    }
    return nil
}

func createUser(name, email string) error {
    if err := validateEmail(email); err != nil {
        return fmt.Errorf("create user: %w", err)
    }

    // Simulate database error
    if err := saveToDatabase(name, email); err != nil {
        return &OpError{
            Op:   "create user",
            Path: "database",
            Err:  err,
            Time: time.Now(),
        }
    }

    return nil
}

func saveToDatabase(name, email string) error {
    // Simulate error
    return fmt.Errorf("connection failed: %w", ErrTimeout)
}

func main() {
    // Test error handling
    err := createUser("Alice", "invalid-email")

    // Handle with errors.Is
    if errors.Is(err, ErrTimeout) {
        fmt.Println("Timeout occurred, retrying...")
    }

    // Handle with errors.As
    var validErr *ValidationError
    if errors.As(err, &validErr) {
        fmt.Printf("Validation failed: %s\n", validErr.Message)
    }

    var opErr *OpError
    if errors.As(err, &opErr) {
        fmt.Printf("Operation %s failed at %v\n", opErr.Op, opErr.Time)
    }

    // Print full error
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Suggerimenti

### Best Practices

1. **Wrap errors con context**: Aggiungi informazioni utili quando wrapping
   ```go
   return fmt.Errorf("failed to process user %d: %w", userID, err)
   ```

2. **Non wrappare sempre**: A volte è meglio ritornare l'errore originale
   ```go
   if errors.Is(err, ErrNotFound) {
       return err  // Non wrappare, passa originale
   }
   ```

3. **Usa sentinel errors per errori pubblici**: Errori che i caller devono gestire
   ```go
   var ErrNotFound = errors.New("not found")
   ```

4. **Custom types per errori complessi**: Con metadata extra
   ```go
   type ValidationError struct { Field, Message string }
   ```

5. **Non usare panic per errori**: Panic solo per bug del programma
   ```go
   // NO: panic(errors.New("user not found"))
   // YES: return ErrNotFound
   ```

### Pattern Comuni

#### Retry with Exponential Backoff
```go
func retryOperation(op func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = op()
        if err == nil {
            return nil
        }

        // Non ritentare errori permanenti
        var validErr *ValidationError
        if errors.As(err, &validErr) {
            return err
        }

        // Backoff
        time.Sleep(time.Duration(1<<uint(i)) * time.Second)
    }
    return fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
}
```

#### Error Logging
```go
func logError(err error) {
    log.Printf("Error: %v", err)

    // Log error chain
    for err != nil {
        log.Printf("  caused by: %v", err)
        err = errors.Unwrap(err)
    }
}
```

## Challenge Extra

- **Error Codes**: Aggiungi error codes machine-readable
- **I18n**: Supporto per messaggi di errore multilingua
- **Error Recovery**: Strategie di recovery automatico
- **Metrics**: Conta tipi di errori per monitoring
- **Structured Logging**: Integra con structured logger (zap, zerolog)
- **gRPC Errors**: Converti tra Go errors e gRPC status
- **HTTP Errors**: Mappa errors a HTTP status codes
- **User-Friendly Messages**: Separa technical vs user-facing messages

## Testing

```go
func TestErrorWrapping(t *testing.T) {
    err := processFile("nonexistent.txt")

    if !errors.Is(err, ErrNotFound) {
        t.Error("Expected ErrNotFound in error chain")
    }

    if err == nil {
        t.Fatal("Expected error, got nil")
    }

    expectedMsg := "process file nonexistent.txt"
    if !strings.Contains(err.Error(), expectedMsg) {
        t.Errorf("Error should contain %q, got: %v", expectedMsg, err)
    }
}

func TestErrorTypes(t *testing.T) {
    err := validateEmail("invalid")

    var validErr *ValidationError
    if !errors.As(err, &validErr) {
        t.Fatal("Expected ValidationError")
    }

    if validErr.Field != "email" {
        t.Errorf("Expected field 'email', got %q", validErr.Field)
    }
}
```

## Risorse

- Go Blog: [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- Go by Example: [Errors](https://gobyexample.com/errors)
- Effective Go: [Errors](https://go.dev/doc/effective_go#errors)
