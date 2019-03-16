package static_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"wallawire/web/assets"
	"wallawire/web/static"
)

const (
	hAcceptRanges    = "Accept-Ranges"
	hContentLength   = "Content-Length"
	hContentType     = "Content-Type"
	hDate            = "Date"
	hEtag            = "Etag"
	hIfModifiedSince = "If-Modified-Since"
	hLastModified    = "Last-Modified"
	hLocation        = "Location"
	mimeTypeHtml     = "text/html; charset=utf-8"
	mimeTypeText     = "text/plain; charset=utf-8"
)

func TestStatic(t *testing.T) {

	const IgnoreValue = "XXX"

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
			Alias:          "index page",
			Path:           "/",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "bytes",
				hContentLength: strconv.Itoa(len(getFileBody("/index.html"))),
				hContentType:   mimeTypeHtml,
				hDate:          IgnoreValue,
				hEtag:          "",
				hLastModified:  getFileTime("/index.html"),
			},
			ResponseBody: getFileBody("/index.html"),
		},
		{
			Alias:         "index page with conditional get",
			Path:          "/",
			RequestMethod: http.MethodGet,
			RequestHeaders: map[string]string{
				hIfModifiedSince: getFileTime("/index.html"),
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusNotModified,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "",
				hContentType:   "",
				hContentLength: "",
				hDate:          IgnoreValue,
				hEtag:          "",
				hLastModified:  getFileTime("/index.html"),
			},
			ResponseBody: []byte{},
		},
		{
			Alias:          "return index for html5 route",
			Path:           "/about",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "bytes",
				hContentLength: strconv.Itoa(len(getFileBody("/index.html"))),
				hContentType:   mimeTypeHtml,
				hDate:          IgnoreValue,
				hEtag:          "",
				hLastModified:  getFileTime("/index.html"),
			},
			ResponseBody: getFileBody("/index.html"),
		},
		{
			Alias:          "index.html redirects to /",
			Path:           "/index.html",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusMovedPermanently,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "",
				hContentLength: "36",
				hContentType:   mimeTypeHtml,
				hDate:          IgnoreValue,
				hLocation:      "/",
			},
			ResponseBody: []byte("<a href=\"/\">Moved Permanently</a>.\n\n"),
		},
		{
			Alias:          "non-root static files are reachable",
			Path:           "/fonts/roboto-latin-100.woff",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "bytes",
				hContentLength: strconv.Itoa(len(getFileBody("/fonts/roboto-latin-100.woff"))),
				hContentType:   "font/woff",
				hDate:          IgnoreValue,
				hEtag:          "",
				hLastModified:  getFileTime("/fonts/roboto-latin-100.woff"),
			},
			ResponseBody: nil,
		},
		{
			Alias:          "static directory returns root",
			Path:           "/fonts",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "bytes",
				hContentLength: strconv.Itoa(len(getFileBody("/index.html"))),
				hContentType:   mimeTypeHtml,
				hDate:          IgnoreValue,
				hEtag:          "",
				hLastModified:  getFileTime("/index.html"),
			},
			ResponseBody: getFileBody("/index.html"),
		},
		{
			Alias:          "static directory with trailing slash returns root",
			Path:           "/fonts/",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "bytes",
				hContentLength: strconv.Itoa(len(getFileBody("/index.html"))),
				hContentType:   mimeTypeHtml,
				hDate:          IgnoreValue,
				hEtag:          "",
				hLastModified:  getFileTime("/index.html"),
			},
			ResponseBody: getFileBody("/index.html"),
		},
		{
			Alias:          "index page POST",
			Path:           "/",
			RequestMethod:  http.MethodPost,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusMethodNotAllowed,
			ResponseHeaders: map[string]string{
				hAcceptRanges:  "",
				hContentLength: "19",
				hContentType:   mimeTypeText,
				hDate:          IgnoreValue,
			},
			ResponseBody: []byte("Method Not Allowed\n"),
		},
	}

	newReader := func(b []byte) io.Reader {
		if b == nil {
			return nil
		}
		return bytes.NewReader(b)
	}

	for _, testCase := range testCases {

		testFn := func(tt *testing.T) {

			embeddedStore, _ := static.NewEmbeddedStore()
			server := httptest.NewServer(static.Handler(embeddedStore))
			defer server.Close()

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			req, err := http.NewRequest(testCase.RequestMethod, server.URL+testCase.Path, newReader(testCase.RequestBody))
			if err != nil {
				tt.Fatalf("Cannot create request: %s", err.Error())
			}
			for key, value := range testCase.RequestHeaders {
				req.Header.Add(key, value)
			}

			rsp, errRsp := client.Do(req)
			if errRsp != nil {
				tt.Fatalf("Error getting response: %s", errRsp.Error())
			}

			body, errBody := ioutil.ReadAll(rsp.Body)
			if errBody != nil {
				tt.Fatalf("Error reading response: %s", errBody.Error())
			}
			defer rsp.Body.Close()

			if got, want := rsp.StatusCode, testCase.ResponseStatus; got != want {
				tt.Errorf("Bad status: %d, expected: %d", got, want)
			}

			// test that expected headers are present
			// that headers are not present (empty string)
			// that headers are present but do not check value (IgnoreValue)
			for key, value := range testCase.ResponseHeaders {
				if got, want := rsp.Header.Get(key), value; got != want && want != IgnoreValue {
					tt.Errorf("Bad response header %s: %s, expected %s", key, got, want)
				}
			}

			// test that no unexpected headers are present
			for key := range rsp.Header {
				if _, ok := testCase.ResponseHeaders[key]; !ok {
					tt.Errorf("Unexpected response header %s", key)
				}
			}

			if testCase.ResponseBody != nil {
				if bytes.Compare(body, testCase.ResponseBody) != 0 {
					tt.Errorf("Bad body: %s, expected %s", body, testCase.ResponseBody)
				}
			}

		} // fn

		t.Run(testCase.Alias, testFn)

	} // cases

}

func getFileBody(path string) []byte {
	return assets.Asset(path)
}

func getFileTime(path string) string {
	return assets.AssetFile(path).ModTime().UTC().Truncate(time.Second).Format(http.TimeFormat)
}
