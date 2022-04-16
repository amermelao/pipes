package pump

import (
	"sync"
	"testing"
)

type products struct {
	A, B float64
}

func multiply(input <-chan Input[products, float64]) {
	for value := range input {
		argmuments := value.Data
		total := argmuments.A * argmuments.B

		value.SendBack(total)
	}

}

func TestProcessor(t *testing.T) {
	apply := Apply(multiply)
	param := products{A: 4, B: 5}
	if apply(NewInput[products, float64](param)) != 20 {
		t.Error("bad result")
	}
}

func FuzzProcessorMultiProcess(f *testing.F) {
	apply := ApplyNOut(multiply, 3)

	f.Add(42.0, 108.0)
	f.Add(4.0, 18.0)
	f.Add(34.0, 18.0)
	f.Add(4.0, 218.0)
	f.Add(42.0, 18.0)

	f.Fuzz(func(t *testing.T, a, b float64) {
		param := products{A: a, B: b}
		if apply(NewInput[products, float64](param)) != a*b {
			t.Error("bad result")
		}
	})

}

func TestParrallel(t *testing.T) {
	apply := ApplyMIn(multiply)

	var wg sync.WaitGroup

	for cont := 0; cont < 1000; cont++ {
		wg.Add(1)
		cont := cont
		go func() {
			a := 4 + float64(cont)
			var b float64 = 5
			param := products{A: a, B: b}
			if apply(NewInput[products, float64](param)) != a*b {
				t.Error("bad result")
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
