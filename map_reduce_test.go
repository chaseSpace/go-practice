package main_test

import (
	"fmt"
	"testing"

	"github.com/chaseSpace/bear"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	ID   int
	Name string
	Age  int
}

func TestMapReduce(t *testing.T) {
	var s1 = bear.NewSlice("1", "2", "2", "3")
	s1.Unique().Reverse().Map(func(i string) string { return i + "x" })
	fmt.Println(s1.Slice()) // 3x 2x 1x

	var s2 = bear.NewOrderedSlice(1, 2, 2, 3)
	var joined = s2.Unique().Sort(true).Join(".")
	fmt.Println(joined) // 3.2.1

	var s3 = bear.NewComputableSlice(2.0, 3.0, 4.2)
	var sum = s3.Append(5.0).PopLeft().Append(6.0).Avg()
	fmt.Println("Sum:", sum) // 4.55

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

	fmt.Println(string(hashedPassword), string(hashedPassword2))

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte("123x"))
	fmt.Println(err)
}
