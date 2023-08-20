package main

import "testing"

func TestX(t *testing.T) {
	umap := make(map[int]int)

	x := func(u map[int]int) {
		u[1] = 1
	}
	x(umap)

	println(umap[1])
}
