package util

import "strings"

func SanitizeQuery(query string) {
	strings.ToLower(query)
	strings.Replace(query, " ", "%20", 10)
}
