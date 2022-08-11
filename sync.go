package sync

import (
	"context"
	"sync"
)

type (
	// Func represents any function that returns a Promise when passed to Execute.
	Func[T, V any] func(context.Context, T) (V, error)

	// Promise represents a potential or actual result from running a Func.
	Promise[V any] struct {
		val  V
		err  error
		done <-chan struct{}
	}
)

// Execute produces a Promise for the supplied Func, evaluating the supplied context.Context and data. The Promise is
// returned immediately, no matter how long it takes for the Func to complete processing.
func Execute[T, V any](ctx context.Context, t T, f Func[T, V]) *Promise[V] {
	done := make(chan struct{})
	p := Promise[V]{
		done: done,
	}
	go func() {
		defer close(done)
		p.val, p.err = f(ctx, t)
	}()
	return &p
}

// Wait takes in zero or more Waiter instances and paused until one returns an error or all d
func Wait[T any](ws []*Promise[T]) {
	var wg sync.WaitGroup
	wg.Add(len(ws))
	done := make(chan struct{})
	for _, w := range ws {
		go func(w *Promise[T]) {
			defer wg.Done()
		}(w)
	}
	go func() {
		defer close(done)
		wg.Wait()
	}()
}

// Get returns the value and the error (if exists) for the Promise. Get waits until the Func associated with this
// Promise has completed. If the Func has completed, Get returns immediately.
func (p *Promise[V]) Get() (V, error) {
	<-p.done
	return p.val, p.err
}
