package pkg

import "strings"

// LevenshteinDistance compares two strings and returns the levenshtein distance between them.
func LevenshteinDistance(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				mn := d[i-1][j]
				if d[i][j-1] < mn {
					mn = d[i][j-1]
				}
				if d[i-1][j-1] < mn {
					mn = d[i-1][j-1]
				}
				d[i][j] = mn + 1
			}
		}

	}
	return d[len(s)][len(t)]
}
