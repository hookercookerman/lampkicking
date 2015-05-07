package keyvalue

import (
	"fmt"

	"github.com/hookercookerman/lampkicking/persist"
)

type Client struct {
	B persist.Backend
}

func path(collection, key string) string {
	return fmt.Sprintf("%v/%v", collection, key)
}

func (client *Client) Set(collection string, key string, data []byte) ([]byte, error) {
	return client.B.Call("POST", path(collection, key), nil, data)
}

func (client *Client) Get(collection, key string) ([]byte, error) {
	return client.B.Call("GET", path(collection, key), nil, nil)
}
