package git

import git "github.com/libgit2/git2go"

type TaggedCommit struct {
	tags      []*git.Oid
	oid       *git.Oid
	next      *git.Oid
	prev      *git.Oid
	commitCnt int
}
