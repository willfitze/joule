package joule

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func TestWorkerpool(t *testing.T) {
	length := 10
	vals := make([]int, length)

	workerFn := func(payload interface{}) error {
		intVal := payload.(int)
		vals[intVal] = intVal * 2

		return nil
	}

	pool := NewPool(workerFn, nil, 10, 50)
	pool.Start(runtime.NumCPU())

	for i := 0; i < length; i++ {
		pool.Enqueue(i)
	}

	pool.Stop()

	for ix, val := range vals {
		if ix*2 != val {
			fmt.Printf("Val %d at index %d does not equal %d\n", val, ix, ix*2)
			t.FailNow()
		}
	}
}

func TestWorkerpoolError(t *testing.T) {
	workerFn := func(payload interface{}) error {
		return errors.New("Test error")
	}

	var expectedErr error
	errorFn := func(payload interface{}, err error) {
		expectedErr = err
	}

	pool := NewPool(workerFn, errorFn, 10, 50)
	pool.Start(runtime.NumCPU())
	pool.Enqueue(10)

	pool.Stop()

	if expectedErr == nil {
		fmt.Println("Error is nil")
		t.FailNow()
	}
}

func BenchmarkWorker(b *testing.B) {
	workerFn := func(payload interface{}) error {
		intVal := payload.(int)
		_ = intVal * 2
		return nil
	}

	pool := NewPool(workerFn, nil, 10, 50)
	pool.Start(runtime.NumCPU())

	for n := 0; n < b.N; n++ {
		pool.Enqueue(10)
	}

	pool.Stop()
}
