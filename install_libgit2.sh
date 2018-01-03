#!/bin/sh

go get -u github.com/libgit2/git2go
cd $GOPATH/src/github.com/libgit2/git2go
git checkout v26
