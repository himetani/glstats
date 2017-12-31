package repo

import (
	"strings"
	"time"

	git "github.com/libgit2/git2go"
)

type Tag struct {
	Time time.Time
	Cnt  int
}

const (
	layout string = "200601021504"
)

func CountTagBy(repo *git.Repository, tagSubstr string, times []time.Time) ([]Tag, error) {
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	timestamps := getTagTimestamps(repo, tagSubstr)

	counts := []Tag{}
	for _, time := range times {
		cnt := 0
		for _, tag := range timestamps {
			if tag.After(time) && time.AddDate(0, 1, 0).After(tag) {
				cnt++
			}
		}
		counts = append(counts, Tag{Time: time, Cnt: cnt})
	}

	return counts, nil
}

func getTagTimestamps(repo *git.Repository, tagSubstr string) []time.Time {
	var timestamps []time.Time

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, tagSubstr) {
			var t time.Time

			o, _ := repo.Lookup(oid)
			switch o.Type() {
			case git.ObjectTag: // For annotated tag
				tag, _ := o.AsTag()
				t = tag.Tagger().When
			case git.ObjectCommit: // For lightweight tag
				tstr := strings.Replace(name, "refs/tags/"+tagSubstr+"/", "", -1)
				t, _ = time.Parse(layout, tstr)
			}
			timestamps = append(timestamps, t)
		}
		return nil
	})

	return timestamps
}
