package pump

import (
	"github.com/amermelao/pipes/conduit"
)

type Input[K any, M any] struct {
	Data          K
	returnChannel chan M
}

func (i Input[K, M]) SendBack(value M) {
	go func() { i.returnChannel <- value }()
}

func NewInput[K any, M any](data K) Input[K, M] {
	channel := make(chan M)
	return Input[K, M]{Data: data, returnChannel: channel}
}

type Process[In any, Out any] func(<-chan Input[In, Out])
type Wrapper[K any, M any] func(Input[K, M]) M
type WrapperF[K any, M any] func(Input[K, M]) func() M

func Apply[In any, Out any](p Process[In, Out]) Wrapper[In, Out] {
	inputChannel := make(chan Input[In, Out])
	go p(inputChannel)

	wrapper := func(data Input[In, Out]) Out {
		inputChannel <- data
		return <-data.returnChannel
	}
	return wrapper
}

func ApplyMIn[In any, Out any](p Process[In, Out]) WrapperF[In, Out] {
	pipe := conduit.SimpleOneConsumer[Input[In, Out]]()
	pipe.NewInput()
	go p(pipe.Recieve())

	wrapper := func(data Input[In, Out]) func() Out {
		inputChannel := pipe.NewInput()
		inputChannel <- data
		close(inputChannel)
		return func() Out { return <-data.returnChannel }
	}
	return wrapper
}

func ApplyNOut[In any, Out any](p Process[In, Out], n int) Wrapper[In, Out] {

	input := conduit.SimpleOneProducer[Input[In, Out]]()

	for range make([]struct{}, n) {
		go p(input.NewOutput())
	}

	wrapper := func(data Input[In, Out]) Out {
		input.Send(data)
		return <-data.returnChannel
	}
	return wrapper
}
