package graph

import (
	"fmt"

	"github.com/hookercookerman/lampkicking/persist"
)

type Client struct {
	B persist.Backend
}

func (client *Client) AddRelation(collection, key, relatedCollection, relatedKey, relation string) ([]byte, error) {
	path := fmt.Sprintf("/%v/%v/relation/%v/%v/%v", collection, key, relation, relatedCollection, relatedKey)
	return client.B.Call("PUT", path, nil, nil)
}

func (client *Client) GetRelation(collection, key, relation string) ([]byte, error) {
	path := fmt.Sprintf("/%v/%v/relation/%v", collection, key, relation)
	return client.B.Call("GET", path, nil, nil)
}
