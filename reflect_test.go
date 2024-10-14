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

func TestSetNilPtr(t *testing.T) {

	a := 1
	var x *int = &a

	rf := reflect.ValueOf(&x).Elem() // 必须要获取指针的引用
	rf.Set(reflect.Zero(rf.Type()))
	println(x)
}
