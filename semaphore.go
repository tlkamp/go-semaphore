package semaphore

import (
	"context"
	"fmt"
	"time"

	"github.com/eapache/channels"
)

type ResizableSemaphore struct {
	ch *channels.ResizableChannel
}

// ResizeableSemaphore returns an initialized semaphore with n slots.
func New(n int) *ResizableSemaphore {
	c := channels.NewResizableChannel()
	c.Resize(channels.BufferCap(n))

	return &ResizableSemaphore{
		ch: c,
	}
}

// Acquire will attempt to acquire a slot. Will return an error if the context is canceled.
func (r *ResizableSemaphore) Acquire(ctx context.Context) error {
	select {
	case r.ch.In() <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire will acquire a slot if any are available. A boolean indicating whether or not a slot was acquired is returned.
func (r *ResizableSemaphore) TryAcquire() bool {
	// Configure a small timeout in case the final semaphore is acquired after the check and before
	// this function acquires a semaphore.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	if r.Len() < r.Cap() {
		if err := r.Acquire(ctx); err == nil {
			return true
		}
	}

	return false
}

// Release frees up a slot.
func (r *ResizableSemaphore) Release() {
	<-r.ch.Out()
}

// Resize resizes the underlying channel, increasing or reducing available slots.
func (r *ResizableSemaphore) Resize(n int) {
	if n > 0 {
		r.ch.Resize(channels.BufferCap(n))
	}
}

// Len returns the length of the semaphore (the actively used slots)
func (r *ResizableSemaphore) Len() int {
	return int(r.ch.Len())
}

// Cap returns the capacity of the semaphore (total slots available)
func (r *ResizableSemaphore) Cap() int {
	return int(r.ch.Cap())
}

func (r *ResizableSemaphore) String() string {
	return fmt.Sprintf("Length: %d -- Capacity: %d", r.Len(), r.Cap())
}
