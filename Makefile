.DEFAULT_GOAL := all

.PHONY: help run test-unit test-bench test-unit cover-clean

BIN_PATH            := bin
COVER_FILE_PATH     := $(BIN_PATH)/coverage.out
TEST_PATH           := ./...
DOC_ADDR             := :8081

help:  ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean up build files
	rm -rf $(BIN_PATH)

$(BIN_PATH):
	mkdir -p $(BIN_PATH)

$(COVER_FILE_PATH): $(BIN_PATH)
	go test -coverprofile=$(COVER_FILE_PATH) $(TESTS_PATH)

cover:  ## Run coverage tests
	go test -cover $(TESTS_PATH)
	
cover-func: $(COVER_FILE_PATH)  ## Run coverage tests by function
	go tool cover -func=$(COVER_FILE_PATH)

cover-file: $(COVER_FILE_PATH) ## Create a file with coverage test

cover-browser: cover-file ## Show coverage test in a browser
	go tool cover -html=$(COVER_FILE_PATH)

cover-clean: ## Show coverage test in a browser
	rm -rf $(COVER_FILE_PATH)

test-unit: ## Run all unit tests
	go test $(TEST_PATH)

test-unit-race: ## Run all unit tests
	go test -race $(TEST_PATH)

fmt: ## Format all code with gofmt
	gofmt -s -w .

test-bench:
	go test -bench=. -benchtime 10s -count 3 -benchmem ./...

doc: ## Start a go doc server, need to have installed go tools: go get -u golang.org/x/tools/...
	godoc -http $(DOC_ADDR)
