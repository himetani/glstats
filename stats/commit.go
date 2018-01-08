package stats

import (
	"strings"

	git "github.com/libgit2/git2go"
)

// CommitStats is struct having the stats data of the commit
type CommitStats struct {
	Tags     []string
	Revision string
	Cnt      int // Number of commit
	Ins      int
	Del      int
}

type commitStatsIterator struct {
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

func (csi *commitStatsIterator) cb(c *git.Commit) bool {
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
			Tags:     csi.tags,
			Revision: csi.id.String(),
			Cnt:      csi.cnt,
			Ins:      ins,
			Del:      del,
		})
		csi.id = c.Id()
		csi.tags = tags
		csi.cnt = 0
		return true
	}

	if c.Parent(0) == nil {
		childCommit, _ := csi.repo.LookupCommit(csi.id)
		ins, del, _ := getInsAndDel(csi.repo, c, childCommit, &csi.diffOpts, &csi.diffFindOpts)

		csi.stats = append(csi.stats, CommitStats{
			Tags:     csi.tags,
			Revision: csi.id.String(),
			Cnt:      csi.cnt,
			Ins:      ins,
			Del:      del,
		})
	}

	csi.cnt++
	return true
}

type tagIterator struct {
	repo            *git.Repository
	tagSubStr       string
	taggedCommitMap map[string][]string
}

func (ti *tagIterator) cb(name string, oid *git.Oid) error {
	short := strings.Replace(name, "refs/tags/", "", -1)
	if strings.Contains(short, ti.tagSubStr) {
		obj, _ := ti.repo.Lookup(oid)
		switch obj.Type() {
		case git.ObjectTag: // For annotated tag
			tag, _ := obj.AsTag()
			revision := tag.TargetId().String()
			ti.taggedCommitMap[revision] = append(ti.taggedCommitMap[revision], short)
		case git.ObjectCommit: // For lightweight tag
			revision := obj.Id().String()
			ti.taggedCommitMap[revision] = append(ti.taggedCommitMap[revision], short)
		}
	}
	return nil
}

// GetTaggedCommitMap returns the map which key is revision and the value is the slice of tags attatched to its revision.
func GetTaggedCommitMap(repo *git.Repository, tagSubStr string) (map[string][]string, error) {
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	ti := &tagIterator{
		repo:            repo,
		tagSubStr:       tagSubStr,
		taggedCommitMap: map[string][]string{},
	}
	repo.Tags.Foreach(ti.cb)

	return ti.taggedCommitMap, nil
}

// GetStats returns slice of CommitsStats
func GetStats(repo *git.Repository, taggedCommitMap map[string][]string) ([]CommitStats, error) {
	walk, _ := repo.Walk()
	err := walk.PushHead()
	if err != nil {
		return nil, err
	}

	defaultOpts, _ := git.DefaultDiffOptions()
	defaultDiffOpts, _ := git.DefaultDiffFindOptions()

	csi := &commitStatsIterator{
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
