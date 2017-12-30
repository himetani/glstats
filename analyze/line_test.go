package analyze

import (
	"testing"

	git "github.com/libgit2/git2go"
)

func TestHoge(t *testing.T) {
	repoPath := "/Users/takafumi.tsukamoto/dev/src/github.com/himetani/glstats-sample"
	repo, _ := git.OpenRepository(repoPath)

	Hoge(repo)

	f := func() (int, int) {
		return 3, 4
	}

	closure(f)
}
