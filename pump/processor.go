package pump

import "github.com/amermelao/pipes/conduit"

type Input[K any, M any] struct {
	Data   K
	Return chan M
}

type Process[In any, Out any] func(<-chan Input[In, Out])
type Wrapper[K any, M any] func(Input[K, M]) M

func Apply[In any, Out any](p Process[In, Out]) Wrapper[In, Out] {
	inputChannel := make(chan Input[In, Out])
	go p(inputChannel)

	wrapper := func(data Input[In, Out]) Out {
		inputChannel <- data
		return <-data.Return
	}
	return wrapper
}

func ApplyN[In any, Out any](p Process[In, Out], n int) Wrapper[In, Out] {

	input := conduit.SimpleOneProducer[Input[In, Out]]()

	for range make([]struct{}, n) {
		go p(input.NewOutput())
	}

	wrapper := func(data Input[In, Out]) Out {
		input.Send(data)
		return <-data.Return
	}
	return wrapper
}
