package sweetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
)

type T interface {
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Logf(string, ...interface{})
	Helper()
}

// ModifierFunc allows custom extensions to the Request Builder
type ModifierFunc func(req *http.Request) *http.Request

// RequestBuilder builds a HTTP Request
type RequestBuilder interface {
	Method(string) RequestBuilder
	Path(string, ...interface{}) RequestBuilder
	Query(key string, val string, params ...interface{}) RequestBuilder
	Header(key string, val string, params ...interface{}) RequestBuilder
	BodyJSON(body interface{}) RequestBuilder
	BodyBytes(b []byte) RequestBuilder
	BodyReader(io.Reader) RequestBuilder
	With(ModifierFunc) RequestBuilder

	Run(T, http.Handler) HTTPResult
}

// BuildRequest returns a new default RequestBuilder
func BuildRequest() RequestBuilder {
	return &requestBuilder{
		req: httptest.NewRequest("GET", "/", nil),
	}
}

type requestBuilder struct {
	req *http.Request
}

func (rb *requestBuilder) With(fn ModifierFunc) RequestBuilder {
	rb.req = fn(rb.req)
	return rb
}

func (rb *requestBuilder) Method(m string) RequestBuilder {
	rb.req.Method = m
	return rb
}

func (rb *requestBuilder) Path(path string, params ...interface{}) RequestBuilder {
	rb.req.URL.Path = fmt.Sprintf(path, params...)
	return rb
}

func (rb *requestBuilder) Query(key string, val string, params ...interface{}) RequestBuilder {

	query, err := url.ParseQuery(rb.req.URL.RawQuery)
	if err != nil {
		panic(err.Error())
	}
	query.Set(key, fmt.Sprintf(val, params...))
	rb.req.URL.RawQuery = query.Encode()

	return rb
}
func (rb *requestBuilder) Header(key string, val string, params ...interface{}) RequestBuilder {
	rb.req.Header.Set(key, fmt.Sprintf(val, params...))
	return rb
}

func (rb *requestBuilder) BodyJSON(body interface{}) RequestBuilder {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err.Error())
	}
	return rb.BodyBytes(b)
}

func (rb *requestBuilder) BodyBytes(b []byte) RequestBuilder {
	return rb.BodyReader(bytes.NewBuffer(b))
}

func (rb *requestBuilder) BodyReader(r io.Reader) RequestBuilder {
	if closer, ok := r.(io.ReadCloser); ok {
		rb.req.Body = closer
	} else {
		rb.req.Body = ioutil.NopCloser(r)
	}
	return rb
}

func (rb *requestBuilder) Run(t T, handler http.Handler) HTTPResult {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, rb.req)
	return &httpResult{
		rw: rec,
		t:  t,
	}
}

type HTTPResult interface {
	Status(int) HTTPResult
	Header(key string, expect string, params ...interface{}) HTTPResult
	BodyJSON(callback interface{}) HTTPResult
	Raw(callback func(*httptest.ResponseRecorder)) HTTPResult
}

type httpResult struct {
	rw *httptest.ResponseRecorder
	t  T
}

func (tr *httpResult) Status(expect int) HTTPResult {
	tr.t.Helper()
	if tr.rw.Code != expect {
		tr.Log()
		tr.t.Fatalf("Status %d, expected %d", tr.rw.Code, expect)
	}
	return tr
}

func (tr *httpResult) Log() {
	tr.t.Helper()
	tr.t.Logf("Body: %s", tr.rw.Body.String())
}

func (tr *httpResult) Header(key string, expect string, params ...interface{}) HTTPResult {
	tr.t.Helper()
	if got := tr.rw.Header().Get(key); got != expect {
		tr.t.Errorf(`Expect header %s to be "%s", got "%s"`, key, expect, got)
	}
	return tr
}

func (tr *httpResult) BodyJSON(callback interface{}) HTTPResult {
	tr.t.Helper()
	callbackValue := reflect.ValueOf(callback)
	val := reflect.New(callbackValue.Type().In(0))
	valInterface := val.Interface()
	json.Unmarshal(tr.rw.Body.Bytes(), valInterface)
	callbackValue.Call([]reflect.Value{reflect.ValueOf(valInterface).Elem()})
	return tr
}

func (tr *httpResult) Raw(cb func(*httptest.ResponseRecorder)) HTTPResult {
	cb(tr.rw)
	return tr
}
