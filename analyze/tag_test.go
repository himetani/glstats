package analyze

import (
	"reflect"
	"testing"
	"time"

	git "github.com/libgit2/git2go"
)

func TestTagAnalyze(t *testing.T) {

	repoPath := "../glstats-sample-submodule"
	repo, _ := git.OpenRepository(repoPath)

	t0 := time.Date(2017, time.Month(9), 1, 0, 0, 0, 0, time.Local)
	t1 := time.Date(2017, time.Month(10), 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(2017, time.Month(11), 1, 0, 0, 0, 0, time.Local)

	times := []time.Time{t0, t1, t2}

	expected := []Tag{
		{Time: t0, Cnt: 0},
		{Time: t1, Cnt: 0},
		{Time: t2, Cnt: 3},
	}

	actual, err := CountTag(repo, "deploy", times)
	if err != nil {
		t.Fatal("Analyze returnen non-nil\n")
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual was %x\n", actual)
	}
}
