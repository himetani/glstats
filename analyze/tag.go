package analyze

import (
	"strings"
	"time"

	git "github.com/libgit2/git2go"
)

type TagCount struct {
	Time time.Time
	Cnt  int
}

type TaggedCommit struct {
	tags      []*git.Oid
	oid       *git.Oid
	next      *git.Oid
	prev      *git.Oid
	commitCnt int
}

func CountTag(repo *git.Repository, str string, times []time.Time) ([]TagCount, error) {
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	tagTimes := []time.Time{}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, str) {
			var tTime time.Time
			layout := "200601021504"

			o, _ := repo.Lookup(oid)
			switch o.Type() {
			// For annotated tag
			case git.ObjectTag:
				tag, _ := o.AsTag()
				tTime = tag.Tagger().When
			// For lightweight tag
			case git.ObjectCommit:
				tstr := strings.Replace(name, "refs/tags/"+str+"/", "", -1)
				tTime, err = time.Parse(layout, tstr)
				if err != nil {
					return nil
				}
			}
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
