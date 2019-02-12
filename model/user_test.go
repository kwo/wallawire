package model_test

import (
	"fmt"
	"strings"
	"testing"

	"wallawire/model"
)

func TestPasswordCreate(t *testing.T) {

	// t.SkipNow()

	u := model.User{}
	password := "demouser"

	if err := u.SetPassword(password); err != nil {
		t.Fatal(err)
	}

	if !u.MatchPassword(password) {
		t.Errorf("Expected password to match")
	}

	t.Log(u.PasswordHash)

}

func TestPassword(b *testing.T) {

	type testCase struct {
		Alias         string
		Password      string
		MatchPassword string
		ExpectedMatch bool
	}

	testCases := []testCase{
		{
			Alias:         "success",
			Password:      "password",
			MatchPassword: "password",
			ExpectedMatch: true,
		},
	}

	maxLength := 72 // bcrypt max

	for n := 71; n <= 73; n++ {
		password := strings.Repeat("x", n)
		matchword := strings.Repeat("x", n)
		if n > maxLength {
			matchword = matchword[0:maxLength]
		}
		testCases = append(testCases, testCase{
			Alias:         fmt.Sprintf("test%0d", n),
			Password:      password,
			MatchPassword: matchword,
			ExpectedMatch: true,
		})
	}

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			u := model.User{}

			if err := u.SetPassword(tCase.Password); err != nil {
				t.Fatal(err)
			}

			if got, want := u.MatchPassword(tCase.MatchPassword), tCase.ExpectedMatch; got != want {
				t.Errorf("Bad match: %t, expected %t", got, want)
			}

		}

		b.Run(tCase.Alias, testFn)

	}

}
