package repo

import (
	"strings"

	git "github.com/libgit2/git2go"
)

type CommitWithLine struct {
	Tags []string
	Oid  *git.Oid
	next *git.Oid
	prev *git.Oid
	Ins  int
	Del  int
	Cnt  int
}

func CountLine(repo *git.Repository, str string) ([]CommitWithLine, error) {
	cwl := []CommitWithLine{}
	commitWithTagMap := map[string][]string{}

	// mutable
	var tag_p *git.Oid
	prev_tags := []string{}
	next := &git.Oid{}
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		shortName := strings.Replace(name, "refs/tags/", "", -1)
		if strings.Contains(name, str) {
			o, _ := repo.Lookup(oid)
			switch o.Type() {
			case git.ObjectTag: // For annotated tag
				tag, _ := o.AsTag()
				revision := tag.TargetId().String()
				commitWithTagMap[revision] = append(commitWithTagMap[revision], shortName)
			case git.ObjectCommit: // For lightweight tag
				revision := o.Id().String()
				commitWithTagMap[revision] = append(commitWithTagMap[revision], shortName)
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

	defaultOpts, _ := git.DefaultDiffOptions()
	defaultDiffOpts, _ := git.DefaultDiffFindOptions()

	construct := func(c *git.Commit) bool {
		if tags, ok := commitWithTagMap[c.Id().String()]; ok {
			childCommit, _ := repo.LookupCommit(tag_p)
			ins, del, _ := getInsAndDel(repo, c, childCommit, &defaultOpts, &defaultDiffOpts)
			cwl = append(cwl, CommitWithLine{
				Tags: prev_tags,
				Oid:  tag_p,
				next: next,
				prev: c.Id(),
				Ins:  ins,
				Del:  del,
			})

			next = tag_p
			tag_p = c.Id()
			prev_tags = tags
		}

		return true
	}

	err = walk.Iterate(construct)
	cwl = append(cwl, CommitWithLine{
		Tags: prev_tags,
		Oid:  tag_p,
		next: next,
		prev: nil,
	})

	if err != nil {
		return nil, nil
	}

	return cwl, nil
}

func getInsAndDel(r *git.Repository, o, c *git.Commit, opts *git.DiffOptions, diffOpts *git.DiffFindOptions) (int, int, error) {
	if o == nil {
		return 0, 0, nil
	}
	tree, _ := c.Tree()
	oldTree, _ := o.Tree()

	diff, _ := r.DiffTreeToTree(oldTree, tree, opts)
	diff.FindSimilar(diffOpts)
	stats, _ := diff.Stats()
	return stats.Insertions(), stats.Deletions(), nil
}
