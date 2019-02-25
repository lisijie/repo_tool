
PROJECT_NAME := repo_tool

GIT_COMMIT := $(shell git rev-parse --short HEAD || echo unsupported)
GIT_VERSION := $(shell git describe --tags)
BUILD_TIME := $(shell date +%FT%T%z)
LD_FLAGS := -ldflags "-X main.Version=$(GIT_VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: build
build:
	GOARCH=amd64 GOOS=linux go build $(LD_FLAGS) -o $(PROJECT_NAME)
	zip -q $(PROJECT_NAME)-$(GIT_VERSION)-linux-amd64.zip $(PROJECT_NAME) home.html config.ini
	GOARCH=amd64 GOOS=darwin go build $(LD_FLAGS) -o $(PROJECT_NAME)
	zip -q $(PROJECT_NAME)-$(GIT_VERSION)-darwin-amd64.zip $(PROJECT_NAME) home.html config.ini
	rm $(PROJECT_NAME)

.PHONY: run
run: build
	$(TARGET_BIN)

.PHONY: clean
clean:
	@rm -rf *.zip

