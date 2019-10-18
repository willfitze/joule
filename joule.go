package joule

import (
	"sync"
	"time"
)

type WorkerFunc func(payload interface{}) error
type ErrorFunc func(payload interface{}, err error)

type Pool struct {
	workerFn   WorkerFunc
	errorFn    ErrorFunc
	nRetries   int
	retryDelay int
	wg         *sync.WaitGroup
	in         chan interface{}
	shutdown   chan bool
}

func NewPool(workerFn WorkerFunc, errorFn ErrorFunc, nRetries, retryDelay int) *Pool {
	return &Pool{
		workerFn:   workerFn,
		errorFn:    errorFn,
		nRetries:   nRetries,
		retryDelay: retryDelay,
		wg:         &sync.WaitGroup{},
		in:         make(chan interface{}),
		shutdown:   make(chan bool),
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
					p.handle(payload)
				}
			}
		}()

	}
}

func (p *Pool) handle(payload interface{}) {
	retries := 0

	for {
		err := p.workerFn(payload)
		if err == nil {
			break
		}

		retries++

		if retries == p.nRetries {
			if p.errorFn != nil {
				p.errorFn(payload, err)
			}

			break
		}

		time.Sleep(time.Duration(p.retryDelay) * time.Millisecond)
	}
}

func (p *Pool) Stop() {
	close(p.shutdown)
	p.wg.Wait()
}
