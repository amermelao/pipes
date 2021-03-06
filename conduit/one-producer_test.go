package conduit

import (
	"sort"
	"sync"
	"testing"

	"github.com/amermelao/pipes/pipeserrors"
)

func TestOneWayPipe(t *testing.T) {

	pipe := SimpleSplitter[float64]()

	outPipe := pipe.NewOutput()

	var sendData float64 = 1
	err := pipe.Send(sendData)

	if err != nil {
		t.Fatal("fail to send file")
	}

	recievedData, more := <-outPipe

	if !more {
		t.Fatal("pipe is closed")
	}

	if recievedData != sendData {
		t.Error("wrong data value")
	}

	pipe.Close()

	err = pipe.Send(sendData)

	if err != pipeserrors.ErrorIsClosed {
		t.Fatal("fail send error")
	}

	_, more = <-outPipe

	if more {
		t.Error("pipe is alive")
	}

}

func TestAddExternalPipe(t *testing.T) {

	pipe := SimpleSplitter[float64]()

	outPipe := pipe.NewOutput()
	outExternal := make(chan float64)
	pipe.Add(outExternal)

	var sendData float64 = 1
	err := pipe.Send(sendData)

	if err != nil {
		t.Fatal("fail to send file")
	}

	recievedData, more := <-outPipe

	if !more {
		t.Fatal("pipe is closed")
	}

	if recievedData != sendData {
		t.Error("wrong data value")
	}

	recievedData, more = <-outExternal

	if !more {
		t.Fatal("external pipe is closed")
	}

	if recievedData != sendData {
		t.Error("external wrong data value")
	}

	pipe.Close()

	err = pipe.Send(sendData)

	if err != pipeserrors.ErrorIsClosed {
		t.Fatal("fail send error")
	}

	_, more = <-outPipe

	if more {
		t.Error("pipe is alive")
	}

	_, more = <-outExternal

	if more {
		t.Error("external pipe is alive")
	}

}
func TestMultipleOutOneIn(t *testing.T) {

	pipe := SimpleSplitter[float64]()

	outPipes := []<-chan float64{}

	for cont := 0; cont < 10; cont++ {
		outPipe := pipe.NewOutput()
		outPipes = append(outPipes, outPipe)
	}

	var sendData float64 = 1
	err := pipe.Send(sendData)

	if err != nil {
		t.Fatal("fail to send file")
	}

	for _, outPipe := range outPipes {
		recievedData, more := <-outPipe

		if !more {
			t.Fatal("pipe is closed")
		}

		if recievedData != sendData {
			t.Error("wrong data value")
		}
	}

	pipe.Close()

	err = pipe.Send(sendData)

	if err != pipeserrors.ErrorIsClosed {
		t.Fatal("fail send error")
	}

	for _, outPipe := range outPipes {
		_, more := <-outPipe

		if more {
			t.Error("pipe is alive")
		}
	}

}

func messageOrderCase(t logTestBench) {

	pipe := SimpleSplitter[int]()

	outPipes := []<-chan int{}

	for cont := 0; cont < 300; cont++ {
		outPipe := pipe.NewOutput()
		outPipes = append(outPipes, outPipe)
	}

	NN := 500

	var wg sync.WaitGroup
	wg.Add(1)

	t.Log("run producer")
	go func() {

		for cont := 0; cont < NN; cont++ {
			err := pipe.Send(cont)
			if err != nil {
				panic("fail to send file")
			}
		}

		pipe.Close()
		wg.Done()
	}()

	results := make([][]int, len(outPipes))

	for consumerID, outPipe := range outPipes {
		outPipe := outPipe
		resultID := consumerID
		wg.Add(1)
		go func() {
			for {
				recievedData, more := <-outPipe

				if !more {
					wg.Done()
					break
				} else {
					results[resultID] = append(results[resultID], recievedData)
				}

			}
		}()

	}

	wg.Wait()

	for consumerID, result := range results {
		if len(result) != NN {
			t.Errorf("consumer %d recieved an unexpected amount of data, %d", consumerID, len(results))
		}

		if !sort.IsSorted(sort.IntSlice(result)) {
			t.Errorf("consumer %d data revieved in bad order", consumerID)
		}
	}
}

func TestMessageOrder(t *testing.T) {
	messageOrderCase(t)
}

func BenchmarkMessageOrder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messageOrderCase(b)
	}
}
