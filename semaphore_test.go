package semaphore 

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	rs := New(5)
	assert.Equal(t, 0, rs.Len())
	assert.Equal(t, 5, rs.Cap())
}

func TestResize(t *testing.T) {

	a := assert.New(t)

	t.Run("Basic", func(t *testing.T) {
		rs := New(5)
		a.Equal(5, rs.Cap())

		rs.Resize(1)
		a.Equal(1, rs.Cap())
	})

	t.Run("ResizeUpWhileActive", func(t *testing.T) {
		rs := New(1)
		fmt.Println(rs.String())

		err := rs.Acquire(context.Background())
		a.NoError(err)

		fmt.Println(rs.String())

		go func() {
			time.Sleep(3 * time.Second)
			rs.Resize(2)
		}()

		err = rs.Acquire(context.Background())
		a.NoError(err)

		defer fmt.Println(rs.String())
	})

	t.Run("ResizeDownWhileActive", func(t *testing.T) {
		rs := New(2)
		fmt.Println(rs.String())

		err := rs.Acquire(context.Background())
		a.NoError(err)

		err = rs.Acquire(context.Background())
		a.NoError(err)

		fmt.Println(rs.String())

		rs.Resize(1)
		fmt.Println(rs.String())
		a.Equal(2, rs.Len())
		a.Equal(1, rs.Cap())

		rs.Release()
		a.Equal(1, rs.Len())
		a.Equal(1, rs.Cap())

		fmt.Println(rs.String())
	})
}

func TestAcquire(t *testing.T) {
	a := assert.New(t)

	t.Run("Basic", func(t *testing.T) {
		r := New(1)
		a.NoError(r.Acquire(context.Background()))
	})

	t.Run("Block", func(t *testing.T) {
		r := New(1)
		a.NoError(r.Acquire(context.Background()))

		go func() {
			time.Sleep(3 * time.Second)
			r.Release()
		}()

		a.NoError(r.Acquire(context.Background()))
	})

	t.Run("Cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		r := New(1)
		a.NoError(r.Acquire(ctx))

		go func() {
			time.Sleep(3 * time.Second)
			cancel()
		}()

		err := r.Acquire(ctx)

		a.Error(err)
		a.ErrorIs(err, context.Canceled)
	})
}
