package helper

import (
	"fmt"
	"runtime"
)

func LogErrors(a ...any) {
	for _, x := range a {
		if err, ok := x.(error); ok {
			fmt.Printf("error: %v\n", err)

			// Print the stack trace
			stackBuf := make([]byte, 1024) // Start with 1KB buffer
			n := runtime.Stack(stackBuf, false)
			if n == len(stackBuf) {
				// Buffer might not be big enough; resize if necessary
				stackBuf = make([]byte, 2*n)
				n = runtime.Stack(stackBuf, false)
			}
			fmt.Printf("stack trace:\n%s\n", stackBuf[:n])
		}
	}
}
