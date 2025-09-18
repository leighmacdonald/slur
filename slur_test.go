package slur_test

import (
	"fmt"
	"testing"

	"github.com/leighmacdonald/slur"
)

type tc struct {
	line     string
	match    string
	common   string
	category int
	nomatch  bool
}

func TestFind(t *testing.T) {
	cases := []tc{
		{line: "asdf apes adf", match: "apes", common: "apes", category: 10},
		{line: "ashke-nazi asdf asdf ", match: "ashke-nazi", common: "jew", category: 10},
		{line: "a sdf adsf ashkenazi asdf", match: "ashkenazi", common: "jew", category: 10},
		{line: "a sdf adsf jiggerboo asdf", match: "jiggerboo", common: "nigger", category: 50},
		{line: "a sdf adsf faggggggg asdf", match: "faggggggg", common: "fag", category: 50},
		{line: "a sdf adsf ywnbaw asdf", match: "ywnbaw", common: "tranny", category: 50},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("case:%d", index), func(t *testing.T) {
			result, matched := slur.Check(tc.line)
			if tc.nomatch && matched {
				t.Fatalf("expected no match, got match")
			} else if !matched {
				t.Fatalf("expected match, got no match: %s", tc.line)
			}
			if result.Category != tc.category {
				t.Fatalf("expected category %d, got %d (%s)", tc.category, result.Category, tc.line)
			}
			if result.Word != tc.match {
				t.Fatalf("expected word %s, got %s", tc.match, result.Word)
			}
			if result.Common != tc.common {
				t.Fatalf("expected word %s, got %s", tc.common, result.Common)
			}
		})
	}
}

func ExampleCheck() {
	var check slur.Checker
	check = slur.NewString(10, "example")

	if match, found := check.Check("this is an example"); found {
		fmt.Printf("%s %d\n", match.Word, match.Category)
		// Output: example 10
	}
}
