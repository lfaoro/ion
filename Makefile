
VERSION ?= 1.0.0
EPOCH ?= 1
MAINTAINER ?= "Community"


install:
	@go install ./cmd/ncrypt/

build:
	@go build ./cmd/ncrypt/
