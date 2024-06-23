package utils

// Zipper goes through two sorted slices and calls:
// - the left function for elements only in the first slice,
// - the center function for elements in both slices
// - the right function for elements only in the second slice.
func Zipper[U, V any](u []U, v []V, cmp func(U, V) int, left func(u U), center func(u U, v V), right func(v V)) []any {
	var i, j int
	var out []any
	for i < len(u) && j < len(v) {
		switch cmp(u[i], v[j]) {
		case -1:
			left(u[i])
			i++
		case 1:
			right(v[j])
			j++
		default:
			center(u[i], v[j])
			i++
			j++
		}
	}

	for ; i < len(u); i++ {
		left(u[i])
	}

	for ; j < len(v); j++ {
		right(v[j])
	}

	return out
}
