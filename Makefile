.PHONY: install clean test

NAME := glstats
VERSION  := 0.0.9
REVISION  := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"github.com/himetani/glstats/cmd.Version=$(VERSION)\" -X \"github.com/himetani/glstats/cmd.Revision=$(REVISION)\""

SRCS    := $(shell find . -path ./vendor -prune -o -name '*.go' -print)
LIBGIT2 := $(shell brew ls libgit2)
GIT2GO := github.com/libgit2/git2go

UNAME := $(shell uname)

ifeq ($(LIBGIT2),)
    LIBGIT2 = must-rebuild
endif

bin/$(NAME): $(SRCS) $(LIBGIT2)
ifeq ($(UNAME), Darwin)
	go build $(LDFLAGS) -o bin/$(NAME)
endif

$(LIBGIT2):
ifeq ($(UNAME), Darwin)
	brew install libgit2
endif

test:
	go test -v github.com/himetani/glstats/...

install: git2go $(LIBGIT2)
	go install $(LDFLAGS)

git2go:
	./script/install.sh

clean:
	rm -rf bin/*
