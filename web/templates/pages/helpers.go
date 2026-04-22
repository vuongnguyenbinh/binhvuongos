package pages

import "strings"

// initialsOf returns 1-2 uppercase letters for avatar fallback.
// "Nguyễn Văn A" → "NA"; "Owner" → "O"; "" → "?".
func initialsOf(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "?"
	}
	if len(parts) == 1 {
		return strings.ToUpper(firstRune(parts[0]))
	}
	return strings.ToUpper(firstRune(parts[0]) + firstRune(parts[len(parts)-1]))
}

// firstRune returns the first UTF-8 rune of s as a string; safe for Vietnamese/Unicode.
func firstRune(s string) string {
	for _, r := range s {
		return string(r)
	}
	return ""
}
