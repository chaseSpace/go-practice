package main

import "sort"

type X struct {
	id int
}
type a []*X

func (r a) sort() {
	sort.Slice(r, func(i, j int) bool {
		return r[i].id > r[j].id
	})
}

func test() {

	var xx = a{
		{id: 1},
		{id: 3},
	}
	xx.sort()

	println(xx[0].id)

}
