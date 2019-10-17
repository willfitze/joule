package joule

import (
	"sync"
)

type WorkerFunc func(payload interface{}) error
type ErrorFunc func(payload interface{}, err error)

type Pool struct {
	workerFn WorkerFunc
	errorFn  ErrorFunc
	wg       *sync.WaitGroup
	in       chan interface{}
	shutdown chan bool
}

func NewPool(workerFn WorkerFunc, errorFn ErrorFunc) *Pool {
	return &Pool{
		workerFn: workerFn,
		errorFn:  errorFn,
		wg:       &sync.WaitGroup{},
		in:       make(chan interface{}),
		shutdown: make(chan bool),
	}
}

func (p *Pool) Enqueue(payload interface{}) {
	p.in <- payload
}

func (p *Pool) Start(nWorkers int) {
	for i := 0; i < nWorkers; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-p.shutdown:
					return

				case payload := <-p.in:
					err := p.workerFn(payload)
					if err != nil && p.errorFn != nil {
						p.errorFn(payload, err)
					}
				}
			}
		}()
	}
}

func (p *Pool) Stop() {
	close(p.shutdown)
	p.wg.Wait()
}
