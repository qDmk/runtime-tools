package util

import "github.com/mndrix/tap-go"

func checkOptionalValue[T comparable](t *tap.T, name string, expected, actual *T) {
	if expected == nil {
		return
	}

	t.Ok(actual != nil && *actual == *expected, name+" is set correctly")
	if actual == nil {
		t.Diagnosticf("expect: %v, actual: nil", *expected)
		return
	}
	t.Diagnosticf("expect: %v, actual: %v", *expected, *actual)
}
