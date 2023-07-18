package main

import (
	"os"
	"runtime"
)

func main() {
	f, _ := os.OpenFile("111.log", os.O_WRONLY|os.O_APPEND, 666)
	defer f.Close()
	os.Stderr = f
	runtime.WriteFile = func() {
		f.WriteString("222")
	}
	var c chan int
	<-c
}
