package pipes

import "github.com/amermelao/pipes/conduit"

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

func NewSimpleOneConsmer[K any]() NInOneOut[K] {
	return conduit.SimpleOneConsumer[K]()
}

func NewSimpleSplitter[K any]() OneInNOut[K] {
	return conduit.SimpleSplitter[K]()
}

func NewSimpleOneProducer[K any]() OneInNOut[K] {
	return conduit.SimpleOneProducer[K]()
}
