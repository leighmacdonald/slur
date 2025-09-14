package slur_test

import (
	"fmt"
	"testing"

	"github.com/leighmacdonald/slur"
)

type tc struct {
	line    string
	match   string
	weight  int
	nomatch bool
}

func TestFind(t *testing.T) {
	cases := []tc{
		{line: "asdf ape adf", match: "ape", weight: 5},
		{line: "ashke-nazi asdf asdf ", match: "ashke-nazi", weight: 75},
		{line: "a sdf adsf ashkenazi asdf", match: "ashkenazi", weight: 75},
		{line: "a sdf adsf jiggerboo asdf", match: "jiggerboo", weight: 25},
		{line: "a sdf adsf faggggggg asdf", match: "faggggggg", weight: 50},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("case:%d", index), func(t *testing.T) {
			result, matched := slur.Find(tc.line)
			if tc.nomatch && matched {
				t.Fatalf("expected no match, got match")
			} else if !matched {
				t.Fatalf("expected match, got no match")
			}
			if result.Weight != tc.weight {
				t.Fatalf("expected weight %d, got %d", tc.weight, result.Weight)
			}
			if result.Word != tc.match {
				t.Fatalf("expected word %s, got %s", tc.match, result.Word)
			}
		})
	}
}
