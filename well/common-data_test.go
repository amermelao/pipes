package well

import (
	"sync"
	"testing"
)

func TestStoreData(t *testing.T) {
	vault := UniqElementVault()

	vault.AddUpdata(1, 1)
	if v, _ := vault.Get(1); v != 1 {
		t.Error("fail to get the correct value")
	}

	vault.AddUpdata(1, 3)
	if v, _ := vault.Get(1); v != 3 {
		t.Error("fail to update the value")
	}
}

func TestDelete(t *testing.T) {
	vault := UniqElementVault()
	vault.AddUpdata(2, 42)

	vault.Remove(2)

	_, ok := vault.Get(2)

	if ok {
		t.Error("fail to delete value")
	}
}

func TestConcurrentAsk(t *testing.T) {
	vault := UniqElementVault()

	var wg sync.WaitGroup
	NN := int64(2000)

	for cont := int64(0); cont < NN; cont++ {
		cont := cont
		wg.Add(1)
		go func() {
			vault.AddUpdata(cont, float64(cont))
			wg.Done()
		}()

	}

	wg.Wait()

	if int64(vault.Len()) != NN {
		t.Error("not all elements were inserted")
	}
}
