package fsm

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	counter int
}

func TestMachine_Run(t *testing.T) {
	t.Parallel()

	t.Run("successful transitions", func(t *testing.T) {
		t.Parallel()

		data := &testData{}
		machine := NewMachine(data)

		state1 := StateFunc[testData](func(ctx context.Context, d *testData) (State[testData], error) {
			d.counter++
			return StateFunc[testData](func(ctx context.Context, d *testData) (State[testData], error) {
				d.counter++
				return nil, nil
			}), nil
		})

		err := machine.Run(context.Background(), state1)
		assert.NoError(t, err)
		assert.Equal(t, 2, data.counter)
	})

	t.Run("stops on error", func(t *testing.T) {
		t.Parallel()

		data := &testData{}
		machine := NewMachine(data)
		expectedErr := errors.New("something went wrong")

		state1 := StateFunc[testData](func(ctx context.Context, d *testData) (State[testData], error) {
			d.counter++
			return nil, expectedErr
		})

		err := machine.Run(context.Background(), state1)
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, 1, data.counter)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		t.Parallel()

		data := &testData{}
		machine := NewMachine(data)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		state1 := StateFunc[testData](func(ctx context.Context, d *testData) (State[testData], error) {
			d.counter++
			return nil, nil
		})

		err := machine.Run(ctx, state1)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 0, data.counter)
	})
}
