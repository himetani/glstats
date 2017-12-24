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

const (
	layout string = "200601021504"
)

func CountTag(repo *git.Repository, substr string, times []time.Time) ([]TagCount, error) {
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	timestamps := getTimestamps(repo, substr)

	counts := []TagCount{}
	for _, time := range times {
		cnt := 0
		for _, tag := range timestamps {
			if tag.After(time) && time.AddDate(0, 1, 0).After(tag) {
				cnt++
			}
		}
		counts = append(counts, TagCount{Time: time, Cnt: cnt})
	}

	return counts, nil
}

func getTimestamps(repo *git.Repository, substr string) []time.Time {
	var timestamps []time.Time

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, substr) {
			var t time.Time

			o, _ := repo.Lookup(oid)
			switch o.Type() {
			case git.ObjectTag: // For annotated tag
				tag, _ := o.AsTag()
				t = tag.Tagger().When
			case git.ObjectCommit: // For lightweight tag
				tstr := strings.Replace(name, "refs/tags/"+substr+"/", "", -1)
				t, _ = time.Parse(layout, tstr)
			}
			timestamps = append(timestamps, t)
		}
		return nil
	})

	return timestamps
}
