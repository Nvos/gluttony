package sliceutil

func Merge[T any](slices ...[]T) []T {
	var out []T

	for i := range slices {
		out = append(out, slices[i]...)
	}

	return out
}
