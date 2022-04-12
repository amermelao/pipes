package pipes

import (
	"testing"
)

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
	apply := Apply(multiply)
	if apply(newInput(4, 5)) != 20 {
		t.Error("bad result")
	}
}

func FuzzProcessorMultiProcess(f *testing.F) {
	apply := ApplyN(multiply, 3)

	f.Add(42.0, 108.0)
	f.Add(4.0, 18.0)
	f.Add(34.0, 18.0)
	f.Add(4.0, 218.0)
	f.Add(42.0, 18.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		if apply(newInput(a, b)) != a*b {
			t.Error("bad result")
		}
	})

}
