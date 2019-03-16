package user_test

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
	"wallawire/web/user"
)

func TestChangeUsername(b *testing.T) {

	now := time.Now().Truncate(time.Second)

	demouser := func() *model.User {
		u := &model.User{
			ID:       "id",
			Username: "demouser",
			Name:     "Demo User",
			Created:  now,
			Updated:  now,
		}
		if err := u.SetPassword("demouser"); err != nil {
			b.Fatal(err)
		}
		return u
	}

	demouser2 := func() *model.User {
		u := &model.User{
			ID:       "id",
			Username: "demouser2",
			Name:     "Demo User",
			Created:  now,
			Updated:  now,
		}
		if err := u.SetPassword("demouser"); err != nil {
			b.Fatal(err)
		}
		return u
	}

	demoRoles := []model.UserRole{
		{
			ID:   "roleid",
			Name: "users",
		},
	}

	demouserS := model.ToSessionToken("S123", demouser(), demoRoles, now.Truncate(time.Minute), now.Truncate(time.Minute).Add(model.LoginTimeout))
	demouser2S := model.ToSessionToken("S123", demouser2(), demoRoles, now.Truncate(time.Minute), now.Truncate(time.Minute).Add(model.LoginTimeout))

	testCases := []struct {
		Alias           string
		Path            string
		OutputResponse  model.ChangeUsernameResponse
		RequestMethod   string
		RequestHeaders  map[string]string
		RequestBody     []byte
		ResponseStatus  int
		ResponseHeaders map[string]string
		ResponseBody    []byte
	}{
		{
			Alias: "success",
			Path:  "/changeusername",
			OutputResponse: model.ChangeUsernameResponse{
				Code:         http.StatusOK,
				SessionToken: demouser2S,
			},
			RequestMethod: http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
				hCookie:      getCookieString(demouserS, testPassword),
			},
			RequestBody:    []byte(`{"newusername": "demouser2", "password": "demouser"}`),
			ResponseStatus: http.StatusOK,
			ResponseHeaders: map[string]string{
				hContentLength: "33",
				hContentType:   mimeTypeJson,
				hDate:          ignoreValue,
				hSetCookie:     getCookieString(demouser2S, testPassword),
			},
			ResponseBody: []byte(`{"statusCode":200,"message":"OK"}`),
		},
		{
			Alias:          "unauthorized",
			Path:           "/changeusername",
			OutputResponse: model.ChangeUsernameResponse{},
			RequestMethod:  http.MethodPost,
			RequestHeaders: nil,
			RequestBody:    []byte(`{"newusername": "demouser2", "password": "demouser2"}`),
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
			Alias:          "no content-type",
			Path:           "/changeusername",
			OutputResponse: model.ChangeUsernameResponse{},
			RequestMethod:  http.MethodPost,
			RequestHeaders: map[string]string{
				hCookie: getCookieString(demouserS, testPassword),
			},
			RequestBody:    []byte(`{"newusername": "demouser2", "password": "demouser"}`),
			ResponseStatus: http.StatusBadRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "58",
				hContentType:   mimeTypeJson,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte(`{"statusCode":400,"message":"bad or missing content type"}`),
		},
		{
			Alias:          "no body",
			Path:           "/changeusername",
			OutputResponse: model.ChangeUsernameResponse{},
			RequestMethod:  http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
				hCookie:      getCookieString(demouserS, testPassword),
			},
			RequestBody:    nil,
			ResponseStatus: http.StatusBadRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "47",
				hContentType:   mimeTypeJson,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte(`{"statusCode":400,"message":"bad json payload"}`),
		},
		{
			Alias:          "bogus request payload",
			Path:           "/changeusername",
			OutputResponse: model.ChangeUsernameResponse{},
			RequestMethod:  http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
				hCookie:      getCookieString(demouserS, testPassword),
			},
			RequestBody:    []byte(`{"bogus"}`),
			ResponseStatus: http.StatusBadRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "47",
				hContentType:   mimeTypeJson,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte(`{"statusCode":400,"message":"bad json payload"}`),
		},
		{
			Alias: "any backend error",
			Path:  "/changeusername",
			OutputResponse: model.ChangeUsernameResponse{
				Code:    http.StatusMisdirectedRequest,
				Message: "any old error",
			},
			RequestMethod: http.MethodPost,
			RequestHeaders: map[string]string{
				hContentType: mimeTypeJson,
				hCookie:      getCookieString(demouserS, testPassword),
			},
			RequestBody:    []byte(`{"oldpassword": "demouser", "newpassword": "demouser2"}`),
			ResponseStatus: http.StatusMisdirectedRequest,
			ResponseHeaders: map[string]string{
				hContentLength: "44",
				hContentType:   mimeTypeJson,
				hDate:          ignoreValue,
			},
			ResponseBody: []byte(`{"statusCode":421,"message":"any old error"}`),
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

			us := &UserServiceMock{
				ChangeUsernameResponse: testCase.OutputResponse,
			}

			handler := chi.NewRouter()
			handler.Use(auth.NewAuthenticator(testPassword)...)
			handler.Post("/changeusername", user.ChangeUsername(us, testPassword))

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
				t.Fatalf("Error getting response: %s", errRsp.Error())
			}

			body, errBody := ioutil.ReadAll(rsp.Body)
			if errBody != nil {
				t.Fatalf("Error reading response: %s", errBody.Error())
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
