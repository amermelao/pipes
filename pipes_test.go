package pipes

import "testing"

func TestOneWayPipe(t *testing.T) {

	pipe := NewPipe[float64]()

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

	if err != ErrorIsClosed {
		t.Fatal("fail send error")
	}

	_, more = <-outPipe

	if more {
		t.Error("pipe is alive")
	}

}

func TestMultipleOutOneIn(t *testing.T) {

	pipe := NewPipe[float64]()

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

	if err != ErrorIsClosed {
		t.Fatal("fail send error")
	}

	for _, outPipe := range outPipes {
		_, more := <-outPipe

		if more {
			t.Error("pipe is alive")
		}
	}

}
