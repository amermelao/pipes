package pipes

type Input[K any, M any] struct {
	Data   K
	Return chan M
}

type Process[In any, Out any] func(<-chan Input[In, Out])

func Apply[In any, Out any](p Process[In, Out]) chan<- Input[In, Out] {
	inputChannel := make(chan Input[In, Out])
	go p(inputChannel)
	return inputChannel
}
