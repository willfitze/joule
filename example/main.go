package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/qezel/joule"
)

func main() {
	workerFn := func(payload interface{}) error {
		v := payload.(int)
		log.Printf("%d * 2 = %d\n", v, v*2)

		// Fake some work.
		time.Sleep(time.Second * time.Duration(rand.Intn(3)))

		return nil
	}

	pool := joule.NewPool(workerFn, nil, 0, 0)
	pool.Start(2)

	for i := 1; i <= 10; i++ {
		pool.Enqueue(i)
	}

	pool.Stop()
}
