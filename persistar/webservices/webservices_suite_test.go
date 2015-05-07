package webservices_test

import (
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/hookercookerman/lampkicking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stvp/tempredis"

	"testing"
)

var redisServer *tempredis.Server

func init() {
	os.Setenv("GO_ENV", "test")
}

func TestWebservices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webservices Suite")
}

var Conn redis.Conn

var _ = BeforeSuite(func() {
	var err error
	redisServer, err = tempredis.Start(
		tempredis.Config{
			"port":      lampkicking.Getenv("REDIS_PORT"),
			"databases": "1",
		},
	)
	Conn, err = redis.Dial("tcp", ":"+lampkicking.Getenv("REDIS_PORT"))
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	Conn.Do("FLUSHDB")
	Conn.Close()
	redisServer.Term()
})
