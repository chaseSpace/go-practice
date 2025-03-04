package main

import (
	"fmt"
	"reflect"
	"testing"
)

// printReadableTypeValue returns a string representation of the type and value of the given value.
func printReadableTypeValue(value interface{}) string {

	v := reflect.ValueOf(value)
	t := v.Type()

	var valueStr string
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			valueStr = "nil"
		} else {
			valueStr = fmt.Sprintf("%v", v.Elem().Interface())
		}
	} else {
		valueStr = fmt.Sprintf("%v", v.Interface())
	}

	var typeStr string
	if v.Kind() == reflect.Ptr {
		typeStr = fmt.Sprintf("%s(%s)", t, valueStr)
	} else {
		typeStr = fmt.Sprintf("%s(%s)", t, valueStr)
	}

	return typeStr
}

func TestX(t *testing.T) {
	g := func() int { return 10 }

	for i := 0; i < g(); i++ {
		t.Log(i)
	}
}
