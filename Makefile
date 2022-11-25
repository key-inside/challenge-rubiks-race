GIT_VERSION = $(shell git describe --always --dirty --tags 2> /dev/null || echo 'unversioned')

APP_NAME = rubiks-race
APP_VERSION ?= $(GIT_VERSION)
PLUGIN_NAME ?= meow

## build: Builds the package
build:
	go build -o build/$(APP_NAME) -ldflags "-s -w"

## build-plugin: Builds the plugin package
build-plugin:
	go build -o build/plugins/$(PLUGIN_NAME).so -ldflags "-s -w" -buildmode=plugin "rubiks-race/plugins/$(PLUGIN_NAME)"

## run: Runs the application
run:
	@./build/$(APP_NAME) $(PLUGIN_NAME) $(PUZZLE_DATA)

## clean: Removes build artifacts and logs
clean:
	rm -rf build
	rm -f rubiks.log

## test: Tests packages
test:
	@go test ./...

## version: Shows the current git version
version:
	@echo $(GIT_VERSION)

help: Makefile
	@sed -n 's/^##//p' $< | awk 'BEGIN {FS = ": "}; {printf "\033[94m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build build-plugin run clean test version help
