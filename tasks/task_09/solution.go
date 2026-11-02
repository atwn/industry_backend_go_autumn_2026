package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}

	if err := ctx.Err(); err != nil {
		return nil, err // если контекст уже отменён
	}

	if len(in) == 0 {
		return []Result[R]{}, nil // empty result set, but no error
	}

	results := make([]Result[R], len(in))

	var wg sync.WaitGroup
	queue := make(chan int) // works with indices of 'in'

	// initializing workers:
	for range min(workers, len(in)) {
		wg.Go(func() {
			for i := range queue { // pick next job if available
				if ctx.Err() != nil { // NB! also need to check if context was cancelled in the main goroutine to return error
					return // hits each worker
				}

				value, err := fn(ctx, in[i])
				if err != nil {
					results[i].Err = err
					continue // goes to the next job
				}

				results[i].Value = value
			}
		})
	}

	// feed workers with indices from 'in' through the 'queue' channel
loop:
	for i := range in {
		if ctx.Err() != nil {
			break
		}

		select {
		case queue <- i: // feed next worker
		case <-ctx.Done():
			break loop
		}
	}

	close(queue)
	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
