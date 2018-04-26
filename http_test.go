package sweetest

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTP(t *testing.T) {

	tt := &testTest{t: t}
	e := NewExpect(t)

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

			e.Expect(req.URL.Path).
				Name("path").
				Equal("/path")
			e.Expect(req.URL.Query().Get("query1")).
				Name("query1").
				Equal("query-val-2")
			e.Expect(req.URL.Query().Get("query2")).
				Name("query2").
				Equal("v2")
			e.Expect(req.Header.Get("hdr1")).
				Name("hdr1").
				Equal("hdr1-2")
			e.Expect(req.Header.Get("hdr2")).
				Name("hdr2").
				Equal("hdr2")
			e.Expect(req.Method).Name("method").Equal("PUT")
			body, _ := ioutil.ReadAll(req.Body)
			e.Expect(string(body)).Equal(`{"key":"val"}`)

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
			e.Expect(b.Resp).Name("body.val").Equal("OK")
		})

	tt.expectNothing()

	// Bad Expectations

	hasRun.Status(200)
	tt.expectError("expect status 400 to equal 200")

	hasRun.Header("missing", "val")
	tt.expectError(`expect header missing "" to equal "val"`)

	/*
		BodyJSON(func(b struct {
			Resp string `json:"resp"`
		}) {
			e.Expect(b.Resp).Equal("val")
		})
	*/

}
