PROJECT = bunny
GOPATH := $(GOPATH)
USER = $(shell id -un)
IMAGE = golang:1.4.2

DOCKER := docker run --rm -v $(PWD):/go/src/github.com/trayio/$(PROJECT) -w /go/src/github.com/trayio/$(PROJECT) -v /etc/passwd:/etc/passwd:ro -v /etc/group:/etc/group:ro -u $(USER):$(USER) $(IMAGE)

test:
	$(DOCKER) go test -v ./...
