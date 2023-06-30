package util

import (
	"encoding/json"
	sugar "go_accost/log"
	"runtime/debug"
)

const (
	Datetime = "2006-01-02 15:04:05"
)

func Pretty(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func Recover() {
	if e := recover(); e != nil {
		sugar.Panic("panic-err", e, "STACK=", debug.Stack())
	}
}
