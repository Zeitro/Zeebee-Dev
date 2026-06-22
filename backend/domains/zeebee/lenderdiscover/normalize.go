package lenderdiscover

import (
	"regexp"
	"strings"
)

var nonAlnumSpace = regexp.MustCompile(`[^a-z0-9\s]+`)

var companySuffixes = map[string]struct{}{
	"inc": {}, "incorporated": {}, "llc": {}, "ltd": {}, "limited": {},
	"corp": {}, "corporation": {}, "company": {}, "co": {}, "llp": {},
	"plc": {}, "dba": {},
}

// NormalizeCompanyName lowercases, removes punctuation, and strips common legal suffix tokens
// so that "Acme, Inc." and "ACME LLC" match.
func NormalizeCompanyName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlnumSpace.ReplaceAllString(s, " ")
	fields := strings.Fields(s)
	var kept []string
	for _, w := range fields {
		if _, skip := companySuffixes[w]; skip {
			continue
		}
		kept = append(kept, w)
	}
	return strings.Join(kept, " ")
}
