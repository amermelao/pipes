package pipes

type simpleOneConsmer[K any] struct {
	out chan K
}

func newSimpleOneConsmer[K any]() NInOneOut[K] {
	return &simpleOneConsmer[K]{
		out: make(chan K),
	}
}

func (pipe *simpleOneConsmer[K]) Recieve() <-chan K {
	return pipe.out
}

func (pipe *simpleOneConsmer[K]) NewInput() chan<- K {
	newIn := make(chan K)
	go func() {
		for {
			value, more := <-newIn
			if !more {
				break
			} else {
				pipe.out <- value
			}
		}
	}()
	return newIn
}
