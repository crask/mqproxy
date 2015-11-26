SHELL=/bin/bash

# Program version
VERSION := $(shell grep "VERSION " main.go | sed -E 's/.*"(.+)"$$/\1/')

# Binary name for bintray
BIN_NAME=kafkaproxy

# Project owner for bintray
OWNER=crask

# Project name for bintray
PROJECT_NAME=$(shell basename $(abspath ./))

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com

# Grab the current commit
GIT_COMMIT="$(shell git rev-parse --short HEAD)"

# Check if there are uncommited changes
GIT_DIRTY="$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

build:
	@echo -e "\033[32;1mBuilding\033[0m \033[33;1m${OWNER} ${BIN_NAME}\033[0m \033[31m${VERSION}-${GIT_COMMIT}${GIT_DIRTY}\033[0m"
	godep go build -ldflags "-X main.gitCommit=${GIT_COMMIT}${GIT_DIRTY}" -o ${BIN_NAME}

clean:
	@test ! -e ${BIN_NAME} || rm -v ${BIN_NAME}

save:
	@echo -e "\033[32;1mUpdating godeps\033[0m"
	godep update `go list -json | grep github.com | grep -v mqproxy | awk '{print $1;}' | tr -d '",'`
	godep save -r
	godep save

.PHONY: build dist clean save
