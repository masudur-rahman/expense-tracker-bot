package transaction

import (
	"strings"

	"github.com/masudur-rahman/khorcha-pati/models"
	"github.com/masudur-rahman/khorcha-pati/pkg"
)

// normalizePhrase is a thin alias over pkg.NormalizePhrase, the single source of
// truth for cache keys shared with the AI-cache writers.
func normalizePhrase(s string) string {
	return pkg.NormalizePhrase(s)
}

// localClassify matches a phrase against subcategory keywords/names locally, without
// calling the AI. It returns the best subcategory ID when a whole-word keyword hits,
// preferring more specific (longer, multi-word) keywords. This keeps the common case
// off the rate-limited AI endpoint. When locked is true, candidates whose transaction
// types exclude wantType are skipped so an Income verb never matches an expense-only
// subcategory (e.g. "sold rickshaw" must not resolve to trans-taxi).
func localClassify(phrase string, wantType models.TransactionType, locked bool) (string, bool) {
	phrase = normalizePhrase(phrase)
	if phrase == "" {
		return "", false
	}
	padded := " " + phrase + " "

	bestID := ""
	bestScore := 0
	for _, sub := range models.TxnSubcategories {
		if locked && !models.ContainsType(models.SubcategoryTypes[sub.ID], wantType) {
			continue
		}
		for _, kw := range subcategoryKeywords(sub) {
			if !strings.Contains(padded, " "+kw+" ") {
				continue
			}
			// Longer, multi-word keywords are more specific → higher score.
			score := len(strings.Fields(kw))*100 + len(kw)
			if score > bestScore {
				bestScore, bestID = score, sub.ID
			}
		}
	}
	if bestID == "" {
		return "", false
	}
	return bestID, true
}

// subcategoryKeywords returns the normalized whole-word search terms for a subcategory:
// its comma-separated Keywords plus its Name.
func subcategoryKeywords(sub models.TxnSubcategory) []string {
	terms := make([]string, 0)
	for _, kw := range strings.Split(sub.Keywords, ",") {
		if kw = normalizePhrase(kw); kw != "" {
			terms = append(terms, kw)
		}
	}
	if name := normalizePhrase(sub.Name); name != "" {
		terms = append(terms, name)
	}
	return terms
}

// isRateLimitErr reports whether an AI error is a quota/rate-limit failure, so the
// caller can degrade gracefully instead of dropping the user's transaction.
func isRateLimitErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, needle := range []string{"rate limit", "ratelimit", "429", "quota", "resource_exhausted", "too many requests"} {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}
