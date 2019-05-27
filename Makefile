APP ?= "./cmd/ion"

VERSION ?= 1.0.0
EPOCH ?= 1
MAINTAINER ?= "Community"

LDFLAGS += -X "main.date=$(shell date '+%Y-%m-%d %I:%M:%S %Z')"

install:
	go install -mod vendor -ldflags='$(LDFLAGS)' "$(APP)"

build:
	@go build -mod vendor "$(APP)"

release:
	cd cmd/ion && \
	goreleaser release --rm-dist --config=../../.goreleaser.yml

reltest:
	cd cmd/ion && \
	goreleaser release --snapshot --rm-dist --skip-publish --config=../../.goreleaser.yml

test:
	@go test -mod vendor -cover "$(APP)"

repmod:
	go mod edit -replace=github.com/lfaoro/pkg="$(GOPATH)/src/github.com/lfaoro/pkg"

dropmod:
	go mod edit -dropreplace github.com/lfaoro/pkg

tag?=""
tag:
	git tag -f -a $(tag) -m "$(tag)"
	git push -f origin $(tag)

lines:
	@find . -type f -name "*.go" -not -path "*/vendor/*" -not -path "./docs/*" \
	 | xargs wc -l | sort

codecov:
	go test -race -coverprofile=coverage.txt -covermode=atomic
	bash <(curl -s https://codecov.io/bash) -t b01ec8ef-cab3-428c-8a44-f943035b8714