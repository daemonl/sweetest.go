package sweetest

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTP(t *testing.T) {

	tt := &testTest{t: t}

	req := BuildRequest().
		Method("PUT").
		Path("/path").
		Query("query1", "query-val-%d", 1).
		Query("query1", "query-val-%d", 2).
		Query("query2", "v2").
		Header("Hdr1", "hdr1-%d", 1).
		Header("Hdr2", "hdr2").
		Header("Hdr1", "hdr1-%d", 2).
		BodyJSON(map[string]interface{}{
			"key": "val",
		})

	hasRun := req.Run(tt, http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {

			if path := req.URL.Path; path != "/path" {
				t.Errorf("Path: %s", path)
			}
			if got := req.URL.Query().Get("query1"); got != "query-val-2" {
				t.Errorf("Query1: %s", got)
			}
			if got := req.URL.Query().Get("query2"); got != "v2" {
				t.Errorf("Query2: %s", got)
			}
			if got := req.Header.Get("hdr1"); got != "hdr1-2" {
				t.Errorf("hdr1: %s", got)
			}
			if got := req.Header.Get("hdr2"); got != "hdr2" {
				t.Errorf("hdr2: %s", got)
			}
			if req.Method != "PUT" {
				t.Errorf("Method: %s", req.Method)
			}
			body, _ := ioutil.ReadAll(req.Body)
			if got := string(body); got != `{"key":"val"}` {
				t.Errorf("Bad Body: %s", got)
			}

			rw.Header().Add("RespHdr", "val")
			rw.Header().Add("OtherHeader", "other")
			rw.WriteHeader(400)
			rw.Write([]byte(`{"resp":"OK"}`))
		}))

	// Good Expectations
	hasRun.
		Status(400).
		Header("resphdr", "val").
		BodyJSON(func(b struct {
			Resp string `json:"resp"`
		}) {
			if b.Resp != "OK" {
				t.Errorf("Expected OK, got %s", b.Resp)
			}
		})

	tt.expectNothing()

	// Bad Expectations

	hasRun.Status(200)
	tt.expectFatal("Status 400, expected 200")

	hasRun.Header("missing", "val")
	tt.expectError(`Expect header missing to be "val", got ""`)

	/*
		BodyJSON(func(b struct {
			Resp string `json:"resp"`
		}) {
			e.Expect(b.Resp).Equal("val")
		})
	*/

}
