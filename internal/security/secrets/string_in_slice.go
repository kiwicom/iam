package secrets

import "strings"

// stringInSlice checks whether str matches (case-insensitive) any string in slice.
func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
