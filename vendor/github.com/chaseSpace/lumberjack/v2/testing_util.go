package lumberjack

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// isNilUp is like isNil, but used inside helper functions, to ensure that the
// file and line number reported by failures corresponds to one or more levels
// up the stack.
func isNilUp(obtained interface{}, t testing.TB, caller int) {
	if !_isNil(obtained) {
		_, file, line, _ := runtime.Caller(caller + 1)
		fmt.Printf("%s:%d: expected nil, got: %v\n", filepath.Base(file), line, obtained)
		t.FailNow()
	}
}

// _isNil is a helper function for isNil and notNil, and should not be used
// directly.
func _isNil(obtained interface{}) bool {
	if obtained == nil {
		return true
	}

	switch v := reflect.ValueOf(obtained); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}

// equalsUp is like equals, but used inside helper functions, to ensure that the
// file and line number reported by failures corresponds to one or more levels
// up the stack.
func equalsUp(exp, act interface{}, t testing.TB, caller int) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(caller + 1)
		fmt.Printf("%s:%d: exp: %v (%T), got: %v (%T)\n",
			filepath.Base(file), line, exp, exp, act, act)
		t.FailNow()
	}
}
