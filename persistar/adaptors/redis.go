package adaptors

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var KeyMissingError = errors.New("Key Not Found")

type RedisAdaptor struct {
	pool *redis.Pool
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
				return c, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewRedisAdaptor(server string, password string) (*RedisAdaptor, error) {
	if server == "" {
		server = "localhost:6379"
	}
	pool := newPool(server, password) // surely need to generate error here as well
	return &RedisAdaptor{pool: pool}, nil
}

func fullKey(collection, key string) string {
	return fmt.Sprintf("%v:%v", collection, key)
}

func (r *RedisAdaptor) Set(collection string, key string, value []byte) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", fullKey(collection, key), value); err != nil {
		return false, err
	}
	return true, nil
}

func relationKey(collection, key, relation string) string {
	return fmt.Sprintf("%v:%v:%v", collection, key, relation)
}

func (r *RedisAdaptor) GetRelation(collection, key, relation string) ([]string, [][]byte, error) {
	conn := r.pool.Get()
	members, err := redis.Strings(conn.Do("SMEMBERS", relationKey(collection, key, relation)))

	if err != nil {
		return nil, nil, err
	}
	if len(members) == 0 {
		return nil, nil, errors.New("No Members")
	}

	var argList []interface{}
	for _, member := range members {
		argList = append(argList, member)
	}

	reply, err := redis.Strings(conn.Do("MGET", argList...))
	if err != nil {
		return nil, nil, err
	}

	var result [][]byte
	for _, o := range reply {
		result = append(result, []byte(o))
	}

	return members, result, nil
}

func (r *RedisAdaptor) exists(conn redis.Conn, collection, key string) (bool, error) {
	return redis.Bool(conn.Do("EXISTS", fullKey(collection, key)))
}

func (r *RedisAdaptor) AddRelation(collection, key, relation, relatedCollection, relatedKey string) (bool, error) {
	conn := r.pool.Get()
	if ok, _ := r.exists(conn, collection, key); ok == false {
		return false, KeyMissingError
	}
	if ok, _ := r.exists(conn, relatedCollection, relatedKey); ok == false {
		return false, KeyMissingError
	}
	_, err := conn.Do("SADD", relationKey(collection, key, relation), fullKey(relatedCollection, relatedKey))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisAdaptor) Get(collection string, key string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	redisKey := fullKey(collection, key)
	value, err := redis.Bytes(conn.Do("GET", redisKey))
	if err != nil {
		return nil, err
	}
	return value, nil
}
