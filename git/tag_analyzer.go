package git

import (
	"strings"
	"time"

	"github.com/libgit2/git2go"
)

type TagAnalyzer struct {
	Path string
}

func (g *TagAnalyzer) Analyze(str string, times []time.Time) ([]TagCount, error) {
	repo, _ := git.OpenRepository(g.Path)
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	tagTimes := []time.Time{}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, str) {
			o, _ := repo.Lookup(oid)
			tag, _ := o.AsTag()
			tTime := tag.Tagger().When
			tagTimes = append(tagTimes, tTime)
		}
		return nil
	})

	tcs := []TagCount{}
	for _, time := range times {
		cnt := 0
		for _, tag := range tagTimes {
			if tag.After(time) && time.AddDate(0, 1, 0).After(tag) {
				cnt++
			}
		}
		tcs = append(tcs, TagCount{Time: time, Cnt: cnt})
	}

	return tcs, nil
}
