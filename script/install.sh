#!/bin/sh

GIT2GO="github.com/libgit2/git2go"

go get $GIT2GO
cd $GOPATH/src/$GIT2GO
git checkout v26
