package webservices

import (
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/hookercookerman/lampkicking/persistar/interactors"
)

type GraphService struct {
	Interactor *interactors.GraphInteractor
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	output := `{"error" : "%v" }`
	http.Error(w, fmt.Sprintf(output, err.Error()), 500)
}

func (service *GraphService) GetRelation(w http.ResponseWriter, request *http.Request) {
	collection := bone.GetValue(request, "collection")
	key := bone.GetValue(request, "key")
	relation := bone.GetValue(request, "relation")
	result, err := service.Interactor.GetRelation(collection, key, relation)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		writeErrorResponse(err, w)
	} else {
		w.Write(result)
	}
}

func (service *GraphService) AddRelation(w http.ResponseWriter, request *http.Request) {
	collection := bone.GetValue(request, "collection")
	key := bone.GetValue(request, "key")
	relatedCollection := bone.GetValue(request, "relatedCollection")
	relatedKey := bone.GetValue(request, "relatedKey")
	relation := bone.GetValue(request, "relation")
	w.Header().Set("Content-Type", "application/json")
	if _, err := service.Interactor.AddRelation(collection, key, relation, relatedCollection, relatedKey); err != nil {
		writeErrorResponse(err, w)
	} else {
		w.WriteHeader(201)
	}
}
