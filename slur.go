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
		NewString(50, "nigger", "nigga", "nig", "ni", "niggress", "darkie", "tar-baby", "tarbaby", "coon", "jiggaboo", "jiggerboo", "niggerboo", "jiggabo", "jigarooni", "jijjiboo", "zigab", "jig", "jigg", "jigger"),
		NewString(50, "fag", "faggot", "faggy"),
		NewString(10, "apes"),
		NewString(50, "tranny", "shemale", "ywnbaw", "troon", "transtrender", "tr4nny"),
		NewString(10, "spic", "wetback", "spik", "spick", "beaner"),
		NewString(10, "chink", "gook"),
		NewString(10, "half-breed", "halfbreed"),
		NewString(10, "injun"),
		NewString(10, `paki`),
		NewString(10, "jap", "nip", "slant–eye", "slanteye", "zipperhead"),
		NewString(10, "jew", "jewboy", "yid", "ashke-nazi", "ashkenazi", "kyke", "kike"),
		NewString(10, "卍", "卐"),
		NewRegex(50, "fag", `(fa[g6]+|f[a4][g6]{1,}[o0][t7])s?$`, "fag$"),
		NewRegex(75, "tranny", `^tr[a4]nn?y?$`),
		NewRegex(25, "retard", `^(retards?|retarded)$`),
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
func NewRegex(category int, common string, patterns ...string) Regex {
	matcher := Regex{common: common, category: category}
	for _, pattern := range patterns {
		matcher.pattern = append(matcher.pattern, regexp.MustCompile(pattern))
	}
	return matcher
}

type Regex struct {
	common   string
	pattern  []*regexp.Regexp
	category int
}

func (m Regex) Check(word string) (Match, bool) {
	for _, pattern := range m.pattern {
		match := pattern.FindStringSubmatch(word)
		if match != nil {
			return Match{Word: match[0], Category: m.category, Common: m.common}, true
		}
	}
	return Match{}, false
}

func NewString(category int, common string, patterns ...string) *String {
	return &String{
		common:   common,
		patterns: append([]string{common}, patterns...),
		category: category,
	}
}

// String handles performing simple string matches.
//
// TODO (maybe): Add levenshtein distance option
type String struct {
	common   string
	patterns []string
	category int
}

func (m String) Check(word string) (Match, bool) {
	if slices.Contains(m.patterns, word) {
		return Match{Word: word, Category: m.category, Common: m.common}, true
	}

	return Match{}, false
}

type Match struct {
	// Common or simplified representation of the word
	Common string
	// Word contains the matched word
	Word string
	// Category is a value that can be used to group similar terms.
	Category int
}

func normalize(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(strings.ToLower(text))), " ")
}
