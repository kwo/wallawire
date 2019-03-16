package repository_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"wallawire/idgen"
	"wallawire/model"
	"wallawire/repository"
)

const (
	userIDFakeuser   = "0f644276-eace-40a4-8adc-8f39992a4389"
	userIDGuest      = "856c18d6-fa9c-42d0-9c9f-35406b9a6702"
	roleIDEditor     = "7934260a-305d-4381-be64-3560f0885881"
	roleIDReporter   = "b62a31ac-97ab-472c-8908-0e189ea3994e"
	roleIDCopywriter = "6062c2d0-98ce-4cfe-a33d-29590406f25b"
	roleIDStaff      = "f9ce6462-1d65-4223-ac08-5e29ae93d918"
)

func init() {

	sStatements := []string{
		fmt.Sprintf("INSERT INTO users (id, disabled, username, name, created, updated, password_hash) VALUES ('%s', TRUE, 'guestuser', 'Guest User', %d, %d, '2432612431302463717561632f524161654d32594251717a63717a356578562e46305934503175673241304673505855686965797059525465684665')", userIDGuest, now.Unix(), now.Unix()),
		fmt.Sprintf("INSERT INTO users (id, disabled, username, name, created, updated, password_hash) VALUES ('%s', FALSE, 'fakeuser', 'Fake User', %d, %d, '243261243130244b546b346e4649547463526458745079555a79685775725651456a654650375a517853442e77532e64356b4367596979426c795871')", userIDFakeuser, now.Unix(), now.Unix()),
		fmt.Sprintf("INSERT INTO roles (id, name) VALUES ('%s', 'editor')", roleIDEditor),
		fmt.Sprintf("INSERT INTO roles (id, name) VALUES ('%s', 'reporter')", roleIDReporter),
		fmt.Sprintf("INSERT INTO roles (id, name) VALUES ('%s', 'copywriter')", roleIDCopywriter),
		fmt.Sprintf("INSERT INTO roles (id, name) VALUES ('%s', 'staff')", roleIDStaff),
		fmt.Sprintf("INSERT INTO user_role (user_id, role_id, valid_from, valid_to) VALUES ('%s', '%s', %d, %d)", userIDFakeuser, roleIDEditor, now1h.Unix(), now3d.Unix()),
		fmt.Sprintf("INSERT INTO user_role (user_id, role_id, valid_from, valid_to) VALUES ('%s', '%s', %d, NULL)", userIDFakeuser, roleIDReporter, now1h.Unix()),
		fmt.Sprintf("INSERT INTO user_role (user_id, role_id, valid_from, valid_to) VALUES ('%s', '%s', NULL, %d)", userIDFakeuser, roleIDCopywriter, now1h.Unix()),
		fmt.Sprintf("INSERT INTO user_role (user_id, role_id, valid_from, valid_to) VALUES ('%s', '%s', NULL, NULL)", userIDFakeuser, roleIDStaff),
	}

	tStatements := []string{
		fmt.Sprintf("DELETE FROM user_role WHERE user_id = '%s'", userIDGuest),
		fmt.Sprintf("DELETE FROM user_role WHERE user_id = '%s'", userIDFakeuser),
		fmt.Sprintf("DELETE FROM users WHERE id = '%s'", userIDGuest),
		fmt.Sprintf("DELETE FROM users WHERE id = '%s'", userIDFakeuser),
		fmt.Sprintf("DELETE FROM roles WHERE id = '%s'", roleIDReporter),
		fmt.Sprintf("DELETE FROM roles WHERE id = '%s'", roleIDEditor),
		fmt.Sprintf("DELETE FROM roles WHERE id = '%s'", roleIDCopywriter),
		fmt.Sprintf("DELETE FROM roles WHERE id = '%s'", roleIDStaff),
	}

	addTestStatements(sStatements, tStatements)

}

func TestUser(b *testing.T) {

	testCases := []struct {
		Alias                     string
		User                      model.User
		User2                     model.User
		ExpectedAvailibility      bool
		ExpectedAvailabilityError error
		ExpectedGetError          error
		ExpectedSetError          error
		ExpectedUpdateError       error
		ExpectedReGetError        error
		ExpectedDeleteError       error
	}{
		{
			Alias: "success",
			User: model.User{
				ID:           "50b2a050-7b90-4cb1-84c3-327b39847fa6",
				Disabled:     false,
				Username:     "TestUser1",
				Name:         "Test User1",
				PasswordHash: "passwordhash",
				Created:      now.UTC(),
				Updated:      now.UTC(),
			},
			User2: model.User{
				ID:           "50b2a050-7b90-4cb1-84c3-327b39847fa6",
				Disabled:     true,
				Username:     "TestUser11",
				Name:         "Test User11",
				PasswordHash: "passwordhash2",
				Created:      now.UTC(),
				Updated:      now.UTC(),
			},
			ExpectedAvailibility:      true,
			ExpectedAvailabilityError: nil,
			ExpectedGetError:          nil,
			ExpectedSetError:          nil,
			ExpectedUpdateError:       nil,
			ExpectedReGetError:        nil,
			ExpectedDeleteError:       nil,
		},
	}

	database := repository.NewDatabase(db)
	us := repository.New(idgen.NewUUIDGenerator())

	for _, tc := range testCases {

		testFn := func(t *testing.T) {

			err := database.Run(func(tx model.Transaction) error {

				ctx := context.Background()

				// Available
				available, errAvailable := us.IsUsernameAvailable(ctx, tx, tc.User.Username)
				if available != tc.ExpectedAvailibility {
					t.Errorf("Bad availibility: %t, expected %t", available, tc.ExpectedAvailibility)
				}
				if errAvailable != tc.ExpectedAvailabilityError {
					t.Errorf("Bad availibility error: %s, expected %s", errAvailable, tc.ExpectedAvailabilityError)
				}

				// Set
				errSet := us.SetUser(ctx, tx, tc.User)
				if errSet != tc.ExpectedSetError {
					t.Errorf("Bad set error: %s, expected %s", errSet, tc.ExpectedSetError)
				}

				// errRolesSet := us.S

				// Get
				u, errGet := us.GetActiveUserByUsername(ctx, tx, tc.User.Username)
				if errGet != tc.ExpectedGetError {
					t.Errorf("Bad get error: %s, expected %s", errGet, tc.ExpectedGetError)
				}
				if !reflect.DeepEqual(u, &tc.User) {
					t.Errorf("Bad user: %v, expected %v", u, tc.User)
				}

				// Update
				errUpdate := us.SetUser(ctx, tx, tc.User2)
				if errUpdate != tc.ExpectedUpdateError {
					t.Errorf("Bad update error: %s, expected %s", errUpdate, tc.ExpectedUpdateError)
				}

				// ReGet
				u2, errReGet := us.GetUser(ctx, tx, tc.User.ID)
				if errReGet != tc.ExpectedReGetError {
					t.Errorf("Bad re-get error: %s, expected %s", errReGet, tc.ExpectedReGetError)
				}
				if !reflect.DeepEqual(u2, &tc.User2) {
					t.Errorf("Bad user: %v, expected %v", u2, tc.User2)
				}

				// Delete
				errDelete := us.DeleteUser(ctx, tx, tc.User.ID)
				if errDelete != tc.ExpectedDeleteError {
					t.Errorf("Bad delete error: %s, expected %s", errDelete, tc.ExpectedDeleteError)
				}

				return nil // always nil, so don't test database.Run return value

			})

			if err != nil {
				t.Error(err)
			}

		}

		b.Run(tc.Alias, testFn)

	}

}

func TestGetUser(b *testing.T) {

	testCases := []struct {
		Alias            string
		UserID           string
		User             *model.User
		ExpectedGetError error
	}{
		{
			Alias:  "normal user",
			UserID: userIDFakeuser,
			User: &model.User{
				ID:           userIDFakeuser,
				Disabled:     false,
				Username:     "fakeuser",
				Name:         "Fake User",
				PasswordHash: "243261243130244b546b346e4649547463526458745079555a79685775725651456a654650375a517853442e77532e64356b4367596979426c795871",
				Created:      now.UTC(),
				Updated:      now.UTC(),
			},
			ExpectedGetError: nil,
		},
		{
			Alias:  "disabled user",
			UserID: userIDGuest,
			User: &model.User{
				ID:           userIDGuest,
				Disabled:     true,
				Username:     "guestuser",
				Name:         "Guest User",
				PasswordHash: "2432612431302463717561632f524161654d32594251717a63717a356578562e46305934503175673241304673505855686965797059525465684665",
				Created:      now.UTC(),
				Updated:      now.UTC(),
			},
			ExpectedGetError: nil,
		},
	}

	database := repository.NewDatabase(db)
	us := repository.New(idgen.NewUUIDGenerator())

	for _, tc := range testCases {

		testFn := func(t *testing.T) {

			err := database.Run(func(tx model.Transaction) error {

				ctx := context.Background()

				// Get
				u, errGet := us.GetUser(ctx, tx, tc.UserID)
				if errGet != tc.ExpectedGetError {
					t.Errorf("Bad get error: %s, expected %s", errGet, tc.ExpectedGetError)
				}
				if !reflect.DeepEqual(u, tc.User) {
					t.Errorf("Bad user: %v, expected %v", u, tc.User)
				}

				return nil // always nil, so don't test database.Run return value

			})

			if err != nil {
				t.Error(err)
			}

		}

		b.Run(tc.Alias, testFn)

	}

}

func TestGetActiveUserByUsername(b *testing.T) {

	testCases := []struct {
		Alias            string
		Username         string
		User             *model.User
		ExpectedGetError error
	}{
		{
			Alias:    "normal user",
			Username: "fakeuser",
			User: &model.User{
				ID:           userIDFakeuser,
				Disabled:     false,
				Username:     "fakeuser",
				Name:         "Fake User",
				PasswordHash: "243261243130244b546b346e4649547463526458745079555a79685775725651456a654650375a517853442e77532e64356b4367596979426c795871",
				Created:      now.UTC(),
				Updated:      now.UTC(),
			},
			ExpectedGetError: nil,
		},
		{
			Alias:            "disabled user",
			Username:         "guestuser",
			User:             nil,
			ExpectedGetError: nil,
		},
	}

	database := repository.NewDatabase(db)
	us := repository.New(idgen.NewUUIDGenerator())

	for _, tc := range testCases {

		testFn := func(t *testing.T) {

			err := database.Run(func(tx model.Transaction) error {

				ctx := context.Background()

				// Get
				u, errGet := us.GetActiveUserByUsername(ctx, tx, tc.Username)
				if errGet != tc.ExpectedGetError {
					t.Errorf("Bad get error: %s, expected %s", errGet, tc.ExpectedGetError)
				}
				if !reflect.DeepEqual(u, tc.User) {
					t.Errorf("Bad user: %v, expected %v", u, tc.User)
				}

				return nil // always nil, so don't test database.Run return value

			})

			if err != nil {
				t.Error(err)
			}

		}

		b.Run(tc.Alias, testFn)

	}

}

func TestIsUsernameAvailable(b *testing.T) {

	testCases := []struct {
		Alias          string
		LookupUsername string
		ExpectedError  error
		ExpectedResult bool
	}{
		{
			Alias:          "available",
			LookupUsername: "fakeuser1",
			ExpectedError:  nil,
			ExpectedResult: true,
		},
		{
			Alias:          "not available",
			LookupUsername: "fakeuser",
			ExpectedError:  nil,
			ExpectedResult: false,
		},
		{
			Alias:          "not available mixed case",
			LookupUsername: "FakeUser",
			ExpectedError:  nil,
			ExpectedResult: false,
		},
		{
			Alias:          "not available - disabled",
			LookupUsername: "guestuser",
			ExpectedError:  nil,
			ExpectedResult: false,
		},
	}

	database := repository.NewDatabase(db)
	us := repository.New(idgen.NewUUIDGenerator())

	for _, testCase := range testCases {

		testFn := func(t *testing.T) {

			err := database.Run(func(tx model.Transaction) error {

				ctx := context.Background()

				available, errAvailable := us.IsUsernameAvailable(ctx, tx, testCase.LookupUsername)
				if errAvailable != testCase.ExpectedError {
					t.Errorf("bad error: %s, expected %s", errAvailable, testCase.ExpectedError)
				}
				if available != testCase.ExpectedResult {
					t.Errorf("bad result: %t, expeted %t", available, testCase.ExpectedResult)
				}

				return nil // always nil, so don't test database.Run return value

			})

			if err != nil {
				t.Error(err)
			}

		}

		b.Run(testCase.Alias, testFn)

	}

}

func TestGetUserRoles(b *testing.T) {

	testCases := []struct {
		Alias         string
		UserID        string
		ReferenceTime *time.Time
		ExpectedRoles []model.UserRole
	}{
		{
			Alias:         "all",
			UserID:        userIDFakeuser,
			ReferenceTime: nil,
			ExpectedRoles: []model.UserRole{
				{
					ID:        roleIDCopywriter,
					Name:      "copywriter",
					ValidFrom: nil,
					ValidTo:   &now1h,
				},
				{
					ID:        roleIDEditor,
					Name:      "editor",
					ValidFrom: &now1h,
					ValidTo:   &now3d,
				},
				{
					ID:        roleIDReporter,
					Name:      "reporter",
					ValidFrom: &now1h,
					ValidTo:   nil,
				},
				{
					ID:        roleIDStaff,
					Name:      "staff",
					ValidFrom: nil,
					ValidTo:   nil,
				},
			},
		},
		{
			Alias:         "now",
			UserID:        userIDFakeuser,
			ReferenceTime: &now,
			ExpectedRoles: []model.UserRole{
				{
					ID:        roleIDCopywriter,
					Name:      "copywriter",
					ValidFrom: nil,
					ValidTo:   &now1h,
				},
				{
					ID:        roleIDStaff,
					Name:      "staff",
					ValidFrom: nil,
					ValidTo:   nil,
				},
			},
		},
		{
			Alias:         "2h",
			UserID:        userIDFakeuser,
			ReferenceTime: &now2h,
			ExpectedRoles: []model.UserRole{
				{
					ID:        roleIDEditor,
					Name:      "editor",
					ValidFrom: &now1h,
					ValidTo:   &now3d,
				},
				{
					ID:        roleIDReporter,
					Name:      "reporter",
					ValidFrom: &now1h,
					ValidTo:   nil,
				},
				{
					ID:        roleIDStaff,
					Name:      "staff",
					ValidFrom: nil,
					ValidTo:   nil,
				},
			},
		},
		{
			Alias:         "5d",
			UserID:        userIDFakeuser,
			ReferenceTime: &now5d,
			ExpectedRoles: []model.UserRole{
				{
					ID:        roleIDReporter,
					Name:      "reporter",
					ValidFrom: &now1h,
					ValidTo:   nil,
				},
				{
					ID:        roleIDStaff,
					Name:      "staff",
					ValidFrom: nil,
					ValidTo:   nil,
				},
			},
		},
		{
			Alias:         "none",
			UserID:        userIDGuest,
			ReferenceTime: nil,
			ExpectedRoles: []model.UserRole{},
		},
	}

	database := repository.NewDatabase(db)
	us := repository.New(idgen.NewUUIDGenerator())

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			err := database.Run(func(tx model.Transaction) error {

				ctx := context.Background()

				roles, errRoles := us.GetUserRoles(ctx, tx, tCase.UserID, tCase.ReferenceTime)
				if errRoles != nil {
					t.Fatal(errRoles)
				}

				if roles == nil && tCase.ExpectedRoles != nil {
					t.Errorf("nil roles, expected %#v", tCase.ExpectedRoles)
				} else if roles != nil && tCase.ExpectedRoles == nil {
					t.Errorf("bad roles %#v, expected nil", roles)
				}

				if got, want := len(roles), len(tCase.ExpectedRoles); got != want {
					t.Fatalf("bad number of roles %d, expected %d", got, want)
				}

				for i, role := range roles {

					expectedRole := tCase.ExpectedRoles[i]

					if got, want := role.ID, expectedRole.ID; got != want {
						t.Errorf("bad role ID %s, expected %s", got, want)
					}
					if got, want := role.Name, expectedRole.Name; got != want {
						t.Errorf("bad role name %s, expected %s", got, want)
					}

					compareTimes(t, role.ValidFrom, expectedRole.ValidFrom)
					compareTimes(t, role.ValidTo, expectedRole.ValidTo)

				}

				return nil // always nil, so don't test database.Run return value

			})

			if err != nil {
				t.Error(err)
			}

		}

		b.Run(tCase.Alias, testFn)

	}

}

func TestSetUserRoles(b *testing.T) {

	testCases := []struct {
		Alias  string
		UserID string
		Roles  []model.UserRole
	}{
		{
			Alias:  "all",
			UserID: userIDGuest,
			Roles: []model.UserRole{
				{
					ID:        roleIDCopywriter,
					Name:      "copywriter",
					ValidFrom: nil,
					ValidTo:   &now1h,
				},
				{
					ID:        roleIDEditor,
					Name:      "editor",
					ValidFrom: &now1h,
					ValidTo:   &now3d,
				},
				{
					ID:        roleIDReporter,
					Name:      "reporter",
					ValidFrom: &now1h,
					ValidTo:   nil,
				},
				{
					ID:        roleIDStaff,
					Name:      "staff",
					ValidFrom: nil,
					ValidTo:   nil,
				},
			},
		},
		{
			Alias:  "just user",
			UserID: userIDGuest,
			Roles: []model.UserRole{
				{
					ID:        roleIDStaff,
					Name:      "staff",
					ValidFrom: &now1h,
					ValidTo:   &now3d,
				},
			},
		},
		{
			Alias:  "none",
			UserID: userIDGuest,
			Roles:  []model.UserRole{},
		},
	}

	database := repository.NewDatabase(db)
	us := repository.New(idgen.NewUUIDGenerator())

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			err := database.Run(func(tx model.Transaction) error {

				ctx := context.Background()

				if err := us.SetUserRoles(ctx, tx, tCase.UserID, tCase.Roles); err != nil {
					t.Fatal(err)
				}

				roles, errRoles := us.GetUserRoles(ctx, tx, tCase.UserID, nil)
				if errRoles != nil {
					t.Fatal(errRoles)
				}

				if roles == nil && tCase.Roles != nil {
					t.Errorf("nil roles, expected %#v", tCase.Roles)
				} else if roles != nil && tCase.Roles == nil {
					t.Errorf("bad roles %#v, expected nil", roles)
				}

				if got, want := len(roles), len(tCase.Roles); got != want {
					t.Fatalf("bad number of roles %d, expected %d", got, want)
				}

				for i, role := range roles {

					expectedRole := tCase.Roles[i]

					if got, want := role.ID, expectedRole.ID; got != want {
						t.Errorf("bad role ID %s, expected %s", got, want)
					}
					if got, want := role.Name, expectedRole.Name; got != want {
						t.Errorf("bad role name %s, expected %s", got, want)
					}

					compareTimes(t, role.ValidFrom, expectedRole.ValidFrom)
					compareTimes(t, role.ValidTo, expectedRole.ValidTo)

				}

				return nil // always nil, so don't test database.Run return value

			})

			if err != nil {
				t.Error(err)
			}

		}

		b.Run(tCase.Alias, testFn)

	}

}
