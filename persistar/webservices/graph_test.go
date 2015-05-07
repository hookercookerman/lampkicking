package webservices_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/go-zoo/bone"
	"github.com/hookercookerman/lampkicking/persistar/adaptors"
	"github.com/hookercookerman/lampkicking/persistar/interactors"
	. "github.com/hookercookerman/lampkicking/persistar/webservices"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keyvalue", func() {
	var graphService *GraphService
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		redisAdaptor, _ := adaptors.NewRedisAdaptor(os.Getenv("REDIS_SERVER"), "")
		interactor := &interactors.GraphInteractor{
			Adaptor: redisAdaptor,
		}
		graphService = &GraphService{
			Interactor: interactor,
		}
	})

	Describe("GetRelations", func() {
		Context("when relation has members", func() {
			BeforeEach(func() {
				Conn.Do("SET", "test1:egg", `{"name" : "beans"}`)
				Conn.Do("SET", "test2:beans2", `{"name" : "egg"}`)
				Conn.Do("SADD", "test1:egg:friends", "test2:beans2")
				payload := []byte{}
				req, _ := http.NewRequest("GET", "/test1/egg/relation/friends", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				mux := bone.New()
				mux.Get("/:collection/:key/relation/:relation", http.HandlerFunc(graphService.GetRelation))
				mux.ServeHTTP(recorder, req)
			})

			It("responds with result of relations", func() {
				result, _ := ioutil.ReadAll(recorder.Body)
				Expect(string(result)).To(MatchJSON(`{"count": 1, "results" : [{"path": "test2:beans2", "value" : {"name" : "egg"}}] }`))
			})
		})
	})

	Describe("AddRelation", func() {
		Context("when successful", func() {
			BeforeEach(func() {
				Conn.Do("SET", "test1:egg", `{"name" : "beans"}`)
				Conn.Do("SET", "test2:beans", `{"name" : "egg"}`)
				payload := []byte{}
				req, _ := http.NewRequest("PUT", "/test1/egg/relation/friends/test2/beans", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				mux := bone.New()
				mux.Put("/:collection/:key/relation/:relation/:relatedCollection/:relatedKey", http.HandlerFunc(graphService.AddRelation))
				mux.ServeHTTP(recorder, req)
			})
			It("returns 201", func() {
				Expect(recorder.Code).To(Equal(201))
			})
		})
		Context("when unsuccessful, missing keys", func() {
			BeforeEach(func() {
				payload := []byte{}
				req, _ := http.NewRequest("PUT", "/test1/egg1/relation/friends/test2/beans2", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				mux := bone.New()
				mux.Put("/:collection/:key/relation/:relation/:relatedCollection/:relatedKey", http.HandlerFunc(graphService.AddRelation))
				mux.ServeHTTP(recorder, req)
			})
			It("returns 500", func() {
				Expect(recorder.Code).To(Equal(500))
			})
		})
	})
})
