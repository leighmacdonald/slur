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
		NewString(50, "nigger", "nig", "ni"),
		NewString(50, "fag", "faggot"),
		NewString(10, "apes"),
		NewString(50, "tranny", "shemale", "ywnbaw", "troon"),
		NewString(10, "spic", "wetback"),
		Regex{regexp.MustCompile(`^apes?$`), 10},       //Referring to outdated theories ascribing cultural differences between racial groups as being linked to their evolutionary distance from chimpanzees, with which humans share common ancestry.
		Regex{regexp.MustCompile(`^ashke-?nazi$`), 75}, // Pronounced like "AshkeNatzi". Used mostly by Mizrachi Jews.
		Regex{regexp.MustCompile(`^beaners?$`), 30},
		Regex{regexp.MustCompile(`^japs?$`), 30},
		Regex{regexp.MustCompile(`^jiggaboo|jiggerboo|niggerboo|jiggabo|jigarooni|jijjiboo|zigab|jig|jigg|jigger$`), 50},
		Regex{regexp.MustCompile(`^nigger|niger|nig|nigga|niggress$`), 50},
		Regex{regexp.MustCompile(`^k[iy]ke$`), 50},
		Regex{regexp.MustCompile(`^yid$`), 30},
		Regex{regexp.MustCompile(`^paki$`), 30},
		Regex{regexp.MustCompile(`(fa[g6]+|f[a4][g6]{1,}[o0][t7])s?$`), 50},
		Regex{regexp.MustCompile(`^chinks?$`), 50},
		Regex{regexp.MustCompile(`^tr[a4]nn?y?$`), 75},
		Regex{regexp.MustCompile(`^ni$`), 30},
		Regex{regexp.MustCompile(`^gooks?$`), 30},
		Regex{regexp.MustCompile(`^(spic|wetback)s$`), 30},
		Regex{regexp.MustCompile(`^homos?$`), 30},
		Regex{regexp.MustCompile(`^wops?$`), 30},
		Regex{regexp.MustCompile(`^coons?$`), 50},
		Regex{regexp.MustCompile(`^(shemale|troon)s?$`), 50},
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
	pattern *regexp.Regexp
	weight  int
}

func (m Regex) Check(word string) (Match, bool) {
	match := m.pattern.FindStringSubmatch(word)
	if match != nil {
		return Match{Word: match[0], Weight: m.weight}, true
	}

	return Match{}, false
}

func NewString(weight int, patterns ...string) *String {
	return &String{
		patterns: patterns,
		weight:   weight,
	}
}

type String struct {
	patterns []string
	weight   int
}

func (m String) Check(word string) (Match, bool) {
	if slices.Contains(m.patterns, word) {
		return Match{Word: word, Weight: m.weight}, true
	}

	return Match{}, false
}

type Match struct {
	// Word contains the matched word
	Word string
	// Weight is a 0-100 score indicating the severity of the slur, 0 being the least severe.
	Weight int
}

func normalize(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(strings.ToLower(text))), " ")
}
