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

func NewSplitter[K any]() OneInNOut[K] {
	return newSimpleSplitter[K]()
}

func NewPipeOneConsumer[K any]() NInOneOut[K] {
	return newSimpleOneConsmer[K]()
}

func NewPipeOneProducer[K any]() OneInNOut[K] {
	return newSimpleSplitter[K]()
}
