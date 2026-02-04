package main

import (
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucketLimiter struct {
	tokens     chan struct{}
	ticker     *time.Ticker
	maxTokens  int
	refillRate time.Duration
}

func main() {
	rate := flag.Int("rate", 5, "requests per second")
	workers := flag.Int("workers", 10, "number of workers")
	duration := flag.Duration("duration", 10*time.Second, "test duration")
	flag.Parse()

	if *rate <= 0 || *workers <= 0 || *duration <= 0 {
		fmt.Println("invalid config")
		return
	}
	fmt.Println("Rate Limiter Test")
	fmt.Printf("Rate: %d req/s\n", *rate)
	fmt.Printf("Workers: %d\n", *workers)
	fmt.Printf("Duration: %s\n", duration.String())
	fmt.Println("Starting test...")

	refill := time.Second / time.Duration(*rate)
	limiter := NewTokenBucketLimiter(*rate, refill)
	defer limiter.Stop()

	var total int64
	done := time.After(*duration)
	var wg sync.WaitGroup
	wg.Add(*workers)
	for i := 0; i < *workers; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					limiter.Wait()
					atomic.AddInt64(&total, 1)
				}
			}
		}()
	}
	wg.Wait()

	fmt.Printf("Total requests: %d\n", total)
	fmt.Printf("Actual rate: %.2f req/s\n", float64(total)/duration.Seconds())

}

func NewTokenBucketLimiter(maxTokens int, refillRate time.Duration) *TokenBucketLimiter {
	tokens := make(chan struct{}, maxTokens)
	for i := 0; i < maxTokens; i++ {
		tokens <- struct{}{}
	}

	ticker := time.NewTicker(refillRate)
	rl := &TokenBucketLimiter{
		tokens:     tokens,
		ticker:     ticker,
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}

	go func() {
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}

func (rl *TokenBucketLimiter) Wait() {
	<-rl.tokens
}

func (rl *TokenBucketLimiter) TryWait(timeout time.Duration) bool {
	select {
	case <-rl.tokens:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (rl *TokenBucketLimiter) Stop() {
	rl.ticker.Stop()
}
