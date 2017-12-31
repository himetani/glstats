package stats

import (
	"strings"

	git "github.com/libgit2/git2go"
)

type CommitStats struct {
	Tags []string
	Oid  *git.Oid
	Cnt  int // Number of commit
	Ins  int
	Del  int
}

func CountCommit(repo *git.Repository, tagSubStr string) (map[string][]string, error) {
	taggedCommitMap := map[string][]string{}

	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	repo.Tags.Foreach(func(name string, oid *git.Oid) error {
		short := strings.Replace(name, "refs/tags/", "", -1)
		if strings.Contains(short, tagSubStr) {
			obj, _ := repo.Lookup(oid)
			switch obj.Type() {
			case git.ObjectTag: // For annotated tag
				tag, _ := obj.AsTag()
				revision := tag.TargetId().String()
				taggedCommitMap[revision] = append(taggedCommitMap[revision], short)
			case git.ObjectCommit: // For lightweight tag
				revision := obj.Id().String()
				taggedCommitMap[revision] = append(taggedCommitMap[revision], short)
			}
		}
		return nil
	})
	return taggedCommitMap, nil
}

func GetStats(repo *git.Repository, taggedCommitMap map[string][]string) ([]CommitStats, error) {
	// mutable
	commits := []CommitStats{}
	tag_p := &git.Oid{}
	prev_tags := []string{}
	next := &git.Oid{}
	cnt := 0

	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	setFirst := func(c *git.Commit) bool {
		if tags, ok := taggedCommitMap[c.Id().String()]; ok {
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
		if tags, ok := taggedCommitMap[c.Id().String()]; ok {
			childCommit, _ := repo.LookupCommit(tag_p)
			ins, del, _ := getInsAndDel(repo, c, childCommit, &defaultOpts, &defaultDiffOpts)

			commits = append(commits, CommitStats{
				Tags: prev_tags,
				Oid:  tag_p,
				Cnt:  cnt,
				Ins:  ins,
				Del:  del,
			})
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
	commits = append(commits, CommitStats{
		Tags: prev_tags,
		Oid:  tag_p,
		Cnt:  cnt},
	)

	if err != nil {
		return nil, nil
	}

	return commits, nil
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
