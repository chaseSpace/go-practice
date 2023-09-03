package main_test

import (
	"fmt"
	"github.com/lyuangg/mr"
	"testing"
)

type Person struct {
	ID   int
	Name string
	Age  int
}

func TestMapReduce(t *testing.T) {
	// Reduce
	numbers := []int{1, 2, 3, 4, 5}
	sum := mr.Reduce(numbers, func(a, b int) int {
		return a + b
	}, 0)
	fmt.Println("Sum:", sum)

	// Contains
	names := []string{"Alice", "Bob", "Charlie"}
	contains := mr.Contains(names, "Bob", func(a string) string { return a })
	fmt.Println("Contains Bob:", contains)

	// Map
	squares := mr.Map(numbers, func(n int) int {
		return n * n
	})
	fmt.Println("Squares:", squares)

	// ToMap
	persons := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	personsMap := mr.ToMap(persons, func(p Person) int {
		return p.ID
	})
	fmt.Println("Persons Map:", personsMap)

	// Filter
	adults := mr.Filter(persons, func(p Person) bool {
		return p.Age >= 18
	})
	fmt.Println("Adults:", adults)
}
