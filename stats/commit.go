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

type CommitStatsIterator struct {
	stats []CommitStats

	// Passed to next invocation
	cnt  int
	tags []string
	id   *git.Oid

	diffOpts        git.DiffOptions
	diffFindOpts    git.DiffFindOptions
	taggedCommitMap map[string][]string
	repo            *git.Repository
}

func (csi *CommitStatsIterator) cb(c *git.Commit) bool {
	if tags, ok := csi.taggedCommitMap[c.Id().String()]; ok {
		if csi.id == nil {
			csi.id = c.Id()
			csi.tags = tags
			csi.cnt = 0
			return true
		}

		childCommit, _ := csi.repo.LookupCommit(csi.id)
		ins, del, _ := getInsAndDel(csi.repo, c, childCommit, &csi.diffOpts, &csi.diffFindOpts)

		csi.stats = append(csi.stats, CommitStats{
			Tags: csi.tags,
			Oid:  csi.id,
			Cnt:  csi.cnt,
			Ins:  ins,
			Del:  del,
		})
		csi.id = c.Id()
		csi.tags = tags
		csi.cnt = 0
		return true
	}
	csi.cnt++
	return true
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
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	defaultOpts, _ := git.DefaultDiffOptions()
	defaultDiffOpts, _ := git.DefaultDiffFindOptions()

	csi := &CommitStatsIterator{
		diffOpts:        defaultOpts,
		diffFindOpts:    defaultDiffOpts,
		taggedCommitMap: taggedCommitMap,
		repo:            repo,
	}
	err = walk.Iterate(csi.cb)

	return csi.stats, nil
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
