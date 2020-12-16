# all: run a complete build
.PHONY: all
all: \
	commitlint \
	go-lint \
	go-review \
	go-test \
	go-build \
	go-mod-tidy \
	git-verify-nodiff

export GO111MODULE := on

include tools/commitlint/rules.mk
include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk
include tools/semantic-release/rules.mk

.PHONY: go-build
go-build:
	GOOS=darwin go build ./...
	GOOS=windows go build ./...
	GOOS=linux go build ./...

.PHONY: go-test
go-test:
	go test -cover -race ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v
