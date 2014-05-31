package main

import (
	"reflect"
	"testing"
)

func AssertTypeMatch(t *testing.T, a interface{}, expectedType interface{}) {
	expected := reflect.TypeOf(expectedType).Elem()
	if reflect.TypeOf(a) != expected {
		t.Errorf("%v was expected to be of type: %v - Got: %v", a, expected, reflect.TypeOf(a))
	}
}

func AssertTrue(t *testing.T, a bool) {
	if !a {
		t.Errorf("'true' (type bool) was expected - Got %v (type %v)", a, reflect.TypeOf(a))
	}
}

func AssertFalse(t *testing.T, a bool) {
	if a {
		t.Errorf("'false' (type bool) was expected - Got %v (type %v)", a, reflect.TypeOf(a))
	}
}

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func AssertDeepEqual(t *testing.T, a interface{}, b interface{}) {
	if eq := reflect.DeepEqual(a, b); !eq {
		t.Errorf("Expected deep equivalency %v (type %v) did not match %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func AssertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Expected %v (type %v) NOT to match: %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
