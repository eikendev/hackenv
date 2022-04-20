OUT_DIR := ./out

SCRIPTS := $(wildcard ./bin/*)

HE_BUILD_VERSION ?= $(shell git describe --tags)
ifeq ($(HE_BUILD_VERSION),)
	_ := $(error Cannot determine build version)
endif

BUILD_VERSION_FLAG := github.com/eikendev/hackenv/internal/buildconfig.Version=$(HE_BUILD_VERSION)

.PHONY: build
build:
	mkdir -p $(OUT_DIR)
	go build -ldflags "-w -s -X $(BUILD_VERSION_FLAG)" -o $(OUT_DIR)/hackenv ./cmd/hackenv

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

.PHONY: test
test: lint_scripts
	stdout=$$(gofumpt -l . 2>&1); if [ "$$stdout" ]; then exit 1; fi
	go vet ./...
	gocyclo -over 10 $(shell find . -type f -iname '*.go')
	staticcheck ./...
	go test -v -cover ./...
	@printf '\n%s\n' "> Test successful"

.PHONY: setup
setup:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install mvdan.cc/gofumpt@latest

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: lint_scripts
lint_scripts:
	shellcheck ${SCRIPTS}

.PHONY: install_scripts
install_scripts:
	ln -i -s -r ${SCRIPTS} ${HOME}/bin/
