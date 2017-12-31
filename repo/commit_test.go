package repo

import (
	"testing"

	git "github.com/libgit2/git2go"
)

func TestCommitCount(t *testing.T) {
	repo, _ := git.OpenRepository("../glstats-sample-submodule")

	expected := []struct {
		revision string
		cnt      int
	}{
		{revision: "e7cec8f34445ef794e54c0d7f6bacef97d99bf5a", cnt: 2},
		{revision: "264f0767fb0cb4f34eb49d63022a443cefb75783", cnt: 1},
		{revision: "01bb16f3d083fa252c4476f6419a4b2761f4a839", cnt: 0},
	}
	taggedCommits, err := CountCommit(repo, "deploy")
	if err != nil {
		t.Fatal("Analyze return non-nil\n")
	}

	for i, tc := range taggedCommits {
		if tc.Oid.String() != expected[i].revision || tc.Cnt != expected[i].cnt {
			t.Fatalf("exected was %x, but was %x\n", expected[i], tc)
		}
	}

}
