FIXLABELS := ./bin/hackenv_fixlabels

BASEDIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build:
	mkdir -p ./out
	go build -ldflags="-w -s" -o ./out/hackenv ./cmd/hackenv

.PHONY: test
test: lint_scripts
	stdout=$$(gofmt -l . 2>&1); \
	if [ "$$stdout" ]; then \
		exit 1; \
	fi
	go vet ./...
	gocyclo -over 10 $(shell find . -iname '*.go' -type f)
	staticcheck ./...
	go test -v -cover ./...

.PHONY: setup
setup:
	go get -u github.com/fzipp/gocyclo/cmd/gocyclo
	go get -u honnef.co/go/tools/cmd/staticcheck

.PHONY: lint_scripts
lint_scripts:
	shellcheck ${FIXLABELS}

.PHONY: install_scripts
install_scripts:
	ln -i -s -r ${BASEDIR}/${FIXLABELS} ${HOME}/bin/
