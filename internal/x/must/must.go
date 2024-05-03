package must

import "fmt"

func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("Must: %v", err))
	}

	return v
}
