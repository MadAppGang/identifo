package model

import "strings"

// SliceIntersect returns only items in as that are found in bs.
// simple intersection of two slices, with complexity: O(n^2)
// there is better algorithms around, this one is simple and scopes are usually 1-3 items in it
func SliceIntersect(a, b []string) []string {
	res := make([]string, 0)

	for _, e := range a {
		if SliceContains(b, e) {
			res = append(res, e)
		}
	}

	return res
}

func SliceContains(s []string, e string) bool {
	el := strings.TrimSpace(e)

	for _, a := range s {
		if strings.EqualFold(strings.TrimSpace(a), el) {
			return true
		}

	}

	return false
}

func SliceExcluding(a []string, exclude string) []string {
	res := make([]string, 0)

	for _, e := range a {
		if e != exclude {
			res = append(res, e)
		}
	}

	return res
}
