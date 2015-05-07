package users

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/hookercookerman/lampkicking/persist"
	"github.com/hookercookerman/lampkicking/persist/client"
	"github.com/nu7hatch/gouuid"
)

const (
	RelationName   string = "connections"
	CollectionName string = "users"
)

type UserService struct {
	PersistAPI *client.API
}

func NewUserService(persistUrl string, backend persist.Backend) *UserService {
	api := &client.API{}
	api.Init(persistUrl, backend)
	return &UserService{PersistAPI: api}
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	output := `{"error" : "%v"}`
	http.Error(w, fmt.Sprintf(output, err.Error()), 500)
}

func (service *UserService) GetConnections(w http.ResponseWriter, req *http.Request) {
	key := bone.GetValue(req, "id")
	result, err := service.PersistAPI.Graph.GetRelation(CollectionName, key, RelationName)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		writeErrorResponse(err, w)
	} else {
		w.Write(result)
	}
}

func (service *UserService) AddConnection(w http.ResponseWriter, req *http.Request) {
	id := bone.GetValue(req, "id")
	connectionId := bone.GetValue(req, "connection_id")
	result, err := service.PersistAPI.Graph.AddRelation(CollectionName, id, CollectionName, connectionId, RelationName)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		writeErrorResponse(err, w)
	} else {
		w.WriteHeader(201)
		w.Write(result)
	}
}

func (service *UserService) Store(w http.ResponseWriter, req *http.Request) {
	key, _ := uuid.NewV4()
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	_, err := service.PersistAPI.KeyValue.Set(CollectionName, key.String(), body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		writeErrorResponse(err, w)
	} else {
		w.WriteHeader(201)
		w.Write([]byte(fmt.Sprintf(`{"id" : "%v"}`, key.String())))
	}
}
