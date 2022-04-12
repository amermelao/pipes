package pipes

import "testing"

type products struct {
	A, B float64
}

func newInput(a, b float64) Input[products, float64] {
	returnChannel := make(chan float64)
	return Input[products, float64]{
		Data:   products{A: a, B: b},
		Return: returnChannel,
	}
}

func multiply(input <-chan Input[products, float64]) {
	for value := range input {
		argmuments := value.Data
		total := argmuments.A * argmuments.B

		value.Return <- total
	}

}

func TestProcessor(t *testing.T) {

	passArguments := Apply(multiply)

	args := newInput(4, 5)

	passArguments <- args

	result := <-args.Return

	if result != 20 {
		t.Error("bad result")
	}

}
