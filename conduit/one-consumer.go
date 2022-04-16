package conduit

import (
	"sync"
)

type simpleOneConsmer[K any] struct {
	out chan K

	guard           sync.Mutex
	openConnections int
}

func SimpleOneConsumer[K any]() *simpleOneConsmer[K] {
	return &simpleOneConsmer[K]{
		out: make(chan K),
	}
}

func (pipe *simpleOneConsmer[K]) Recieve() <-chan K {
	return pipe.out
}

func (pipe *simpleOneConsmer[K]) NewInput() chan<- K {
	newIn := make(chan K)
	go pipe.transmit(newIn)
	pipe.newConnectionOpened()
	return newIn
}

func (pipe *simpleOneConsmer[K]) newConnectionOpened() {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.openConnections++
}

func (pipe *simpleOneConsmer[K]) tryCloseConsumer() {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.openConnections--
	if pipe.openConnections == 0 {
		pipe.close()
	}
}

func (pipe *simpleOneConsmer[K]) transmit(newIn <-chan K) {
	for {
		value, more := <-newIn
		if !more {
			pipe.tryCloseConsumer()
			break
		} else {
			pipe.out <- value
		}
	}
}

func (pipe *simpleOneConsmer[K]) close() {
	close(pipe.out)
}
