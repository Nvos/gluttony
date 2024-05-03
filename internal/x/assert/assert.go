package assert

import "fmt"

func Assert(check bool, message string) {
	if !check {
		panic(fmt.Sprintf("Assert: %s", message))
	}
}
