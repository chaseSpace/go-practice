package panic

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func bar() string {
	errStack := ""
	//stack := string(debug.Stack())
	//fmt.Println("Full stack trace:")
	//fmt.Println(stack)

	// Get the caller information for the third level up
	pc := make([]uintptr, 5)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])

	// Print the caller information
	fmt.Println("Caller information (3 levels up):")
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		errStack += fmt.Sprintf("- %s:%d %s\n", frame.File, frame.Line, frame.Function)
	}
	return errStack
}

func panicFunc(i int) {
	_ = 1 / i
}
func runFormatPanicStack() {

	defer func() {
		if err := recover(); err != nil {
			stack := bar()
			log.Printf("Time: %s\nReason: %v\nStack: \n%s", time.Now(), err, stack)
		}
	}()
	panicFunc(0)
}
