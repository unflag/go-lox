GO_CMD=$(shell which go)
ARTIFACT_NAME=lox

generate:
	$(GO_CMD) generate ./...

test: generate
	$(GO_CMD) test -v ./...

build: generate
	$(GO_CMD) build -o $(ARTIFACT_NAME)

run: generate
	$(GO_CMD) run . $(arg)
