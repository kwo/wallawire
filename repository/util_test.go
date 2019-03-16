package repository

import (
	"reflect"
	"testing"

	"wallawire/model"
)

func TestGetComplement(b *testing.T) {

	testCases := []struct {
		Alias          string
		A              []model.UserRole
		B              []model.UserRole
		ExpectedResult []model.UserRole
	}{
		{
			Alias: "success",
			A: []model.UserRole{
				{
					ID: "123",
				},
				{
					ID: "456",
				},
				{
					ID: "789",
				},
			},
			B: []model.UserRole{
				{
					ID: "123",
				},
				{
					ID: "456",
				},
			},
			ExpectedResult: []model.UserRole{
				{
					ID: "789",
				},
			},
		},
		{
			Alias: "nothing",
			A: []model.UserRole{
				{
					ID: "123",
				},
				{
					ID: "456",
				},
			},
			B: []model.UserRole{
				{
					ID: "123",
				},
				{
					ID: "456",
				},
				{
					ID: "789",
				},
			},
			ExpectedResult: nil,
		},
	}

	for _, tCase := range testCases {

		testFn := func(t *testing.T) {

			us := &Repository{}
			result := us.findComplement(tCase.A, tCase.B)

			if !reflect.DeepEqual(result, tCase.ExpectedResult) {
				t.Errorf("bad result %#v, expected %#v", result, tCase.ExpectedResult)
			}

		}

		b.Run(tCase.Alias, testFn)

	}

}
