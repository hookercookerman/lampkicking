package persistar

type KeyValueAdaptor interface {
	Set(collection string, key string, value []byte) (bool, error)
	Get(collection string, key string) ([]byte, error)
}

type GraphAdaptor interface {
	AddRelation(collection string, key string, relation string, relatedCollection string, relatedKey string) (bool, error)
	GetRelation(collection string, key string, relation string) ([]string, [][]byte, error)
}

type Adaptor interface {
	GraphAdaptor
	KeyValueAdaptor
}
