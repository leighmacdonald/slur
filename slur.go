// Package slur is a very simple package for checking text for slurs. It exists to provide a
// shared implementation of the same functionality for a couple parent packages.
package slur

import (
	"regexp"
	"slices"
	"strings"
)

var (
	Checkers = []Checker{
		NewString(50, "nigger", "nig", "ni", "darkie", "tar-baby", "tarbaby", "coon"),
		NewString(50, "fag", "faggot"),
		NewString(10, "apes"),
		NewString(50, "tranny", "shemale", "ywnbaw", "troon", "transtrender"),
		NewString(10, "spic", "wetback", "spik", "spick"),
		NewString(10, "chink", "gook"),
		NewString(10, "half-breed", "halfbreed"),
		NewString(10, "injun"),
		NewString(10, "jap", "nip", "slant–eye", "slanteye", "zipperhead"),
		NewString(10, "jewboy"),
		NewString(10, "卍", "卐"),
		Regex{regexp.MustCompile(`^ashke-?nazi$`), 75}, // Pronounced like "AshkeNatzi". Used mostly by Mizrachi Jews.
		Regex{regexp.MustCompile(`^beaners?$`), 30},
		Regex{regexp.MustCompile(`^jiggaboo|jiggerboo|niggerboo|jiggabo|jigarooni|jijjiboo|zigab|jig|jigg|jigger$`), 50},
		Regex{regexp.MustCompile(`nigger|nigga|niggress$`), 50},
		Regex{regexp.MustCompile(`^k[iy]ke$`), 50},
		Regex{regexp.MustCompile(`^yid$`), 30},
		Regex{regexp.MustCompile(`^paki$`), 30},
		Regex{regexp.MustCompile(`(fa[g6]+|f[a4][g6]{1,}[o0][t7])s?$`), 50},
		Regex{regexp.MustCompile(`^tr[a4]nn?y?$`), 75},
		Regex{regexp.MustCompile(`^(retards?|retarded)$`), 25},
	}
)

type Checker interface {
	Check(text string) (Match, bool)
}

// Check a line of text for slurs. Returns the first match found, and a true value. Otherwise, returns an empty Match and false.
// Checks each word individually, for example, `^` would match the beginning of the word for regex, not the entire line.
func Check(line string) (Match, bool) {
	for word := range strings.SplitSeq(normalize(line), " ") {
		for _, slur := range Checkers {
			if match, found := slur.Check(word); found {
				return match, true
			}
		}
	}

	return Match{}, false
}

type Regex struct {
	pattern  *regexp.Regexp
	category int
}

func (m Regex) Check(word string) (Match, bool) {
	match := m.pattern.FindStringSubmatch(word)
	if match != nil {
		return Match{Word: match[0], Category: m.category}, true
	}

	return Match{}, false
}

func NewString(category int, patterns ...string) *String {
	return &String{
		patterns: patterns,
		category: category,
	}
}

// String handles performing simple string matches.
//
// TODO (maybe): Add levenshtein distance option
type String struct {
	patterns []string
	category int
}

func (m String) Check(word string) (Match, bool) {
	if slices.Contains(m.patterns, word) {
		return Match{Word: word, Category: m.category}, true
	}

	return Match{}, false
}

type Match struct {
	// Word contains the matched word
	Word string
	// Category is a value that can be used to group similar terms.
	Category int
}

func normalize(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(strings.ToLower(text))), " ")
}
