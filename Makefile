.PHONY: all dep build install clean test

all: ;

NAME := glstats
VERSION  := 0.0.9
REVISION  := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"github.com/himetani/glstats/cmd.Version=$(VERSION)\" -X \"github.com/himetani/glstats/cmd.Revision=$(REVISION)\""

SRCS    := $(shell find . -path ./vendor -prune -o -name '*.go' -print)
LIBGIT2 := $(shell brew ls libgit2)
GIT2GO :=  $(shell find $$GOPATH/src/github.com/libgit2/ -type f -name '*.go')

UNAME := $(shell uname)

ifeq ($(LIBGIT2),)
    LIBGIT2 = must-rebuild
endif

ifeq ($(GIT2GO),)
	GIT2GO = must-rebuild
endif

bin/$(NAME): $(SRCS) dep
ifeq ($(UNAME), Darwin)
	go build $(LDFLAGS) -o bin/$(NAME)
endif

$(GIT2GO):
	./script/install.sh

$(LIBGIT2):
ifeq ($(UNAME), Darwin)
	brew install libgit2
endif

dep: $(LIBGIT2) $(GIT2GO)

build: bin/$(NAME)

install: dep
	go install $(LDFLAGS)

clean:
	rm -rf bin/*

test: dep
	go test -v github.com/himetani/glstats/...
