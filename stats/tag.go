package stats

import (
	"strings"
	"time"

	git "github.com/libgit2/git2go"
)

// Tag is struct hold stats data of Tag
type Tag struct {
	Time time.Time
	Cnt  int
}

const (
	layout string = "200601021504"
)

type tagTimestampIterator struct {
	repo       *git.Repository
	timestamps []time.Time
	tagSubstr  string
}

func (ti *tagTimestampIterator) cb(name string, oid *git.Oid) error {
	if strings.Contains(name, ti.tagSubstr) {
		var t time.Time
		o, _ := ti.repo.Lookup(oid)
		switch o.Type() {
		// For annotated tag
		case git.ObjectTag:
			tag, _ := o.AsTag()
			t = tag.Tagger().When
		// For lightweight tag
		case git.ObjectCommit:
			commit, _ := o.AsCommit()
			t = commit.Committer().When
		}
		ti.timestamps = append(ti.timestamps, t)
	}
	return nil
}

// CountTagBy returns slice of Tag
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
	itr := &tagTimestampIterator{
		repo:      repo,
		tagSubstr: tagSubstr,
	}
	repo.Tags.Foreach(itr.cb)
	return itr.timestamps
}
