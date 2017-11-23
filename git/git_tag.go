package git

import (
	"strings"
	"time"

	"github.com/libgit2/git2go"
)

type Analyzer interface {
	Analyze()
}

type GitTag struct {
	path string
}

func (g *GitTag) Analyze(str string, times []time.Time) ([]TagCount, error) {
	repo, _ := git.OpenRepository(g.path)
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	tags := []time.Time{}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, str) {
			o, _ := repo.Lookup(oid)
			tag, _ := o.AsTag()
			tTime := tag.Tagger().When
			tags = append(tags, tTime)
		}
		return nil
	})

	tcs := []TagCount{}
	for _, time := range times {
		cnt := 0
		for _, tag := range tags {
			if tag.After(time) && time.AddDate(0, 1, 0).After(tag) {
				cnt++
			}
		}
		tcs = append(tcs, TagCount{time: time, cnt: cnt})
	}

	/*
		callback := func(c *git.Commit) bool {
			return true
		}

		err = walk.Iterate(callback)
		if err != nil {
			fmt.Printf(err.Error())
		}
	*/

	return tcs, nil
}
