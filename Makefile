GO=go
GOFMT=gofmt
DELETE=rm
BINARY=jellybeans-workload
BUILD_BINARY=bin/$(BINARY)
DOCKER=docker
DOCKER_REPO=1xyz/jellybeans-workload
# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
SAFE_BRANCH = $(subst /,-,$(BRANCH))
VER = $(shell git rev-parse --short HEAD)
DOCKER_TAG = "$(SAFE_BRANCH)-$(VER)"

info:
	@echo " target         ⾖ Description.                                    "
	@echo " ----------------------------------------------------------------- "
	@echo
	@echo " build          generate a local build ⇨ $(BUILD_BINARY)          "
	@echo " clean          clean up bin/ & go test cache                      "
	@echo " fmt            format go code files using go fmt                  "
	@echo " release/darwin generate a darwin target build                     "
	@echo " release/linux  generate a linux target build                      "
	@echo " tidy           clean up go module file                            "
	@echo " docker-build        build image $(DOCKER_REPO):$(DOCKER_TAG)      "
	@echo " docker-push         push image $(DOCKER_REPO):$(DOCKER_TAG)       "
	@echo
	@echo " ------------------------------------------------------------------"

build: clean fmt
	$(GO) build -o $(BUILD_BINARY) -v main.go


.PHONY: clean
clean:
	$(DELETE) -rf bin/
	$(GO) clean -cache


.PHONY: fmt
fmt:
	$(GOFMT) -l -w $(SRC)


release/%: clean fmt
	@echo "build no race on alpine. https://github.com/golang/go/issues/14481"
	@echo "build GOOS: $(subst release/,,$@) & GOARCH: amd64"
	GOOS=$(subst release/,,$@) GOARCH=amd64 $(GO) build -o bin/$(subst release/,,$@)/$(BINARY) -v main.go

.PHONY: tidy
tidy:
	$(GO) mod tidy

docker-build:
	$(DOCKER) build -t $(DOCKER_REPO):$(DOCKER_TAG) -f Dockerfile .

docker-push: docker-build
	$(DOCKER) push $(DOCKER_REPO):$(DOCKER_TAG)