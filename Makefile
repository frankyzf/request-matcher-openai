# Self documented Makefile
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
UNAME_S=$(shell uname -s)
SED=sed -i
ifeq ($(UNAME_S), Darwin)
	SED=sed -i ''
endif


$(eval GIT_SUMMARY = $(shell  git describe --tags --dirty --always))
$(eval GIT_BRANCH = $(shell  git rev-parse --abbrev-ref HEAD))
$(eval BUILD_TIME = $(shell date +%FT%T+08))
$(eval BUILD_MACHINE = $(shell hostname))
FLAGS:="-X main.GitSummary=$(GIT_SUMMARY) -X main.GitBranch=$(GIT_BRANCH) -X main.BuildTime=$(BUILD_TIME) -X main.BuildMachine=$(BUILD_MACHINE) "


.PHONY: all
all: clean build test

clean:
	@go clean

doc:
	swag init --parseDependency --parseInternal

tidy:
	go mod tidy

build:
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAGS) -o bin/linux/request-matcher-openai .
mac:
	$(GOBUILD) -tags dynamic  -ldflags $(FLAGS) -o bin/request-matcher-openai .
windows:
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -tags 'etcd' -ldflags $(FLAGS) -o bin/windows/request-matcher-openai.exe .


FORCE: ;