package main

import (
	"reflect"
	"testing"
)

func TestReflectSet(t *testing.T) {
	type x struct {
		I int
	}

	var xx interface{} = (*x)(nil)
	v := reflect.ValueOf(&xx)
	println(v.Kind().String(), v.Elem().Kind().String()) // ptr, interface
	println(v.Elem().Elem().Type().String())             // *x

	v.Elem().Set(reflect.New(v.Elem().Elem().Type().Elem()))
	println(xx == nil, xx.(*x).I)
}
