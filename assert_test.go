package mwnd

import (
	"cmp"
	"math"
	"reflect"
	"runtime"
	"testing"
)

func printCaller(t *testing.T) {
	_, file, line, _ := runtime.Caller(2)
	t.Logf("\nat %s:%d", file, line)
}

func printMsg(t *testing.T, msg []string) {
	if len(msg) > 0 {
		t.Log(msg[0])
	}
}

func assertEqual(t *testing.T, expected, actual any, msg ...string) bool {
	ok := expected == actual
	if !ok {
		printCaller(t)
		printMsg(t, msg)
		t.Errorf("\nwant: %v\ngot:  %v", expected, actual)
	}
	return ok
}

func assertNil(t *testing.T, actual any, msg ...string) bool {
	if actual == nil {
		return true
	}

	value := reflect.ValueOf(actual)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		if value.IsNil() {
			return true
		}
	}

	printCaller(t)
	printMsg(t, msg)
	t.Errorf("\nwant: nil\ngot:  %v", actual)
	return false
}

func assertLessOrEqual[T cmp.Ordered](t *testing.T, a, b T, msg ...string) bool {
	if a == b {
		return true
	}

	ok := cmp.Less(a, b)
	if !ok {
		printCaller(t)
		printMsg(t, msg)
		t.Errorf("\nexpected a < b, but it was not\na: %v\nb: %v", a, b)
	}
	return ok
}

func assertInDelta(t *testing.T, expected, actual, delta float64, msg ...string) bool {
	if math.IsNaN(expected) && math.IsNaN(actual) {
		return true
	}

	if math.IsNaN(expected) {
		printCaller(t)
		printMsg(t, msg)
		t.Error("\nexpected must not be NaN")
		return false
	}

	if math.IsNaN(actual) {
		printCaller(t)
		printMsg(t, msg)
		t.Errorf("\nexpected %f with delta %f, but was NaN", expected, delta)
		return false
	}

	if diff := math.Abs(expected - actual); diff > delta {
		printCaller(t)
		printMsg(t, msg)
		t.Errorf("\ndelta is allowed to be %v, but was %v", delta, diff)
		return false
	}

	return true
}
