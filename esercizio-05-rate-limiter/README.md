# Esercizio 5: Rate Limiter

## Obiettivo
Implementare un rate limiter che controlla il numero di operazioni permesse in un intervallo di tempo, usando channels e time.Ticker.

## Descrizione
Creare diversi tipi di rate limiters per controllare la frequenza di esecuzione di operazioni (es. API calls, richieste HTTP, processing tasks).

## Requisiti

### 1. Token Bucket Rate Limiter
Implementare un rate limiter basato sul pattern "token bucket":
- Bucket con capacità massima di N tokens
- Genera 1 token ogni X millisecondi
- Un'operazione consuma 1 token
- Se non ci sono token, l'operazione deve attendere

```go
type RateLimiter struct {
    tokens    chan struct{}
    maxTokens int
    refillRate time.Duration
}
```

### 2. Sliding Window Rate Limiter
Implementare un rate limiter con finestra temporale scorrevole:
- Massimo N richieste per intervallo di tempo
- Traccia timestamp delle richieste
- Rimuove richieste vecchie dalla finestra

### 3. Concurrent Request Limiter
Limitare il numero di operazioni concorrenti:
- Massimo N operazioni simultanee
- Usare buffered channel come semaforo

## Funzionalità da Implementare

### Token Bucket

```go
// Crea un rate limiter: 5 richieste per secondo
limiter := NewRateLimiter(5, time.Second)

// Aspetta per ottenere permesso
limiter.Wait()
// Esegue operazione
doWork()

// Oppure con timeout
if limiter.TryWait(100 * time.Millisecond) {
    doWork()
} else {
    fmt.Println("Rate limit exceeded")
}
```

### Testing del Rate Limiter

Creare un programma che:
1. Simula N workers che fanno richieste
2. Applica rate limiting
3. Misura e stampa le statistiche

## Esempi di Utilizzo

```bash
# Test base del rate limiter
go run main.go

# Con parametri configurabili
go run main.go -rate=10 -workers=20 -duration=5s
```

## Output Atteso

```
Rate Limiter Test
Configuration:
  Rate: 5 requests/second
  Workers: 10
  Duration: 10s

Starting test...

[00:00:00] Worker 1: Request completed
[00:00:00] Worker 2: Request completed
[00:00:00] Worker 3: Request completed
[00:00:00] Worker 4: Request completed
[00:00:00] Worker 5: Request completed
[00:00:01] Worker 6: Request completed (rate limited)
[00:00:01] Worker 7: Request completed (rate limited)
...

Statistics:
  Total requests: 50
  Successful: 50
  Rate limited: 45
  Average wait time: 200ms
  Actual rate: 5.0 req/s
```

## Implementazioni da Realizzare

### 1. Simple Rate Limiter (con time.Ticker)
```go
// Permette 1 operazione ogni 200ms (5/sec)
limiter := time.Tick(200 * time.Millisecond)

for i := 0; i < requests; i++ {
    <-limiter  // Aspetta il prossimo tick
    doWork(i)
}
```

### 2. Burst Rate Limiter
Permette burst iniziale, poi rate costante:
```go
// Burst di 3, poi 1 per secondo
limiter := NewBurstRateLimiter(3, time.Second)
```

### 3. Per-User Rate Limiter
Rate limiting per chiave (es. user ID, IP address):
```go
limiter := NewPerKeyRateLimiter(10, time.Minute)
if limiter.Allow(userID) {
    handleRequest()
} else {
    returnTooManyRequests()
}
```

## Concetti Go da Usare

- `time.Ticker` per generare eventi periodici
- `time.NewTimer()` per timeout
- `chan` come token bucket
- Buffered channels come semaforo
- `select` per timeout e cancellazione
- `sync.Mutex` per proteggere state condiviso
- `time.Since()` per misurare durate
- Goroutines per simulare workers concorrenti

## Struttura Suggerita

```go
package main

import (
    "context"
    "time"
)

// Token Bucket Rate Limiter
type TokenBucketLimiter struct {
    tokens     chan struct{}
    ticker     *time.Ticker
    maxTokens  int
    refillRate time.Duration
}

func NewTokenBucketLimiter(maxTokens int, refillRate time.Duration) *TokenBucketLimiter {
    // TODO: Implementare
    return nil
}

func (rl *TokenBucketLimiter) Wait() {
    // TODO: Aspetta fino a quando un token è disponibile
}

func (rl *TokenBucketLimiter) TryWait(timeout time.Duration) bool {
    // TODO: Prova ad ottenere token con timeout
    return false
}

func (rl *TokenBucketLimiter) Stop() {
    // TODO: Cleanup resources
}

// Sliding Window Rate Limiter
type SlidingWindowLimiter struct {
    requests   []time.Time
    maxRequests int
    window     time.Duration
    mu         sync.Mutex
}

func (sw *SlidingWindowLimiter) Allow() bool {
    // TODO: Controlla se la richiesta è permessa
    return false
}

// Test function
func simulateRequests(limiter *TokenBucketLimiter, workers int, duration time.Duration) {
    // TODO: Simula workers che fanno richieste
}
```

## Suggerimenti

- Inizia con il Token Bucket pattern (più semplice)
- Usa buffered channel con capacità = max tokens
- Lancia goroutine separata per refill dei tokens
- Per testing, usa molti workers per creare pressione
- Misura il tempo effettivo tra richieste
- Gestisci graceful shutdown con context
- Considera burst capacity per casi d'uso reali

## Challenge Extra

- Implementare distributed rate limiter (con Redis)
- Adaptive rate limiting (aggiusta rate dinamicamente)
- Rate limiter middleware per HTTP server
- Statistiche e monitoring (Prometheus metrics)
- Circuit breaker pattern combinato con rate limiting
- Weighted rate limiting (operazioni diverse hanno costi diversi)
- Hierarchical rate limiting (globale + per-user)

## Casi d'Uso Reali

1. **API Client**: Rispettare rate limits di API esterne
2. **HTTP Server**: Proteggere endpoint da abuse
3. **Background Jobs**: Controllare processing rate
4. **Database**: Limitare query rate per evitare overload
5. **External Services**: Rispettare SLA e rate limits

## Testing

Scrivi test che verificano:
- Rate effettivo è rispettato
- Burst handling
- Comportamento sotto carico
- Fairness tra workers
- Resource cleanup
