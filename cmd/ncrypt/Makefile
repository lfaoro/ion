
VERSION ?= 1.0.0
EPOCH ?= 1
MAINTAINER ?= "Community"

LDFLAGS += -X "main.date=$(shell date '+%Y-%m-%d %I:%M:%S %Z')"

install:
	go install -mod vendor -ldflags='$(LDFLAGS)' .

build:
	@go build -mod vendor .

release:
	goreleaser release --rm-dist

reltest:
	goreleaser release --snapshot --rm-dist --skip-publish

test:
	@go test -mod vendor -cover .

tag?=""
tag:
	git tag -f -a $(tag) -m "$(tag)"
	git push -f origin $(tag)


lines:
	@find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./docs/*" \
	 | xargs wc -l | sort
