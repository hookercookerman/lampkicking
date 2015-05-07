package webservices_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/hookercookerman/lampkicking/persistar/adaptors"
	"github.com/hookercookerman/lampkicking/persistar/interactors"
	. "github.com/hookercookerman/lampkicking/persistar/webservices"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keyvalue", func() {
	var keyValueService *KeyValueService
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		redisAdaptor, _ := adaptors.NewRedisAdaptor(os.Getenv("REDIS_SERVER"), "")
		interactor := &interactors.KeyValueInteractor{
			Adaptor: redisAdaptor,
		}
		keyValueService = &KeyValueService{
			Interactor: interactor,
		}
	})

	Describe("SET", func() {
		Context("when successful", func() {
			BeforeEach(func() {
				payload := []byte(`{"name" : "egg"}`)
				req, _ := http.NewRequest("POST", "/test1/egg", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(keyValueService.Set).ServeHTTP(recorder, req)
			})
			It("returns 201", func() {
				Expect(recorder.Code).To(Equal(201))
			})
		})
		Context("when unsuccessful, invalid json", func() {
			BeforeEach(func() {
				payload := []byte(`{"name" "egg"}`)
				req, _ := http.NewRequest("POST", "/test1/egg", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(keyValueService.Set).ServeHTTP(recorder, req)
			})
			It("returns 500", func() {
				Expect(recorder.Code).To(Equal(500))
			})
		})
	})
})
