package analyze

import (
	"strings"

	"github.com/libgit2/git2go"
)

type CommitAnalyzer struct {
	path string
}

func (g *CommitAnalyzer) Analyze(str string) ([]TaggedCommit, error) {
	repo, _ := git.OpenRepository(g.path)

	commitMap := map[string][]*git.Oid{}
	tcs := []TaggedCommit{}

	// mutable
	tag_p := &git.Oid{}
	prev_tags := []*git.Oid{}
	next := &git.Oid{}
	next = nil
	cnt := 0
	walk, _ := repo.Walk()
	err := walk.PushHead()

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		if strings.Contains(name, str) {
			o, _ := repo.Lookup(oid)
			tag, _ := o.AsTag()
			commitMap[tag.TargetId().String()] = append(commitMap[tag.TargetId().String()], tag.Id())
		}
		return nil
	})

	setFirst := func(c *git.Commit) bool {
		if tags, ok := commitMap[c.Id().String()]; ok {
			tag_p = c.Id()
			prev_tags = tags
			return false
		}
		return true
	}

	err = walk.Iterate(setFirst)

	construct := func(c *git.Commit) bool {
		if tags, ok := commitMap[c.Id().String()]; ok {
			tcs = append(tcs, TaggedCommit{
				tags:      prev_tags,
				oid:       tag_p,
				next:      next,
				prev:      c.Id(),
				commitCnt: cnt},
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
		tags:      prev_tags,
		oid:       tag_p,
		next:      next,
		prev:      nil,
		commitCnt: cnt},
	)

	if err != nil {
		return nil, nil
	}

	return tcs, nil
}
