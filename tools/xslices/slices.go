package xslices

import (
	"golang.org/x/exp/maps"
)

// Intersect returns intersection of sl2 and sl1
// the result has only unique values
func Intersect[T comparable](sl1, sl2 []T) []T {
	var intersect []T

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range sl1 {
			found := false
			for _, s2 := range sl2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if found {
				intersect = append(intersect, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			sl1, sl2 = sl2, sl1
		}
	}

	return Unique(intersect)
}

// Unique - returns new slice with unique values only.
func Unique[T comparable](slice []T) []T {
	m := map[T]bool{}
	for _, v := range slice {
		m[v] = true
	}
	return maps.Keys(m)
}

// ConcatUnique concatenate sl2 with sl1, only unique values
func ConcatUnique[T comparable](sl1, sl2 []T) []T {
	res := append(sl1, sl2...)
	return Unique(res)
}
