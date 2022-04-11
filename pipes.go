package pipes

import "fmt"

var ErrorIsClosed error = fmt.Errorf("pipe is closed")

type OneInNOut[K any] interface {
	NewOutput() <-chan K
	Add(chan<- K)
	Send(value K) error
	Close()
}

func NewPipe[K any]() OneInNOut[K] {
	return MewSimplePipe[K]()
}
