package model_test

import (
	"testing"

	"wallawire/model"
)

func TestValidationError(t *testing.T) {

	x := model.NewValidationError("too long")
	if !model.IsValidationError(x) {
		t.Error("Expected validation error")
	}

	if model.IsNotFoundError(x) {
		t.Error("Unexpected notfound error")
	}

}

func TestNotFoundError(t *testing.T) {

	x := model.NewNotFoundError("not found")

	if !model.IsNotFoundError(x) {
		t.Error("Expected notfound error")
	}

	if model.IsValidationError(x) {
		t.Error("Unexpected validation error")
	}


}
