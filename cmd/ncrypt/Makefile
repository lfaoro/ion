
VERSION ?= 1.0.0
EPOCH ?= 1
MAINTAINER ?= "Community"


install:
	@go install -mod vendor .

build:
	@go build -mod vendor .

release:
	goreleaser release --snapshot --rm-dist


test:
	@go test -mod vendor -cover .

tag?=""
tag:
	git tag -a $(tag) -m "$(tag)"
	git push origin $(tag)

