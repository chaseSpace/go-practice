package main

import (
	"testing"
)

func TestX(t *testing.T) {
	g := func() int { return 10 }

	for i := 0; i < g(); i++ {
		t.Log(i)
	}
}
