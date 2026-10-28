OUT_DIR := ./out
GO_FILES := $(shell find . -type f \( -iname '*.go' \))
GO_MODULE := github.com/eikendev/hackenv
GO_TOOL := go tool -modfile=./tools/go.mod

SCRIPT_FILES := $(wildcard ./bin/*)

.PHONY: build
build:
	mkdir -p $(OUT_DIR)
	go build -ldflags "-w -s" -o $(OUT_DIR)/hackenv ./cmd/hackenv

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

.PHONY: test
test: lint_scripts
	if [ -n "$$($(GO_TOOL) gofumpt -l $(GO_FILES))" ]; then echo "Code is not properly formatted"; exit 1; fi
	if [ -n "$$(go fix -diff ./...)" ]; then echo "Code does not use latest idioms of Go"; exit 1; fi
	if [ -n "$$($(GO_TOOL) goimports -l -local $(GO_MODULE) $(GO_FILES))" ]; then echo "Imports are not properly formatted"; exit 1; fi
	go vet ./...
	$(GO_TOOL) misspell -error $(GO_FILES)
	$(GO_TOOL) gocyclo -over 10 $(GO_FILES)
	$(GO_TOOL) staticcheck ./...
	$(GO_TOOL) errcheck -ignoregenerated ./...
	$(GO_TOOL) gocritic check -disable='#experimental,#opinionated' -@ifElseChain.minThreshold 3 ./...
	$(GO_TOOL) revive -set_exit_status ./...
	$(GO_TOOL) nilaway ./...
	go test -v -cover ./...
	$(GO_TOOL) gosec -exclude-generated ./...
	$(GO_TOOL) govulncheck ./...
	@printf '\n%s\n' "> Test successful"

.PHONY: setup
setup:
	go -C ./tools mod download

.PHONY: fmt
fmt:
	$(GO_TOOL) gofumpt -l -w .

.PHONY: lint_scripts
lint_scripts:
	shellcheck ${SCRIPT_FILES}

.PHONY: install_scripts
install_scripts:
	ln -i -s -r ${SCRIPT_FILES} ${HOME}/bin/
