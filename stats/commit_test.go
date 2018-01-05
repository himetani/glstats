package stats

import (
	"reflect"
	"testing"

	git "github.com/libgit2/git2go"
)

func TestGetTaggedCommitMap(t *testing.T) {
	repo, _ := git.OpenRepository("../glstats-sample-submodule")

	expected := map[string][]string{
		"264f0767fb0cb4f34eb49d63022a443cefb75783": []string{"deploy/20171123"},
		"e7cec8f34445ef794e54c0d7f6bacef97d99bf5a": []string{"deploy/20171124"},
		"0c254ca47f0924a4d0874c88499315a0987d8d3f": []string{"deploy/20180102"},
		"65d04b818726554c198866e5f9fbef65d6064a46": []string{"deploy/20170103"},
		"01bb16f3d083fa252c4476f6419a4b2761f4a839": []string{"deploy/20171121"},
	}

	taggedCommitMap, err := GetTaggedCommitMap(repo, "deploy")
	if err != nil {
		t.Errorf("Before Test: unexpected error %s", err)
	}

	if !reflect.DeepEqual(expected, taggedCommitMap) {
		t.Errorf("expected %q, got %q", expected, taggedCommitMap)
	}
}
