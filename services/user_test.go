package services_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"wallawire/idgen"
	"wallawire/model"
	"wallawire/services"
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

	testCases := []struct {
		Alias                string
		OutputUser           *model.User
		OutputAvailability   bool
		OutputAvailableError error
		OutputGetError       error
		OutputSetError       error
		OutputRoles          []model.UserRole
		OutputRolesError     error
		RequestSessionToken  model.SessionToken
		Request              model.ChangeUsernameRequest
		ExpectedResponse     model.ChangeUsernameResponse
	}{
		{
			Alias:                "success",
			OutputUser:           demouser(),
			OutputAvailability:   true,
			OutputAvailableError: nil,
			OutputGetError:       nil,
			OutputSetError:       nil,
			OutputRoles: []model.UserRole{
				{
					ID:   "monsterrole-id",
					Name: "monster",
				},
			},
			OutputRolesError: nil,
			RequestSessionToken: model.SessionToken{
				SessionID: "4567",
			},
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code: http.StatusOK,
				SessionToken: &model.SessionToken{
					SessionID: "4567",
					ID:        "id",
					Username:  "demouser2",
					Name:      "Demo User",
					Roles:     []string{},
				},
			},
		},
		{
			Alias:                "user not found",
			OutputUser:           nil,
			OutputAvailability:   true,
			OutputAvailableError: nil,
			OutputGetError:       nil,
			OutputSetError:       nil,
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusNotFound,
				Message: "user not found",
			},
		},
		{
			Alias:                "bad password",
			OutputUser:           demouser(),
			OutputAvailability:   true,
			OutputAvailableError: nil,
			OutputGetError:       nil,
			OutputSetError:       nil,
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser2",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusBadRequest,
				Message: "password incorrect",
			},
		},
		{
			Alias:                "bad username",
			OutputUser:           demouser(),
			OutputAvailability:   true,
			OutputAvailableError: nil,
			OutputGetError:       nil,
			OutputSetError:       nil,
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "ko",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusBadRequest,
				Message: "username not valid",
			},
		},
		{
			Alias:                "username not available",
			OutputUser:           demouser(),
			OutputAvailability:   false,
			OutputAvailableError: nil,
			OutputGetError:       nil,
			OutputSetError:       nil,
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusBadRequest,
				Message: "username not available",
			},
		},
		{
			Alias:                "availibility fails",
			OutputUser:           demouser(),
			OutputAvailability:   false,
			OutputAvailableError: errors.New("just some error"),
			OutputGetError:       nil,
			OutputSetError:       nil,
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
		{
			Alias:                "get user fails",
			OutputUser:           demouser(),
			OutputAvailability:   true,
			OutputAvailableError: nil,
			OutputGetError:       errors.New("just some error"),
			OutputSetError:       nil,
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
		{
			Alias:                "set user fails",
			OutputUser:           demouser(),
			OutputAvailability:   true,
			OutputAvailableError: nil,
			OutputGetError:       nil,
			OutputSetError:       errors.New("just some error"),
			Request: model.ChangeUsernameRequest{
				UserID:      "id",
				Password:    "demouser",
				NewUsername: "demouser2",
			},
			ExpectedResponse: model.ChangeUsernameResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
	}

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			idg := idgen.NewIdGenerator()

			db := &DatabaseMock{}
			userRepo := &UserRepositoryMock{
				User:           tCase.OutputUser,
				Available:      tCase.OutputAvailability,
				AvailableError: tCase.OutputAvailableError,
				GetError:       tCase.OutputGetError,
				SetError:       tCase.OutputSetError,
			}
			userService := services.NewUserService(db, userRepo, idg)
			ctx := context.Background()
			ctx = context.WithValue(ctx, model.UserKey, tCase.RequestSessionToken)

			rsp := userService.ChangeUsername(ctx, tCase.Request)

			if got, want := rsp.Code, tCase.ExpectedResponse.Code; got != want {
				t.Errorf("bad response code %d, expected %d", got, want)
			}

			if got, want := rsp.Message, tCase.ExpectedResponse.Message; got != want {
				t.Errorf("bad response message %s, expected %s", got, want)
			}

			if rsp.SessionToken == nil && tCase.ExpectedResponse.SessionToken != nil {
				t.Fatal("nil response session, expected non-nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken == nil {
				t.Fatal("non-nil response session, expected nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken != nil {
				if got, want := rsp.SessionToken.SessionID, tCase.ExpectedResponse.SessionToken.SessionID; got != want {
					t.Errorf("bad response sessiontoken SessionID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.ID, tCase.ExpectedResponse.SessionToken.ID; got != want {
					t.Errorf("bad response sessiontoken ID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Username, tCase.ExpectedResponse.SessionToken.Username; got != want {
					t.Errorf("bad response sessiontoken Username %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Name, tCase.ExpectedResponse.SessionToken.Name; got != want {
					t.Errorf("bad response sessiontoken Name %s, expected %s", got, want)
				}

				if got, want := strings.Join(rsp.SessionToken.Roles, ","), strings.Join(tCase.ExpectedResponse.SessionToken.Roles, ","); got != want {
					t.Errorf("bad response sessiontoken Roles %s, expected %s", got, want)
				}
			}

		} // fn

		b.Run(tCase.Alias, testFn)

	} // cases

}

func TestChangeProfile(b *testing.T) {

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

	testCases := []struct {
		Alias               string
		OutputUser          *model.User
		OutputGetError      error
		OutputSetError      error
		OutputRoles         []model.UserRole
		OutputRolesError    error
		RequestSessionToken model.SessionToken
		Request             model.ChangeProfileRequest
		ExpectedResponse    model.ChangeProfileResponse
	}{
		{
			Alias:          "success",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: nil,
			OutputRoles: []model.UserRole{
				{
					ID:   "monsterrole-id",
					Name: "monster",
				},
			},
			OutputRolesError: nil,
			RequestSessionToken: model.SessionToken{
				SessionID: "4567",
			},
			Request: model.ChangeProfileRequest{
				UserID:      "id",
				Displayname: "demouser2",
			},
			ExpectedResponse: model.ChangeProfileResponse{
				Code: http.StatusOK,
				SessionToken: &model.SessionToken{
					SessionID: "4567",
					ID:        "id",
					Username:  "demouser",
					Name:      "demouser2",
				},
			},
		},
		{
			Alias:          "user not found",
			OutputUser:     nil,
			OutputGetError: nil,
			OutputSetError: nil,
			Request: model.ChangeProfileRequest{
				UserID:      "id",
				Displayname: "demouser2",
			},
			ExpectedResponse: model.ChangeProfileResponse{
				Code:    http.StatusNotFound,
				Message: "user not found",
			},
		},
		{
			Alias:          "bad displayname",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: nil,
			Request: model.ChangeProfileRequest{
				UserID:      "id",
				Displayname: "ko",
			},
			ExpectedResponse: model.ChangeProfileResponse{
				Code:    http.StatusBadRequest,
				Message: "displayname not valid",
			},
		},
		{
			Alias:          "get user fails",
			OutputUser:     demouser(),
			OutputGetError: errors.New("just some error"),
			OutputSetError: nil,
			Request: model.ChangeProfileRequest{
				UserID:      "id",
				Displayname: "demouser2",
			},
			ExpectedResponse: model.ChangeProfileResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
		{
			Alias:          "set user fails",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: errors.New("just some error"),
			Request: model.ChangeProfileRequest{
				UserID:      "id",
				Displayname: "demouser2",
			},
			ExpectedResponse: model.ChangeProfileResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
	}

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			db := &DatabaseMock{}
			userRepo := &UserRepositoryMock{
				User:     tCase.OutputUser,
				GetError: tCase.OutputGetError,
				SetError: tCase.OutputSetError,
			}
			userService := services.NewUserService(db, userRepo, nil) // idgen only used for login
			ctx := context.Background()
			ctx = context.WithValue(ctx, model.UserKey, tCase.RequestSessionToken)

			rsp := userService.ChangeProfile(ctx, tCase.Request)

			if got, want := rsp.Code, tCase.ExpectedResponse.Code; got != want {
				t.Errorf("bad response code %d, expected %d", got, want)
			}

			if got, want := rsp.Message, tCase.ExpectedResponse.Message; got != want {
				t.Errorf("bad response message %s, expected %s", got, want)
			}

			if rsp.SessionToken == nil && tCase.ExpectedResponse.SessionToken != nil {
				t.Fatal("nil response session, expected non-nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken == nil {
				t.Fatal("non-nil response session, expected nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken != nil {
				if got, want := rsp.SessionToken.SessionID, tCase.ExpectedResponse.SessionToken.SessionID; got != want {
					t.Errorf("bad response sessiontoken SessionID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.ID, tCase.ExpectedResponse.SessionToken.ID; got != want {
					t.Errorf("bad response sessiontoken ID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Username, tCase.ExpectedResponse.SessionToken.Username; got != want {
					t.Errorf("bad response sessiontoken Username %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Name, tCase.ExpectedResponse.SessionToken.Name; got != want {
					t.Errorf("bad response sessiontoken Name %s, expected %s", got, want)
				}

				if got, want := strings.Join(rsp.SessionToken.Roles, ","), strings.Join(tCase.ExpectedResponse.SessionToken.Roles, ","); got != want {
					t.Errorf("bad response sessiontoken Roles %s, expected %s", got, want)
				}
			}

		} // fn

		b.Run(tCase.Alias, testFn)

	} // cases

}

func TestChangePassword(b *testing.T) {

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

	testCases := []struct {
		Alias               string
		OutputUser          *model.User
		OutputGetError      error
		OutputSetError      error
		OutputRoles         []model.UserRole
		OutputRolesError    error
		RequestSessionToken model.SessionToken
		Request             model.ChangePasswordRequest
		ExpectedResponse    model.ChangePasswordResponse
	}{
		{
			Alias:          "success",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: nil,
			OutputRoles: []model.UserRole{
				{
					ID:   "monsterrole-id",
					Name: "monster",
				},
			},
			OutputRolesError: nil,
			RequestSessionToken: model.SessionToken{
				SessionID: "4567",
			},
			Request: model.ChangePasswordRequest{
				UserID:      "id",
				Password:    "demouser",
				NewPassword: "demouser2",
			},
			ExpectedResponse: model.ChangePasswordResponse{
				Code: http.StatusOK,
				SessionToken: &model.SessionToken{
					SessionID: "4567",
					ID:        "id",
					Username:  "demouser",
					Name:      "Demo User",
					Roles:     []string{"monster"},
				},
			},
		},
		{
			Alias:          "user not found",
			OutputUser:     nil,
			OutputGetError: nil,
			OutputSetError: nil,
			Request: model.ChangePasswordRequest{
				UserID:      "id",
				Password:    "demouser",
				NewPassword: "demouser2",
			},
			ExpectedResponse: model.ChangePasswordResponse{
				Code:    http.StatusNotFound,
				Message: "user not found",
			},
		},
		{
			Alias:          "bad password",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: nil,
			Request: model.ChangePasswordRequest{
				UserID:      "id",
				Password:    "demouser2",
				NewPassword: "demouser2",
			},
			ExpectedResponse: model.ChangePasswordResponse{
				Code:    http.StatusBadRequest,
				Message: "password incorrect",
			},
		},
		{
			Alias:          "bad new password",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: nil,
			Request: model.ChangePasswordRequest{
				UserID:      "id",
				Password:    "demouser",
				NewPassword: "demo",
			},
			ExpectedResponse: model.ChangePasswordResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid new password",
			},
		},
		{
			Alias:          "get user fails",
			OutputUser:     demouser(),
			OutputGetError: errors.New("just some error"),
			OutputSetError: nil,
			Request: model.ChangePasswordRequest{
				UserID:      "id",
				Password:    "demouser",
				NewPassword: "demouser2",
			},
			ExpectedResponse: model.ChangePasswordResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
		{
			Alias:          "set user fails",
			OutputUser:     demouser(),
			OutputGetError: nil,
			OutputSetError: errors.New("just some error"),
			Request: model.ChangePasswordRequest{
				UserID:      "id",
				Password:    "demouser",
				NewPassword: "demouser2",
			},
			ExpectedResponse: model.ChangePasswordResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
	}

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			db := &DatabaseMock{}
			userRepo := &UserRepositoryMock{
				User:       tCase.OutputUser,
				GetError:   tCase.OutputGetError,
				SetError:   tCase.OutputSetError,
				Roles:      tCase.OutputRoles,
				RolesError: tCase.OutputRolesError,
			}
			userService := services.NewUserService(db, userRepo, nil) // only used in login
			ctx := context.Background()
			ctx = context.WithValue(ctx, model.UserKey, tCase.RequestSessionToken)

			rsp := userService.ChangePassword(ctx, tCase.Request)

			if got, want := rsp.Code, tCase.ExpectedResponse.Code; got != want {
				t.Errorf("bad response code %d, expected %d", got, want)
			}

			if got, want := rsp.Message, tCase.ExpectedResponse.Message; got != want {
				t.Errorf("bad response message %s, expected %s", got, want)
			}

			if rsp.SessionToken == nil && tCase.ExpectedResponse.SessionToken != nil {
				t.Fatal("nil response session, expected non-nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken == nil {
				t.Fatal("non-nil response session, expected nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken != nil {
				if got, want := rsp.SessionToken.SessionID, tCase.ExpectedResponse.SessionToken.SessionID; got != want {
					t.Errorf("bad response sessiontoken SessionID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.ID, tCase.ExpectedResponse.SessionToken.ID; got != want {
					t.Errorf("bad response sessiontoken ID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Username, tCase.ExpectedResponse.SessionToken.Username; got != want {
					t.Errorf("bad response sessiontoken Username %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Name, tCase.ExpectedResponse.SessionToken.Name; got != want {
					t.Errorf("bad response sessiontoken Name %s, expected %s", got, want)
				}

				if got, want := strings.Join(rsp.SessionToken.Roles, ","), strings.Join(tCase.ExpectedResponse.SessionToken.Roles, ","); got != want {
					t.Errorf("bad response sessiontoken Roles %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Issued, tCase.ExpectedResponse.SessionToken.Issued; got.Unix() != want.Unix() {
					t.Errorf("bad response sessiontoken Issued %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Expires, tCase.ExpectedResponse.SessionToken.Expires; got.Unix() != want.Unix() {
					t.Errorf("bad response sessiontoken Expires %s, expected %s", got, want)
				}

			}

		} // fn

		b.Run(tCase.Alias, testFn)

	} // cases

}

func TestLogin(b *testing.T) {

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

	userroles := []model.UserRole{
		{
			ID:   "roleid",
			Name: "users",
		},
	}

	testCases := []struct {
		Alias            string
		OutputUser       *model.User
		OutputSessionID  string
		OutputRoles      []model.UserRole
		OutputGetError   error
		OutputRolesError error
		Request          model.LoginRequest
		ExpectedResponse model.LoginResponse
	}{
		{
			Alias:            "success",
			OutputUser:       demouser(),
			OutputRoles:      userroles,
			OutputSessionID:  "S123",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "demouser",
				Password: "demouser",
			},
			ExpectedResponse: model.LoginResponse{
				Code:         http.StatusOK,
				SessionToken: model.ToSessionToken("S123", demouser(), userroles, now.Truncate(time.Minute), now.Truncate(time.Minute).Add(model.LoginTimeout)),
			},
		},
		{
			Alias:            "get user fails",
			OutputUser:       demouser(),
			OutputRoles:      nil,
			OutputSessionID:  "",
			OutputGetError:   errors.New("just some error"),
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "demouser",
				Password: "demouser",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
		{
			Alias:            "user not found",
			OutputUser:       nil,
			OutputRoles:      userroles,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "demouser",
				Password: "demouser",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid username/password",
			},
		},
		{
			Alias:            "bad password",
			OutputUser:       demouser(),
			OutputRoles:      userroles,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "demouser",
				Password: "demouser2",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid username/password",
			},
		},
		{
			Alias:            "get roles fails",
			OutputUser:       demouser(),
			OutputRoles:      nil,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: errors.New("just some error"),
			Request: model.LoginRequest{
				Username: "demouser",
				Password: "demouser",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusInternalServerError,
				Message: "just some error",
			},
		},
		{
			Alias:            "username too short",
			OutputUser:       demouser(),
			OutputRoles:      userroles,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "",
				Password: "demouser",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid username/password",
			},
		},
		{
			Alias:            "username too long",
			OutputUser:       demouser(),
			OutputRoles:      userroles,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: strings.Repeat("x", 65),
				Password: "demouser",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid username/password",
			},
		},
		{
			Alias:            "password too short",
			OutputUser:       demouser(),
			OutputRoles:      userroles,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "demouser",
				Password: "",
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid username/password",
			},
		},
		{
			Alias:            "password too long",
			OutputUser:       demouser(),
			OutputRoles:      userroles,
			OutputSessionID:  "",
			OutputGetError:   nil,
			OutputRolesError: nil,
			Request: model.LoginRequest{
				Username: "demouser",
				Password: strings.Repeat("x", 73),
			},
			ExpectedResponse: model.LoginResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid username/password",
			},
		},
	}

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			idgen := &IdGeneratorMock{
				ID: tCase.OutputSessionID,
			}

			db := &DatabaseMock{}
			userRepo := &UserRepositoryMock{
				User:       tCase.OutputUser,
				Roles:      tCase.OutputRoles,
				GetError:   tCase.OutputGetError,
				RolesError: tCase.OutputRolesError,
			}
			userService := services.NewUserService(db, userRepo, idgen)
			ctx := context.Background()

			rsp := userService.Login(ctx, tCase.Request)

			if got, want := rsp.Code, tCase.ExpectedResponse.Code; got != want {
				t.Errorf("bad response code %d, expected %d", got, want)
			}

			if got, want := rsp.Message, tCase.ExpectedResponse.Message; got != want {
				t.Errorf("bad response message %s, expected %s", got, want)
			}

			if rsp.SessionToken == nil && tCase.ExpectedResponse.SessionToken != nil {
				t.Fatal("nil response session, expected non-nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken == nil {
				t.Fatal("non-nil response session, expected nil")
			}

			if rsp.SessionToken != nil && tCase.ExpectedResponse.SessionToken != nil {
				if got, want := rsp.SessionToken.SessionID, tCase.ExpectedResponse.SessionToken.SessionID; got != want {
					t.Errorf("bad response sessiontoken SessionID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.ID, tCase.ExpectedResponse.SessionToken.ID; got != want {
					t.Errorf("bad response sessiontoken ID %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Username, tCase.ExpectedResponse.SessionToken.Username; got != want {
					t.Errorf("bad response sessiontoken Username %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Name, tCase.ExpectedResponse.SessionToken.Name; got != want {
					t.Errorf("bad response sessiontoken Name %s, expected %s", got, want)
				}

				if got, want := strings.Join(rsp.SessionToken.Roles, ","), strings.Join(tCase.ExpectedResponse.SessionToken.Roles, ","); got != want {
					t.Errorf("bad response sessiontoken Roles %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Issued, tCase.ExpectedResponse.SessionToken.Issued; got.Unix() != want.Unix() {
					t.Errorf("bad response sessiontoken Issued %s, expected %s", got, want)
				}

				if got, want := rsp.SessionToken.Expires, tCase.ExpectedResponse.SessionToken.Expires; got.Unix() != want.Unix() {
					t.Errorf("bad response sessiontoken Expires %s, expected %s", got, want)
				}

			}

		} // fn

		b.Run(tCase.Alias, testFn)

	} // cases

}
