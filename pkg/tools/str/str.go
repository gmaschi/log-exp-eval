package str

import "unicode"

func HasLowers(s string) ([]string, bool) {
	var hasLowers bool
	lowers := make([]string, 0, len(s))

	for _, r := range s {
		if unicode.IsLetter(r) && unicode.IsLower(r) {
			lowers = append(lowers, string(r))
			hasLowers = true
		}
	}

	return lowers, hasLowers
}
