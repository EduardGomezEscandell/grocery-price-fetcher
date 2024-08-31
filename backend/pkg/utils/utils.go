package utils

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

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

// Partition rearranges a slice such that the elements that satisfy the predicate
// come before the elements that do not. It returns the index of the first element
// that does not satisfy the predicate.
//
// The operation is stable, in place, and linear in time.
func Partition[T any](slice []T, predicate func(t T) bool) (p int) {
	if len(slice) == 0 {
		return 0
	}
	for i := range slice {
		if !predicate(slice[i]) {
			continue
		}
		if i == p {
			p++
			continue
		}
		slice[i], slice[p] = slice[p], slice[i]
		p++
	}
	return p
}

func SafeIntConvert[DST constraints.Integer, SRC constraints.Integer](u SRC) (DST, error) {
	v := DST(u)
	if SRC(v) != u {
		return 0, fmt.Errorf("cannot convert: value %d is out of range", u)
	}
	return v, nil
}
