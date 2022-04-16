package conduit

import (
	"sync"
	"time"

	"github.com/amermelao/pipes/pipeserrors"
)

type simpleSplitter[K any] struct {
	guard sync.RWMutex
	out   []chan<- K

	in     chan K
	active bool
}

func SimpleSplitter[K any]() *simpleSplitter[K] {
	pipe := simpleSplitter[K]{
		out:    make([]chan<- K, 0, 1),
		in:     make(chan K, 1),
		active: true,
	}

	go pipe.run()
	return &pipe
}

func (pipe *simpleSplitter[K]) run() {
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

func (pipe *simpleSplitter[K]) sendMessageMultiplePipes(message K) {
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

func (pipe *simpleSplitter[K]) closeOutputPipes() {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	for _, outPipe := range pipe.out {
		outPipe := outPipe
		go func() {
			close(outPipe)
		}()
	}
}

func (pipe *simpleSplitter[K]) NewOutput() <-chan K {
	tmpOut := make(chan K)
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.out = append(pipe.out, tmpOut)

	return tmpOut
}

func (pipe *simpleSplitter[K]) Add(newOut chan<- K) {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.out = append(pipe.out, newOut)
}

func (pipe *simpleSplitter[K]) Send(value K) error {
	if pipe.active {
		pipe.in <- value
		return nil
	} else {
		return pipeserrors.ErrorIsClosed
	}

}

func (pipe *simpleSplitter[K]) Close() {
	close(pipe.in)
	pipe.active = false
}

func SimpleOneProducer[K any]() *simpleOneProducer[K] {
	pipe := simpleOneProducer[K]{
		out:    make([]chan<- K, 0, 1),
		in:     make(chan K, 1),
		active: true,
	}
	go pipe.run()
	return &pipe
}

type simpleOneProducer[K any] struct {
	guard sync.RWMutex
	out   []chan<- K

	in     chan K
	active bool
}

func (pipe *simpleOneProducer[K]) run() {
	for {
		value, more := <-pipe.in

		if more {
			pipe.sendMessage(value)
		} else {
			pipe.closeOutputPipes()
			break
		}
	}
}

func (pipe *simpleOneProducer[K]) closeOutputPipes() {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	for _, outPipe := range pipe.out {
		outPipe := outPipe
		go func() {
			close(outPipe)
		}()
	}
}

func (pipe *simpleOneProducer[K]) sendMessage(message K) {
	pipe.guard.RLock()
	defer pipe.guard.RUnlock()

	if len(pipe.out) < 1 {
		return
	}

	for {
		for _, outPipe := range pipe.out {
			select {
			case outPipe <- message:
				return
			default:
			}
		}
		time.Sleep(time.Microsecond)
	}

}

func (pipe *simpleOneProducer[K]) NewOutput() <-chan K {
	tmpOut := make(chan K)
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.out = append(pipe.out, tmpOut)

	return tmpOut
}

func (pipe *simpleOneProducer[K]) Add(newOut chan<- K) {
	pipe.guard.Lock()
	defer pipe.guard.Unlock()
	pipe.out = append(pipe.out, newOut)
}

func (pipe *simpleOneProducer[K]) Send(value K) error {
	if pipe.active {
		pipe.in <- value
		return nil
	} else {
		return pipeserrors.ErrorIsClosed
	}

}

func (pipe *simpleOneProducer[K]) Close() {
	close(pipe.in)
	pipe.active = false
}
