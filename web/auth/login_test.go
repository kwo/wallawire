package auth_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wallawire/model"
	"wallawire/web/auth"

	"github.com/go-chi/chi"
)

func TestLogin(t *testing.T) {

	now := time.Now().Truncate(time.Second)

	demouser := &model.User{
		ID:       "id",
		Username: "demouser",
		Name:     "Demo User",
		Created:  now,
		Updated:  now,
	}
	if err := demouser.SetPassword("demouser"); err != nil {
		t.Fatal(err)
	}

	demoRoles := []model.UserRole{
		{
			ID:   "roleid",
			Name: "users",
		},
	}

	demouserS := model.ToSessionToken("S123", demouser, demoRoles, now, now.Add(model.LoginTimeout))

	testCases := []struct {
		Alias           string
		Path            string
		OutputResponse  model.LoginResponse
		RequestMethod   string
		RequestHeaders  map[string]string
		RequestBody     []byte
		ResponseStatus  int
		ResponseHeaders map[string]string
		ResponseBody    []byte
	}{
		{
			Alias: "login success",
			Path:  "/login",
			OutputResponse: model.LoginResponse{
				Code:         http.StatusOK,
				SessionToken: demouserS,
			},
			RequestMethod: http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
			},
			RequestBody:    []byte(`{"username": "demouser", "password": "demouser"}`),
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hContentLength: "3",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
				hSetCookie:     getCookieString(demouserS, testPassword),
			},
			ResponseBody: []byte("OK\n"),
		},
		{
			Alias:          "login bad method",
			Path:           "/login",
			OutputResponse: model.LoginResponse{},
			RequestMethod:  http.MethodDelete,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusMethodNotAllowed,
			ResponseHeaders: map[string]string{
				hContentLength: "19",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte("Method Not Allowed\n"),
		},
		{
			Alias: "login no content-type",
			Path:  "/login",
			OutputResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "bad or missing content type",
			},
			RequestMethod:  http.MethodPost,
			RequestHeaders: nil,
			RequestBody:    []byte(`{"username": "demouser", "password": "demouser"}`),
			ResponseStatus: http.StatusBadRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "28",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte("bad or missing content type\n"),
		},
		{
			Alias: "login no body",
			Path:  "/login",
			OutputResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "cannot read request body",
			},
			RequestMethod: http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusBadRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "22",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte("cannot unmarshal json\n"),
		},
		{
			Alias: "login bogus request payload",
			Path:  "/login",
			OutputResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "cannot unmarshal json",
			},
			RequestMethod: http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
			},
			RequestBody:    []byte(`{"username"}`),
			ResponseStatus: http.StatusBadRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "22",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte("cannot unmarshal json\n"),
		},
		{
			Alias: "any backend error",
			Path:  "/login",
			OutputResponse: model.LoginResponse{
				Code:    http.StatusInternalServerError,
				Message: "real bad error",
			},
			RequestMethod: http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
			},
			RequestBody:    []byte(`{"username": "demouser", "password": "demouser"}`),
			ResponseStatus: http.StatusInternalServerError,
			ResponseHeaders: map[string]string{
				hContentLength: "15",
				hContentType:   mimeTypeText,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte("real bad error\n"),
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

			us := &UserServiceMock{
				LoginResponse: testCase.OutputResponse,
			}
			handler := chi.NewRouter()
			handler.Post("/login", auth.Login(us, testPassword))
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
