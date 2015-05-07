package webservices

import (
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/hookercookerman/lampkicking/persistar/interactors"
)

type KeyValueService struct {
	Interactor *interactors.KeyValueInteractor
}

func (service *KeyValueService) Set(w http.ResponseWriter, request *http.Request) {
	collection := bone.GetValue(request, "collection")
	key := bone.GetValue(request, "key")
	defer request.Body.Close()
	body, _ := ioutil.ReadAll(request.Body)
	w.Header().Set("Content-Type", "application/json")
	if _, err := service.Interactor.Set(collection, key, body); err != nil {
		writeErrorResponse(err, w)
	} else {
		w.WriteHeader(201)
	}
}
