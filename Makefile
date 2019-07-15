# =================================================================
#
# Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

ifdef GOPATH
GCFLAGS=-trimpath=$(shell printenv GOPATH)/src
else
GCFLAGS=-trimpath=$(shell go env GOPATH)/src
endif

LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)

ifndef DEST
DEST=bin
endif

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z0-9_-\]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# Dependencies
#

deps_go:  ## Install Go dependencies
	go get -d -t ./...

.PHONY: deps_go_test
deps_go_test: ## Download Go dependencies for tests
	go get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow # download shadow
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow # install shadow
	go get -u github.com/kisielk/errcheck # download and install errcheck
	go get -u github.com/client9/misspell/cmd/misspell # download and install misspell
	go get -u github.com/gordonklaus/ineffassign # download and install ineffassign
	go get -u honnef.co/go/tools/cmd/staticcheck # download and instal staticcheck
	go get -u golang.org/x/tools/cmd/goimports # download and install goimports

#
# Go building, formatting, testing, and installing
#

.PHONY: fmt
fmt:  ## Format Go source code
	go fmt $$(go list ./... )

.PHONY: imports
imports: ## Update imports in Go source code
	goimports -w $$(find . -iname '*.go')

.PHONY: vet
vet: ## Vet Go source code
	go vet $$(go list ./... )

.PHONY: test_go
test_go: ## Run Go tests
	bash scripts/test.sh

#
# Clean
#

clean:
	rm -fr bin