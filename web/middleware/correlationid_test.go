package middleware_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
	"wallawire/ctxutil"
	"wallawire/web/middleware"
)

const (
	ignoreValue    = "XXX"
	hContentLength = "Content-Length"
	hContentType   = "Content-Type"
	hDate          = "Date"
	mimeTypeText   = "text/plain; charset=utf-8"
)

func printCorrelationIdHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cID := ctxutil.CorrelationIDFromContext(r.Context())

		msg := []byte(cID)
		w.Header().Set(hContentLength, strconv.Itoa(len(msg)))
		w.Write(msg)

	})
}

func TestCorrelationID(b *testing.T) {

	testCases := []struct {
		Alias           string
		Path            string
		RequestMethod   string
		RequestHeaders  map[string]string
		RequestBody     []byte
		ResponseStatus  int
		ResponseHeaders map[string]string
		ResponseBody    []byte
	}{
		{
			Alias:          "success",
			Path:           "/",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hContentLength: "36",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
			},
			ResponseBody: nil,
		},
	}

	newReader := func(b []byte) io.Reader {
		if b == nil {
			return nil
		}
		return bytes.NewReader(b)
	}

	for _, testCase := range testCases {

		testFn := func(t *testing.T) {

			idgen, errGen := ctxutil.NewIdGenerator()
			if errGen != nil {
				t.Fatal(errGen)
			}

			handler := chi.NewRouter()
			handler.Use(middleware.CorrelationID(idgen))
			handler.Get("/", printCorrelationIdHandler())

			server := httptest.NewServer(handler)
			defer server.Close()

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			req, err := http.NewRequest(testCase.RequestMethod, server.URL+testCase.Path, newReader(testCase.RequestBody))
			if err != nil {
				t.Fatalf("Cannot create request: %s", err.Error())
			}
			for key, value := range testCase.RequestHeaders {
				req.Header.Add(key, value)
			}

			rsp, errRsp := client.Do(req)
			if errRsp != nil {
				t.Fatalf("UserError getting response: %s", errRsp.Error())
			}

			body, errBody := ioutil.ReadAll(rsp.Body)
			if errBody != nil {
				t.Fatalf("UserError reading response: %s", errBody.Error())
			}
			defer rsp.Body.Close()

			if got, want := rsp.StatusCode, testCase.ResponseStatus; got != want {
				t.Errorf("Bad status: %d, expected: %d", got, want)
			}

			// test that expected headers are present
			// that headers are not present (empty string)
			// that headers are present but do not check value (ignoreValue)
			for key, value := range testCase.ResponseHeaders {
				if got, want := rsp.Header.Get(key), value; got != want && want != ignoreValue {
					t.Errorf("Bad response header %s: %s, expected %s", key, got, want)
				}
			}

			// test that no unexpected headers are present
			for key := range rsp.Header {
				if _, ok := testCase.ResponseHeaders[key]; !ok {
					t.Errorf("Unexpected response header %s", key)
				}
			}

			t.Log(string(body))

			if testCase.ResponseBody != nil {
				if bytes.Compare(body, testCase.ResponseBody) != 0 {
					t.Errorf("Bad body: %s, expected %s", body, testCase.ResponseBody)
				}
			}

		} // fn

		b.Run(testCase.Alias, testFn)

	} // cases

}
