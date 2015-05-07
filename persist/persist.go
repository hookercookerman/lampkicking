package persist

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var InvalidJSONError = errors.New("Invalid JSON")

type Backend interface {
	Call(method string, path string, paramas *url.Values, data []byte) ([]byte, error)
}

type DefaultBackend struct {
	url        string
	httpClient *http.Client
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func (backend *DefaultBackend) Call(method string, path string, params *url.Values, data []byte) ([]byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = backend.url + path

	if params != nil && len(*params) > 0 {
		path += "?" + params.Encode()
	}

	if data != nil && !isJSON(data) {
		return nil, InvalidJSONError
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	res, err := backend.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 500 {
		return nil, errors.New(extractErrorText(resBody))
	} else {
		return resBody, nil
	}
}

func extractErrorText(err []byte) string {
	var jsonError map[string]string
	json.Unmarshal(err, &jsonError)
	return jsonError["error"]
}

func NewDefaultBackend(url string, httpClient *http.Client) *DefaultBackend {
	return &DefaultBackend{
		httpClient: httpClient,
		url:        url,
	}
}
