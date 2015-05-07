package users_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	. "github.com/hookercookerman/lampkicking/users"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockBackend struct {
	Calls       []*Call
	ReturnError error
	ReturnValue []byte
}
type Call struct {
	Method string
	Path   string
	Params *url.Values
	Data   []byte
}

func (backend *MockBackend) storeCall(method, path string, params *url.Values, data []byte) {
	backend.Calls = append(backend.Calls, &Call{
		Method: method,
		Path:   path,
		Params: params,
		Data:   data,
	})
}

func (backend *MockBackend) Call(method, path string, params *url.Values, data []byte) ([]byte, error) {
	backend.storeCall(method, path, params, data)
	if backend.ReturnError != nil {
		return nil, backend.ReturnError
	}
	if backend.ReturnValue != nil {
		return backend.ReturnValue, nil
	}
	return nil, nil
}

var _ = Describe("Webservice", func() {
	var mockBackend *MockBackend
	var userService *UserService
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		mockBackend = &MockBackend{}
		userService = NewUserService("localhost:9001", mockBackend)
	})

	Describe("GetConnection", func() {
		Context("when successful", func() {
			BeforeEach(func() {
				mockBackend.ReturnValue = []byte(`{"count" : 1, "results" : [{"name" : "beans"}]}`)
				payload := []byte{}
				req, _ := http.NewRequest("GET", "", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(userService.GetConnections).ServeHTTP(recorder, req)
			})
			It("returns 200", func() {
				Expect(recorder.Code).To(Equal(200))
			})
			It("returns json result", func() {
				result, _ := ioutil.ReadAll(recorder.Body)
				Expect(string(result)).To(MatchJSON(`{"count" : 1, "results" : [{"name" : "beans"}]}`))
			})
		})
		Context("when error returned from persist service", func() {
			BeforeEach(func() {
				mockBackend.ReturnError = errors.New("test")
				req, _ := http.NewRequest("GET", "", bytes.NewReader([]byte{}))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(userService.GetConnections).ServeHTTP(recorder, req)
			})
			It("returns 500 http code", func() {
				Expect(recorder.Code).To(Equal(500))
			})
			It("returns error response", func() {
				result, _ := ioutil.ReadAll(recorder.Body)
				Expect(string(result)).To(MatchJSON(`{"error" : "test"}`))
			})
		})
	})

	Describe("AddConnection", func() {
		Context("when successful", func() {
			BeforeEach(func() {
				mockBackend.ReturnValue = []byte("")
				payload := []byte{}
				req, _ := http.NewRequest("POST", "", bytes.NewReader(payload))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(userService.AddConnection).ServeHTTP(recorder, req)
			})
			It("returns 201", func() {
				Expect(recorder.Code).To(Equal(201))
			})
		})
		Context("when error returned from persist service", func() {
			BeforeEach(func() {
				mockBackend.ReturnError = errors.New("test")
				req, _ := http.NewRequest("POST", "", bytes.NewReader([]byte{}))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(userService.Store).ServeHTTP(recorder, req)
			})
			It("returns 500 http code", func() {
				Expect(recorder.Code).To(Equal(500))
			})
			It("returns error response", func() {
				result, _ := ioutil.ReadAll(recorder.Body)
				Expect(string(result)).To(MatchJSON(`{"error" : "test"}`))
			})
		})
	})
	Describe("Store", func() {
		Context("when successful", func() {
			BeforeEach(func() {
				mockBackend.ReturnValue = []byte("")
				payload := `{"name" : "lemon"}`
				req, _ := http.NewRequest("POST", "", bytes.NewReader([]byte(payload)))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(userService.Store).ServeHTTP(recorder, req)
			})
			It("returns 201", func() {
				Expect(recorder.Code).To(Equal(201))
			})
			// @todo mock out uuid generator make me a prop on service
			It("returns uuid of new user", func() {
				result, _ := ioutil.ReadAll(recorder.Body)
				Expect(string(result)).ToNot(BeNil())
			})
		})
		Context("when error returned from persist service", func() {
			BeforeEach(func() {
				mockBackend.ReturnError = errors.New("test")
				req, _ := http.NewRequest("POST", "", bytes.NewReader([]byte{}))
				recorder = httptest.NewRecorder()
				http.HandlerFunc(userService.Store).ServeHTTP(recorder, req)
			})
			It("returns a 500 http code", func() {
				Expect(recorder.Code).To(Equal(500))
			})
			It("returns error response", func() {
				result, _ := ioutil.ReadAll(recorder.Body)
				Expect(string(result)).To(MatchJSON(`{"error" : "test"}`))
			})
		})
	})
})
