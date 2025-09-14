package slur

import (
	"regexp"
	"strings"
)

var (
	slurs = []definition{}
)

type definition struct {
	pattern *regexp.Regexp
	weight  int
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

func Find(text string) (Match, bool) {
	for word := range strings.SplitSeq(normalize(text), " ") {
		for _, slur := range slurs {
			match := slur.pattern.FindStringSubmatch(word)
			if match != nil {
				return Match{Word: match[0], Weight: slur.weight}, true
			}
		}
	}
	return Match{}, false
}

func init() {
	slurs = []definition{
		{regexp.MustCompile(`^apes?$`), 10},       //Referring to outdated theories ascribing cultural differences between racial groups as being linked to their evolutionary distance from chimpanzees, with which humans share common ancestry.
		{regexp.MustCompile(`^ashke-?nazi$`), 75}, // Pronounced like "AshkeNatzi". Used mostly by Mizrachi Jews.
		{regexp.MustCompile(`^beaners?$`), 30},
		{regexp.MustCompile(`^japs?$`), 30},
		{regexp.MustCompile(`^jiggaboo|jiggerboo|niggerboo|jiggabo|jigarooni|jijjiboo|zigab|jig|jigg|jigger$`), 50},
		{regexp.MustCompile(`^nigger|niger|nig|nigga|niggress$`), 50},
		{regexp.MustCompile(`^k[iy]ke$`), 50},
		{regexp.MustCompile(`^yid$`), 30},
		{regexp.MustCompile(`^paki$`), 30},
		{regexp.MustCompile(`(fa[g6]+|f[a4][g6]{1,}[o0][t7])s?$`), 50},
		{regexp.MustCompile(`^chinks?$`), 50},
		{regexp.MustCompile(`^tr[a4]nn?y?$`), 75},
		{regexp.MustCompile(`^ni$`), 30},
		{regexp.MustCompile(`^gooks?$`), 30},
		{regexp.MustCompile(`^(spic|wetback)s$`), 30},
		{regexp.MustCompile(`^homos?$`), 30},
		{regexp.MustCompile(`^wops?$`), 30},
		{regexp.MustCompile(`^coons?$`), 50},
		{regexp.MustCompile(`^(shmemale|troon)s?$`), 50},
		{regexp.MustCompile(`^ywnbaw$`), 50},
		{regexp.MustCompile(`^(retards?|retarded)$`), 25},
	}
}
