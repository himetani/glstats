package stats

import (
	"reflect"
	"testing"

	git "github.com/libgit2/git2go"
)

var (
	all = map[string][]string{
		"01bb16f3d083fa252c4476f6419a4b2761f4a839": []string{"deploy/20171121"},
		"264f0767fb0cb4f34eb49d63022a443cefb75783": []string{"deploy/20171123"},
		"e7cec8f34445ef794e54c0d7f6bacef97d99bf5a": []string{"deploy/20171124"},
		"0c254ca47f0924a4d0874c88499315a0987d8d3f": []string{"deploy/20180102"},
		"65d04b818726554c198866e5f9fbef65d6064a46": []string{"deploy/20180106"},
	}

	deploy2017 = map[string][]string{
		"264f0767fb0cb4f34eb49d63022a443cefb75783": []string{"deploy/20171123"},
		"e7cec8f34445ef794e54c0d7f6bacef97d99bf5a": []string{"deploy/20171124"},
		"01bb16f3d083fa252c4476f6419a4b2761f4a839": []string{"deploy/20171121"},
	}

	deploy2018 = map[string][]string{
		"0c254ca47f0924a4d0874c88499315a0987d8d3f": []string{"deploy/20180102"},
		"65d04b818726554c198866e5f9fbef65d6064a46": []string{"deploy/20180106"},
	}
)

func TestGetTaggedCommitMap(t *testing.T) {
	type data struct {
		TestName  string
		TagSubstr string
		Result    map[string][]string
	}

	repo, err := git.OpenRepository("../glstats-sample-submodule")
	if err != nil {
		t.Errorf("Before Test: unexpected error %s", err)
	}

	tests := []data{
		{"All", "", all},
		{"Only 2017", "deploy/2017", deploy2017},
		{"Only 2018", "deploy/2018", deploy2018},
	}

	for i, test := range tests {
		result, err := GetTaggedCommitMap(repo, test.TagSubstr)
		if err != nil {
			t.Errorf("Test #%d %s: unexpected error %s", i, test.TestName, err)
		}

		if !reflect.DeepEqual(test.Result, result) {
			t.Errorf("Test #%d %s: expected %q, got %q", i, test.TestName, test.Result, result)
		}
	}

}

func TestGetCommitStats(t *testing.T) {
	type data struct {
		TestName        string
		TaggedCommitMap map[string][]string
		Result          []CommitStats
	}

	deploy20180106 := CommitStats{
		Tags:     []string{"deploy/20180106"},
		Revision: "65d04b818726554c198866e5f9fbef65d6064a46",
		Cnt:      1,
		Ins:      0,
		Del:      3,
	}

	deploy20180102 := CommitStats{
		Tags:     []string{"deploy/20180102"},
		Revision: "0c254ca47f0924a4d0874c88499315a0987d8d3f",
		Cnt:      1,
		Ins:      0,
		Del:      0,
	}

	deploy20180102_2 := CommitStats{
		Tags:     []string{"deploy/20180102"},
		Revision: "0c254ca47f0924a4d0874c88499315a0987d8d3f",
		Cnt:      6,
		Ins:      9,
		Del:      0,
	}

	deploy20171124 := CommitStats{
		Tags:     []string{"deploy/20171124"},
		Revision: "e7cec8f34445ef794e54c0d7f6bacef97d99bf5a",
		Cnt:      2,
		Ins:      7,
		Del:      0,
	}

	deploy20171123 := CommitStats{
		Tags:     []string{"deploy/20171123"},
		Revision: "264f0767fb0cb4f34eb49d63022a443cefb75783",
		Cnt:      1,
		Ins:      2,
		Del:      0,
	}

	repo, err := git.OpenRepository("../glstats-sample-submodule")
	if err != nil {
		t.Errorf("Before Test: unexpected error %s", err)
	}

	tests := []data{
		{"All", all, []CommitStats{deploy20180106, deploy20180102, deploy20171124, deploy20171123}},
		{"Only 2017", deploy2017, []CommitStats{deploy20171124, deploy20171123}},
		{"Only 2018", deploy2018, []CommitStats{deploy20180106, deploy20180102_2}},
	}

	for i, test := range tests {
		result, err := GetStats(repo, test.TaggedCommitMap)
		if err != nil {
			t.Errorf("Test #%d %s: unexpected error %s", i, test.TestName, err)
		}

		if len(result) != len(test.Result) {
			t.Errorf("Test #%d %s: expected %d, got %d", i, test.TestName, len(test.Result), len(result))
		}

		if !reflect.DeepEqual(test.Result, result) {
			t.Errorf("Test #%d %s: expected %q, got %q", i, test.TestName, test.Result, result)
		}
	}

}
