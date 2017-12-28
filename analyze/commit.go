package analyze

import (
	"strings"

	git "github.com/libgit2/git2go"
)

type TaggedCommit struct {
	Tags []string
	Oid  *git.Oid
	next *git.Oid
	prev *git.Oid
	Cnt  int
	Name string
}

func CountCommit(repo *git.Repository, str string) ([]TaggedCommit, error) {
	tcs := []TaggedCommit{}
	commitWithTagMap := map[string][]string{}

	// mutable
	tag_p := &git.Oid{}
	prev_tags := []string{}
	next := &git.Oid{}
	cnt := 0
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, str) {
			o, _ := repo.Lookup(oid)
			switch o.Type() {
			case git.ObjectTag: // For annotated tag
				tag, _ := o.AsTag()
				revision := tag.TargetId().String()
				commitWithTagMap[revision] = append(commitWithTagMap[revision], name)
			case git.ObjectCommit: // For lightweight tag
				revision := o.Id().String()
				commitWithTagMap[revision] = append(commitWithTagMap[revision], name)
			}
		}
		return nil
	})

	setFirst := func(c *git.Commit) bool {
		if tags, ok := commitWithTagMap[c.Id().String()]; ok {
			tag_p = c.Id()
			prev_tags = tags
			return false
		}
		return true
	}

	walk.Iterate(setFirst)

	construct := func(c *git.Commit) bool {
		if tags, ok := commitWithTagMap[c.Id().String()]; ok {
			tcs = append(tcs, TaggedCommit{
				Tags: prev_tags,
				Oid:  tag_p,
				next: next,
				prev: c.Id(),
				Cnt:  cnt},
			)
			next = tag_p
			tag_p = c.Id()
			prev_tags = tags
			cnt = 0
		} else {
			cnt++
		}
		return true
	}

	err = walk.Iterate(construct)
	tcs = append(tcs, TaggedCommit{
		Tags: prev_tags,
		Oid:  tag_p,
		next: next,
		prev: nil,
		Cnt:  cnt},
	)

	if err != nil {
		return nil, nil
	}

	return tcs, nil
}
