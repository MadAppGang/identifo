package controller

// sliceContains checks if value is inside slice.
func sliceContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// intersect sl2 from sl1
func intersect[T comparable](sl1, sl2 []T) []T {
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

	return intersect
}

// removeDuplicate
func removeDuplicate[T comparable](slice []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range slice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// concatUnique sl2 with sl1, only unique values
func concatUnique[T comparable](sl1, sl2 []T) []T {
	res := append(sl1, sl2...)
	return removeDuplicate(res)
}
