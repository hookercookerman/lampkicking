package client

import (
	"net/http"

	"github.com/hookercookerman/lampkicking/persist"
	"github.com/hookercookerman/lampkicking/persist/graph"
	"github.com/hookercookerman/lampkicking/persist/keyvalue"
)

type API struct {
	backend  persist.Backend
	url      string
	KeyValue *keyvalue.Client
	Graph    *graph.Client
}

func (api *API) Init(url string, backend persist.Backend) {
	if backend == nil {
		backend = persist.NewDefaultBackend(url, http.DefaultClient)
	}
	api.backend = backend
	api.KeyValue = &keyvalue.Client{B: backend}
	api.Graph = &graph.Client{B: backend}
}
