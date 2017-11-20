package git

import (
	"fmt"
	"strings"

	"github.com/libgit2/git2go"
)

type Analyzer interface {
	Analyze()
}

type GitTag struct {
	path string
}

func (g *GitTag) Analyze(str string) ([]TagCount, error) {
	repo, _ := git.OpenRepository(g.path)
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, str) {
			o, _ := repo.Lookup(oid)
			t, _ := o.AsTag()
			c, _ := t.AsCommit()
			fmt.Println(c)
			//fmt.Println(c.Summary())
		}
		return nil
	})

	/*
		callback := func(c *git.Commit) bool {
			return true
		}

		err = walk.Iterate(callback)
		if err != nil {
			fmt.Printf(err.Error())
		}
	*/

	return nil, nil
}
