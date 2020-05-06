PROJECT_NAME := "StoryManager"
PKG := "github.com/AlexisBDL/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/...)
 
.PHONY: all dep lint vet test build clean update_mod
 
all: build

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

build: dep ## Build the binary file
	@go build -i -o build/StoryManager.exe -v $(PKG)

update_mod:
	@go mod tidy
 
clean: ## Remove previous build
	@del /f/q/s build
