package auth_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallawire/web/auth"

	"github.com/go-chi/chi"
)

func TestLogout(t *testing.T) {

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
			Alias:          "logout success POST",
			Path:           "/logout",
			RequestMethod:  http.MethodPost,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hContentLength: "3",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
				hSetCookie:     "jwt=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Secure",
			},
			ResponseBody: []byte("OK\n"),
		},
		{
			Alias:          "logout success GET",
			Path:           "/logout",
			RequestMethod:  http.MethodGet,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hContentLength: "3",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
				hSetCookie:     "jwt=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Secure",
			},
			ResponseBody: []byte("OK\n"),
		},
		{
			Alias:          "logout failure DELETE",
			Path:           "/logout",
			RequestMethod:  http.MethodDelete,
			RequestHeaders: nil,
			RequestBody:    nil,
			ResponseStatus: http.StatusMethodNotAllowed,
			ResponseHeaders: map[string]string{
				hContentLength: "19",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
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

			handler := chi.NewRouter()
			handler.Get("/logout", auth.Logout())
			handler.Post("/logout", auth.Logout())
			handler.MethodNotAllowed(sendMessageHandler(http.StatusMethodNotAllowed))

			server := httptest.NewServer(handler)
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
				tt.Fatalf("UserError getting response: %s", errRsp.Error())
			}

			body, errBody := ioutil.ReadAll(rsp.Body)
			if errBody != nil {
				tt.Fatalf("UserError reading response: %s", errBody.Error())
			}
			defer rsp.Body.Close()

			if got, want := rsp.StatusCode, testCase.ResponseStatus; got != want {
				tt.Errorf("Bad status: %d, expected: %d", got, want)
			}

			// test that expected headers are present
			// that headers are not present (empty string)
			// that headers are present but do not check value (ignoreValue)
			for key, value := range testCase.ResponseHeaders {
				if got, want := rsp.Header.Get(key), value; got != want && want != ignoreValue {
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
