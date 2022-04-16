package well

import (
	"github.com/amermelao/pipes/pump"
)

type mapVault struct {
	communication pump.Wrapper[vaultData, vaultResponse]
}

func UniqElementVault() mapVault {

	return mapVault{
		communication: pump.Apply(dictionary),
	}
}

func (m mapVault) AddUpdata(key int64, value float64) {
	params := pump.NewInput[vaultData, vaultResponse](vaultData{
		id:    key,
		value: value,
		op:    addUpdate,
	})
	m.communication(params)
}

func (m mapVault) Get(key int64) (float64, bool) {
	params := pump.NewInput[vaultData, vaultResponse](vaultData{
		id: key,
		op: get,
	})
	data := m.communication(params)
	return data.value, data.ok
}

func (m mapVault) Remove(key int64) {
	params := pump.NewInput[vaultData, vaultResponse](vaultData{
		id: key, op: remove,
	})
	m.communication(params)
}

func (m mapVault) Len() int {
	params := pump.NewInput[vaultData, vaultResponse](vaultData{
		op: length},
	)
	data := m.communication(params)
	return data.len
}

type vaultData struct {
	id    int64
	value float64
	op    vaultOperations
}

type vaultResponse struct {
	value float64
	len   int
	ok    bool
}

type vaultOperations int

const (
	addUpdate vaultOperations = iota
	get
	remove
	length
)

func dictionary(comming <-chan pump.Input[vaultData, vaultResponse]) {
	storage := map[int64]float64{}
	for request := range comming {
		var response vaultResponse
		switch request.Data.op {
		case addUpdate:
			id := request.Data.id
			storage[id] = request.Data.value
			response.value = storage[id]
			response.ok = true
		case get:
			id := request.Data.id
			response.value, response.ok = storage[id]
		case remove:
			id := request.Data.id
			delete(storage, id)
		case length:
			response.len = len(storage)
		}
		request.SendBack(response)
	}
}
