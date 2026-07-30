package pkg

import "strings"

// normalizeFillerWords are dropped when normalizing a phrase for cache lookup so
// that "had lunch", "a lunch", and "lunch" collapse to the same cache key.
var normalizeFillerWords = map[string]bool{
	"a": true, "an": true, "the": true, "some": true,
	"had": true, "have": true, "my": true, "of": true,
}

// NormalizePhrase lower-cases, strips punctuation and filler words, and collapses
// whitespace so equivalent inputs share one cache key (raising the AI cache hit
// rate). It is the single source of truth for cache keys: the parser normalizes
// before every lookup, and the AI-cache writers normalize before every store, so
// seeded and learned entries always align. It is idempotent.
func NormalizePhrase(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == ' ':
			b.WriteRune(r)
		case r >= 0x0980 && r <= 0x09FF: // Bengali block — keep for Banglish inputs
			b.WriteRune(r)
		default:
			b.WriteRune(' ')
		}
	}
	fields := make([]string, 0)
	for _, w := range strings.Fields(b.String()) {
		if normalizeFillerWords[w] {
			continue
		}
		fields = append(fields, w)
	}
	return strings.Join(fields, " ")
}
