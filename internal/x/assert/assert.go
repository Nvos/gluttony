package assert

func Assert(check bool, message string) {
	if check {
		panic(message)
	}
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
