package auth_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"

	"wallawire/model"
	"wallawire/web/auth"
)

func TestAuthenticator(b *testing.T) {

	now := time.Now().Truncate(time.Second)

	demouser := &model.SessionToken{
		SessionID: "S123",
		ID:        "id",
		Username:  "demouser",
		Name:      "Demo User",
		Roles:     []string{"users", "guests"},
		Issued:    now.Truncate(time.Minute),
		Expires:   now.Truncate(time.Minute).Add(model.LoginTimeout),
	}

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
			Alias:         "success",
			Path:          "/",
			RequestMethod: http.MethodGet,
			RequestHeaders: map[string]string{
				hCookie: getCookieString(demouser, testPassword),
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hContentLength: "98",
				hContentType:   mimeTypeJson,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte(`{"sessionID":"S123","id":"id","username":"demouser","name":"Demo User","roles":["users","guests"]}`),
		},
		{
			Alias:          "missing cookie",
			Path:           "/",
			RequestMethod:  http.MethodGet,
			RequestHeaders: map[string]string{},
			RequestBody:    nil,
			ResponseStatus: http.StatusUnauthorized,
			ResponseHeaders: map[string]string{
				hContentLength:           "13",
				hContentType:             mimeTypeText,
				hDate:                    ignoreValue,
				"X-Content-Type-Options": ignoreValue,
			},
			ResponseBody: []byte("Unauthorized\n"),
		},
		{
			Alias:         "bad token password",
			Path:          "/",
			RequestMethod: http.MethodGet,
			RequestHeaders: map[string]string{
				hCookie: getCookieString(demouser, testPassword+"bogus"),
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusUnauthorized,
			ResponseHeaders: map[string]string{
				hContentLength:           "13",
				hContentType:             mimeTypeText,
				hDate:                    ignoreValue,
				"X-Content-Type-Options": ignoreValue,
			},
			ResponseBody: []byte("Unauthorized\n"),
		},
		{
			Alias:         "bogus cookie",
			Path:          "/",
			RequestMethod: http.MethodGet,
			RequestHeaders: map[string]string{
				hCookie: "jwt=foo",
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusUnauthorized,
			ResponseHeaders: map[string]string{
				hContentLength:           "13",
				hContentType:             mimeTypeText,
				hDate:                    ignoreValue,
				"X-Content-Type-Options": ignoreValue,
			},
			ResponseBody: []byte("Unauthorized\n"),
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

			handler := chi.NewRouter()
			handler.Use(auth.NewAuthenticator(testPassword)...)
			handler.Get("/", auth.Whoami())

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

			if testCase.ResponseBody != nil {
				if bytes.Compare(body, testCase.ResponseBody) != 0 {
					t.Errorf("Bad body: %s, expected %s", body, testCase.ResponseBody)
				}
			}

		} // fn

		b.Run(testCase.Alias, testFn)

	} // cases

}
