package obsidian

import (
	"strings"
)

var tokenMap = []struct {
	token  string
	layout string
}{
	{"YYYY", "2006"},
	{"YY", "06"},
	{"MMMM", "January"},
	{"MMM", "Jan"},
	{"MM", "01"},
	{"M", "1"},
	{"DDDD", "Monday"},
	{"DDD", "Mon"},
	{"DD", "02"},
	{"D", "2"},
	{"HH", "15"},
	{"H", "15"},
	{"hh", "03"},
	{"h", "3"},
	{"mm", "04"},
	{"m", "4"},
	{"ss", "05"},
	{"s", "5"},
	{"ZZ", "-0700"},
	{"A", "PM"},
	{"a", "pm"},
	{"Z", "-0700"},
	{"T", "T"},
}

// momentToGoLayout converts a subset of Moment.js style tokens used by Obsidian
// into Go time layouts. Returns false when conversion fails.
func momentToGoLayout(format string) (string, bool) {
	var builder strings.Builder
	runes := []rune(format)

	for i := 0; i < len(runes); {
		switch runes[i] {
		case '[':
			j := i + 1
			for j < len(runes) && runes[j] != ']' {
				j++
			}
			if j >= len(runes) {
				return "", false
			}
			builder.WriteString(string(runes[i+1 : j]))
			i = j + 1
			continue
		case '\\':
			// Escape next rune literally
			if i+1 < len(runes) {
				builder.WriteRune(runes[i+1])
				i += 2
			} else {
				builder.WriteRune(runes[i])
				i++
			}
			continue
		}

		if token, layout, ok := matchToken(runes, i); ok {
			builder.WriteString(layout)
			i += len(token)
			continue
		}

		// Unknown token, write rune verbatim
		builder.WriteRune(runes[i])
		i++
	}

	return builder.String(), true
}

func matchToken(runes []rune, start int) (token, layout string, ok bool) {
	for _, entry := range tokenMap {
		tokenRunes := []rune(entry.token)
		if len(runes)-start < len(tokenRunes) {
			continue
		}
		if equalRunes(runes[start:start+len(tokenRunes)], tokenRunes) {
			token = entry.token
			layout = entry.layout
			ok = true
			return
		}
	}
	return
}

func equalRunes(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
