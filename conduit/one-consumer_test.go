package conduit

import (
	"sort"
	"testing"
)

func oneConsumerBaseCase(i logTestBench) {
	pipe := SimpleOneConsumer[int]()

	msg0, msg1 := 0, 1
	in0 := pipe.NewInput()
	in1 := pipe.NewInput()

	go func() {
		in0 <- msg0
		in1 <- msg1
	}()

	response := [2]bool{}

	out0, more0 := <-pipe.Recieve()
	out1, more1 := <-pipe.Recieve()
	response[out0] = true
	response[out1] = true
	switch {
	case !more0:
		i.Fatal("bad pipe")
	case !more1:
		i.Fatal("bad pipe")
	case !response[0]:
		i.Error("msg0 not arrived")
	case !response[1]:
		i.Error("msg1 not arrived")
	}
	close(in1)

	response = [2]bool{}
	go func() {
		in0 <- msg0
	}()

	out0, more0 = <-pipe.Recieve()
	response[out0] = true
	switch {
	case !more0:
		i.Fatal("bad pipe")
	case !response[0]:
		i.Error("msg0 not arrived")
	case response[1]:
		i.Error("msg1 not arrived")
	}
}

func TestOneConsumer(t *testing.T) {
	oneConsumerBaseCase(t)
}

func BenchmarkMultipeInOneOut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		oneConsumerBaseCase(b)
	}
}

func TestMultipleConsumers(t *testing.T) {
	collector := SimpleOneConsumer[int]()
	producer := SimpleOneProducer[int]()

	const NN int = 3

	for range [NN]struct{}{} {
		collectChanel := collector.NewInput()
		producer.Add(collectChanel)
	}

	const messages int = 10

	go func() {
		for cont := range [messages]struct{}{} {
			producer.Send(cont)
		}
		producer.Close()
	}()

	result := []int{}

	for value := range collector.Recieve() {
		result = append(result, value)
	}

	sort.IntSlice(result).Sort()

	if len(result) != messages {
		t.Fatalf(
			"not all messages arrived, expected: %d, got: %d",
			messages,
			len(result),
		)
	}

	if result[0] != 0 || result[messages-1] != messages-1 {
		t.Error("bad values, the messages got scrambled")
	}

}
