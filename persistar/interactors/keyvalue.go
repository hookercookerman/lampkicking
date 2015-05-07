package interactors

import (
	"encoding/json"
	"errors"

	"github.com/hookercookerman/lampkicking/persistar"
)

type KeyValueInteractor struct {
	Adaptor persistar.KeyValueAdaptor
}

func (interactor *KeyValueInteractor) Get(collection string, key string) ([]byte, error) {
	return interactor.Adaptor.Get(collection, key)
}

func (interactor *KeyValueInteractor) Set(collection string, key string, value []byte) (bool, error) {
	var parsedValue map[string]interface{}
	if err := json.Unmarshal(value, &parsedValue); err != nil {
		return false, errors.New("Invalid json")
	}
	return interactor.Adaptor.Set(collection, key, value)
}
