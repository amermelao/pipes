package pipes

import "fmt"

var ErrorIsClosed error = fmt.Errorf("pipe is closed")

type OneInNOut[K any] interface {
	NewOutput() <-chan K
	Add(chan<- K)
	Send(value K) error
	Close()
}

type NInOneOut[K any] interface {
	NewInput() chan<- K
	Recieve() <-chan K
}

