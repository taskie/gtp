.PHONY: build install test fmt coverage dep-init dep-ensure dep-graph pre-commit

CMD_DIR := cmd/gtp

build:
	go build -v -ldflags "-s -w -X fwv.Revision=$(shell git rev-parse --short HEAD)"
	$(MAKE) -C $(CMD_DIR) build

install:
	go install -v -ldflags "-s -w -X fwv.Revision=$(shell git rev-parse --short HEAD)"
	$(MAKE) -C $(CMD_DIR) install

test:
	go test

fmt:
	find . -name '*.go' | xargs gofmt -w

coverage:
	mkdir -p test/coverage
	go test -coverprofile=test/coverage/cover.out
	go tool cover -html=test/coverage/cover.out -o test/coverage/cover.html

dep-init:
	-rm -rf vendor/
	-rm -f Gopkg.toml Gopkg.lock
	dep init

dep-ensure:
	dep ensure

dep-graph:
	mkdir -p images
	dep status -dot | dot -Tpng -o images/dependency.png

pre-commit:
	$(MAKE) dep-init
	$(MAKE) dep-ensure
	$(MAKE) dep-graph
	$(MAKE) fmt
	$(MAKE) build
	$(MAKE) coverage