package schema

import (
	"testing"
)

func TestAssetNames(t *testing.T) {

	expectedAssetNames := []string{
		"1_init.sql",
		"2_data.sql",
	}

	names, errNames := getAssetNames("")
	if errNames != nil {
		t.Fatal(errNames)
	}

	if got, want := len(names), len(expectedAssetNames); got < want {
		t.Fatalf("bad asset count: %d, expected at least %d", got, want)
	}

}

func TestAssets(t *testing.T) {

	testAsset := "1_init.sql"

	data, errData := getAsset(testAsset)
	if errData != nil {
		t.Fatal(errData)
	}

	if len(data) == 0 {
		t.Error("bad asset: 0 length")
	}

}
