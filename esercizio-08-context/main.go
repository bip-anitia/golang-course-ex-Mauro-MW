package main

import (
	"context"
	"fmt"
	"time"
)

type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

func main() {
	withTimeoutExample()
	withCancellationExample()
	workerPoolWithContextExample()
	pipelineWithContextExample()
	withValueExample()
}

func withTimeoutExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	resultCh := make(chan string, 1)

	go func() {
		time.Sleep(1 * time.Second) // lavoro lento
		resultCh <- "completed"
	}()

	select {
	case result := <-resultCh:
		fmt.Println("timeout example:", result)
	case <-ctx.Done():
		fmt.Println("timeout example:", ctx.Err())
	}
}

func withCancellationExample() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(700 * time.Millisecond)
		fmt.Println("cancellation example: cancelling now")
		cancel()
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("cancellation example: working...")
		case <-ctx.Done():
			fmt.Println("cancellation example:", ctx.Err())
			return
		}
	}
}

func pipelineWithContextExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	input := []int{1, 2, 3, 4, 5, 6, 7, 8}

	out := pipeline(ctx, input)

	for value := range out {
		fmt.Println("pipeline example:", value)
		time.Sleep(150 * time.Millisecond) // simula consumer lento
	}
	fmt.Println("pipeline example done:", ctx.Err())
}

func pipeline(ctx context.Context, nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, number := range nums {
			select {
			case <-ctx.Done():
				return
			case out <- number * 2:
			}
		}
	}()

	return out
}

func workerPoolWithContextExample() {
	fmt.Println("\nworker pool example: start")

	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	worker := func(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("worker %d stopped: %v\n", id, ctx.Err())
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}
				// Simula lavoro
				time.Sleep(150 * time.Millisecond)

				// invio risultato cancellabile
				select {
				case <-ctx.Done():
					return
				case results <- job * 2:
				}
			}
		}
	}

	for workerID := 1; workerID <= 3; workerID++ {
		go worker(ctx, workerID, jobs, results)
	}

	for jobID := 1; jobID <= 10; jobID++ {
		jobs <- jobID
	}
	close(jobs)

	received := 0
	for received < 10 {
		select {
		case <-ctx.Done():
			fmt.Println("worker pool example done:", ctx.Err())
			return
		case value := <-results:
			fmt.Println("worker pool result:", value)
			received++
		}
	}
	fmt.Println("worker pool example done: all results received")
}

func withValueExample() {
	fmt.Println("\nwith value example: start")

	ctx := context.Background()
	ctx = withRequestID(ctx, "req-12345")
	ctx = withUserID(ctx, 42)

	requestID, okReq := requestIDFromContext(ctx)
	userID, okUser := userIDFromContext(ctx)

	if !okReq || !okUser {
		fmt.Println("with value example: missing context values")
		return
	}

	fmt.Printf("with value example: requestID=%s userID=%d\n", requestID, userID)
}

func withRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func requestIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(requestIDKey).(string)
	return value, ok
}

func withUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(ctx context.Context) (int, bool) {
	value, ok := ctx.Value(userIDKey).(int)
	return value, ok
}
