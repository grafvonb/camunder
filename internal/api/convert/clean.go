package convert

import "strings"

// CleanNULL returns "" for sentinel "NULL" (case-insensitive), otherwise returns s unchanged
func CleanNULL(s string) string {
	t := strings.TrimSpace(s)
	if strings.EqualFold(t, "null") {
		return ""
	}
	return s
}

// DerefClean returns "" for nil or sentinel "NULL"
func DerefClean(p *string) string {
	return Deref(MapPtr(p, CleanNULL), "")
}
