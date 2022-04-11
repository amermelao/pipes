package pipes

import (
	"sync"
)

type SimplePipe[K any] struct {
	guard sync.RWMutex
	out   []chan<- K

	in     chan K
	active bool
}

func MewSimplePipe[K any]() OneInNOut[K] {
	pipe := SimplePipe[K]{
		out:    make([]chan<- K, 0, 1),
		in:     make(chan K, 1),
		active: true,
	}

	go pipe.run()
	return &pipe
}

func (pipe *SimplePipe[K]) run() {
	for {
		value, more := <-pipe.in

		if more {
			pipe.sendMessageMultiplePipes(value)
		} else {
			pipe.closeOutputPipes()
			break
		}
	}
}

func (pipe *SimplePipe[K]) sendMessageMultiplePipes(message K) {
	pipe.guard.RLock()
	defer pipe.guard.RUnlock()
	var wg sync.WaitGroup
	for _, outPipe := range pipe.out {
		wg.Add(1)
		outPipe := outPipe
		value := message
		go func() {
			outPipe <- value
			wg.Done()
		}()
	}
	wg.Wait()
}

func (pipe *SimplePipe[K]) closeOutputPipes() {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	for _, outPipe := range pipe.out {
		outPipe := outPipe
		go func() {
			close(outPipe)
		}()
	}
}

func (pipe *SimplePipe[K]) NewOutput() <-chan K {
	tmpOut := make(chan K)
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.out = append(pipe.out, tmpOut)

	return tmpOut
}

func (pipe *SimplePipe[K]) Add(newOut chan<- K) {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.out = append(pipe.out, newOut)
}

func (pipe *SimplePipe[K]) Send(value K) error {
	if pipe.active {
		pipe.in <- value
		return nil
	} else {
		return ErrorIsClosed
	}

}

func (pipe *SimplePipe[K]) Close() {
	close(pipe.in)
	pipe.active = false
}
