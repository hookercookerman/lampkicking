package adaptors_test

import (
	"os"

	"github.com/garyburd/redigo/redis"
	. "github.com/hookercookerman/lampkicking/persistar/adaptors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redis", func() {
	var redisAdaptor *RedisAdaptor

	BeforeEach(func() {
		redisAdaptor, _ = NewRedisAdaptor(os.Getenv("REDIS_SERVER"), "")
	})

	Describe("GET", func() {
		It("should return error when key not found", func() {
			_, err := redisAdaptor.Get("test", "wrong")
			Expect(err).ToNot(BeNil())
		})

		Context("When key exists", func() {
			BeforeEach(func() {
				Conn.Do("SET", "test:test", "bar")
			})
			It("should return the value", func() {
				result, _ := redisAdaptor.Get("test", "test")
				Expect(string(result)).To(Equal("bar"))
			})
		})
	})

	Describe("SET", func() {
		It("should store value", func() {
			testValue := []byte("testing")
			redisAdaptor.Set("test1", "test1", testValue)
			result, _ := redis.String(Conn.Do("GET", "test1:test1"))
			Expect(string(result)).To(Equal("testing"))
		})
	})

	Describe("GetRelation", func() {
		Context("when relation does not exist", func() {
			It("should return an error", func() {
				_, _, err := redisAdaptor.GetRelation("badness", "whatever", "egg")
				Expect(err).ToNot(BeNil())
			})
		})
		Context("when relation does exist", func() {
			BeforeEach(func() {
				Conn.Do("SET", "r1:key1", "egg")
				Conn.Do("SET", "r2:key2", "beans")
				redisAdaptor.AddRelation("r1", "key1", "connections", "r2", "key2")
			})
			It("should return the members", func() {
				members, _, _ := redisAdaptor.GetRelation("r1", "key1", "connections")
				Expect(members).To(Equal([]string{"r2:key2"}))
			})
			It("should return the value of the members", func() {
				_, values, _ := redisAdaptor.GetRelation("r1", "key1", "connections")
				Expect(values).To(Equal([][]byte{[]byte("beans")}))
			})
		})
	})

	Describe("AddRelation", func() {
		Context("when keys does not exist", func() {
			It("should return an error", func() {
				_, err := redisAdaptor.AddRelation("missing1", "key1", "connections", "missing2", "key2")
				Expect(err).To(Equal(KeyMissingError))
			})
		})

		Context("when keys exist", func() {
			BeforeEach(func() {
				Conn.Do("SET", "c1:key1", "egg")
				Conn.Do("SET", "c2:key2", "beans")
			})

			It("should add member to set", func() {
				_, err := redisAdaptor.AddRelation("c1", "key1", "connections", "c2", "key2")
				Expect(err).To(BeNil())
				result, _ := redis.Strings(Conn.Do("SMEMBERS", "c1:key1:connections"))
				Expect(result).To(Equal([]string{"c2:key2"}))
			})
		})
	})
})
