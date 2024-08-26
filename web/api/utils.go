package api

import "strings"

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(e)) {
			return true
		}
	}
	return false
}
