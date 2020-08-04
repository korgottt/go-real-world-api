package utils

import "strings"

// CreateSlug creates article slug from title
func CreateSlug(title string) string {
	return strings.ToLower(strings.Join(strings.Fields(title), "-"))
}
