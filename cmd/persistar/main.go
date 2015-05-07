package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/hookercookerman/lampkicking"
	"github.com/hookercookerman/lampkicking/persistar/adaptors"
	"github.com/hookercookerman/lampkicking/persistar/interactors"
	"github.com/hookercookerman/lampkicking/persistar/webservices"
)

func main() {
	redisAdaptor, _ := adaptors.NewRedisAdaptor(lampkicking.Getenv("REDIS_SERVER"), lampkicking.Getenv("REDIS_PASSWORD"))

	keyValueWebService := webservices.KeyValueService{
		Interactor: &interactors.KeyValueInteractor{
			Adaptor: redisAdaptor,
		},
	}
	graphWebService := webservices.GraphService{
		Interactor: &interactors.GraphInteractor{
			Adaptor: redisAdaptor,
		},
	}
	mux := bone.New()
	mux.Put("/:collection/:key/relation/:relation/:relatedCollection/:relatedKey", http.HandlerFunc(graphWebService.AddRelation))
	mux.Get("/:collection/:key/relation/:relation", http.HandlerFunc(graphWebService.GetRelation))
	mux.Post("/:collection/:key", http.HandlerFunc(keyValueWebService.Set))

	http.ListenAndServe(":"+lampkicking.Getenv("PERSISTAR_PORT"), mux)
}
